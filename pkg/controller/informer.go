// Copyright 2017 The etcd-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	"github.com/objectrocket/sensu-operator/pkg/util/k8sutil"
	"github.com/objectrocket/sensu-operator/pkg/util/probe"

	corev1 "k8s.io/api/core/v1"

	// "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kwatch "k8s.io/apimachinery/pkg/watch"

	// informers_corev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// TODO: get rid of this once we use workqueue
var pt *panicTimer

func init() {
	pt = newPanicTimer(time.Minute, "unexpected long blocking (> 1 Minute) when handling cluster event")
}

// Start the controller's informer to watch for custom resource update
func (c *Controller) Start(ctx context.Context) {
	var (
		ns string
	)
	// TODO: get rid of this init code. CRD and storage class will be managed outside of operator.
	for {
		err := c.initResource()
		if err == nil {
			break
		}
		c.logger.Errorf("initialization failed: %v", err)
		c.logger.Infof("retry in %v...", initRetryWaitTime)
		time.Sleep(initRetryWaitTime)
	}
	probe.SetReady()

	if c.Config.ClusterWide {
		ns = metav1.NamespaceAll
	} else {
		ns = c.Config.Namespace
	}
	c.addInformer(ns, api.SensuClusterResourcePlural, &api.SensuCluster{})
	c.addInformer(ns, api.SensuAssetResourcePlural, &api.SensuAsset{})
	c.addInformer(ns, api.SensuCheckConfigResourcePlural, &api.SensuCheckConfig{})
	c.addInformer(ns, api.SensuHandlerResourcePlural, &api.SensuHandler{})
	c.addInformer(ns, api.SensuEventFilterResourcePlural, &api.SensuEventFilter{})
	c.addInformerWithCacheGetter(c.Config.KubeCli.CoreV1().RESTClient(), metav1.NamespaceAll, "nodes", &corev1.Node{})
	// c.addNodeInformer()
	c.startProcessing(ctx)
}

func (c *Controller) startProcessing(ctx context.Context) {
	var (
		clusterController     hasSynced
		assetController       hasSynced
		checkconfigController hasSynced
		handlerController     hasSynced
		eventFilterController hasSynced
		nodeController        hasSynced
	)
	clusterController = c.informers[api.SensuClusterResourcePlural].controller
	assetController = c.informers[api.SensuAssetResourcePlural].controller
	checkconfigController = c.informers[api.SensuCheckConfigResourcePlural].controller
	handlerController = c.informers[api.SensuHandlerResourcePlural].controller
	eventFilterController = c.informers[api.SensuEventFilterResourcePlural].controller
	nodeController = c.informers["nodes"].controller
	go clusterController.Run(ctx.Done())
	go assetController.Run(ctx.Done())
	go checkconfigController.Run(ctx.Done())
	go handlerController.Run(ctx.Done())
	go eventFilterController.Run(ctx.Done())
	c.logger.Debugf("about to call run on node Controller %+v", nodeController)
	go nodeController.Run(ctx.Done())
	c.logger.Debugf("waiting for cluster controller caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), clusterController.HasSynced) {
		c.logger.Fatal("Timed out waiting for cluster caches to sync")
	}
	c.logger.Debugf("waiting for assetController caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), assetController.HasSynced) {
		c.logger.Fatal("Timed out waiting for asset caches to sync")
	}
	c.logger.Debugf("waiting for checkconfigController caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), checkconfigController.HasSynced) {
		c.logger.Fatal("Timed out waiting for checkconfig caches to sync")
	}
	c.logger.Debugf("waiting for handlerController caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), handlerController.HasSynced) {
		c.logger.Fatal("Timed out waiting for handler caches to sync")
	}
	c.logger.Debugf("waiting for eventFilterController caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), eventFilterController.HasSynced) {
		c.logger.Fatal("Timed out waiting for event filter caches to sync")
	}
	c.logger.Debugf("waiting for nodeController caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), nodeController.HasSynced) {
		c.logger.Fatal("Timed out waiting for node caches to sync")
	}
	c.logger.Debugf("finished waiting for all caches to sync")
	c.logger.Debugf("about to spawn %d worker threads", c.WorkerThreads)
	for i := 0; i < c.Config.WorkerThreads; i++ {
		go wait.Until(c.run, time.Second, ctx.Done())
	}
	c.logger.Debugf("done spawning %d worker threads", c.WorkerThreads)
	select {
	case <-ctx.Done():
	}
}

func (c *Controller) addNodeInformer() {
	var (
		informer Informer
	)
	c.logger.Debugf("starting adding node informer")
	c.logger.Debugf("getting new rate limiting queue")
	informer.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	source := cache.NewListWatchFromClient(
		c.KubeCli.CoreV1().RESTClient(),
		"nodes",
		"",
		fields.Everything())
	c.logger.Debugf("getting new node shared index informer")
	sharedInformer := cache.NewSharedIndexInformer(
		source,
		&corev1.Node{},
		c.Config.ResyncPeriod,
		cache.Indexers{},
	)
	// sharedInformer := informers_corev1.NewNodeInformer(c.Config.KubeCli, c.Config.ResyncPeriod, cache.Indexers{})
	c.logger.Debugf("got new node shared index informer: %+v", sharedInformer)
	c.logger.Debugf("adding event handlers to node shared index informer")
	sharedInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				informer.queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				informer.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				informer.queue.Add(key)
			}
		},
	})
	informer.controller = sharedInformer.GetController()
	informer.indexer = sharedInformer.GetIndexer()
	c.logger.Debugf("setting informers map for nodes")
	c.informers["nodes"] = &informer
	c.logger.Debugf("done setting up node shared index informer")
}

func (c *Controller) addInformerWithCacheGetter(getter cache.Getter, namespace, resourcePlural string, objType runtime.Object) {
	var (
		informer Informer
		source   *cache.ListWatch
	)
	informer.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	source = cache.NewListWatchFromClient(
		getter,
		resourcePlural,
		namespace,
		fields.Everything())
	// create finalizer to ensure that sensu server objects are deleted when crd is deleted
	finalizer := cache.NewIndexer(cache.DeletionHandlingMetaNamespaceKeyFunc, cache.Indexers{})
	c.logger.Debugf("getting a new indexer informer with resync period of %s", c.ResyncPeriod.String())
	informer.indexer, informer.controller = cache.NewIndexerInformer(source, objType, c.ResyncPeriod, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				informer.queue.Add(key)
				finalizer.Delete(obj)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				informer.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				finalizer.Add(obj)
				informer.queue.Add(key)
			}
		},
	}, cache.Indexers{})
	c.informers[resourcePlural] = &informer
	c.finalizers[resourcePlural] = finalizer
}

func (c *Controller) addInformer(namespace string, resourcePlural string, objType runtime.Object) {
	c.addInformerWithCacheGetter(c.Config.SensuCRCli.ObjectrocketV1beta1().RESTClient(), namespace, resourcePlural, objType)
}

func (c *Controller) run() {
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		defer c.informers[api.SensuClusterResourcePlural].queue.ShutDown()
		c.logger.Debugf("starting processing of cluster items")
		for c.processNextClusterItem() {
		}
	}()
	go func() {
		defer wg.Done()
		defer c.informers[api.SensuAssetResourcePlural].queue.ShutDown()
		c.logger.Debugf("starting processing of asset items")
		for c.processNextAssetItem() {
		}
	}()
	go func() {
		defer wg.Done()
		defer c.informers[api.SensuCheckConfigResourcePlural].queue.ShutDown()
		c.logger.Debugf("starting processing of checkconfig items")
		for c.processNextCheckConfigItem() {
		}
	}()
	go func() {
		defer wg.Done()
		defer c.informers[api.SensuHandlerResourcePlural].queue.ShutDown()
		c.logger.Debugf("starting processing of handler items")
		for c.processNextHandlerItem() {
		}
	}()
	go func() {
		defer wg.Done()
		defer c.informers[api.SensuEventFilterResourcePlural].queue.ShutDown()
		c.logger.Debugf("starting processing of event filter items")
		for c.processNextEventFilterItem() {
		}
	}()
	go func() {
		defer wg.Done()
		defer c.informers["nodes"].queue.ShutDown()
		c.logger.Debugf("starting processing of node items")
		for c.processNextNodeItem() {
		}
	}()
	c.logger.Debugf("waiting on all waitgroups to finish")
	wg.Wait()
	c.logger.Debugf("waitgroups finished")
}

func (c *Controller) processNextClusterItem() bool {
	var clusterInformer = c.informers[api.SensuClusterResourcePlural]
	key, quit := clusterInformer.queue.Get()
	if quit {
		return false
	}
	defer clusterInformer.queue.Done(key)
	obj, exists, err := clusterInformer.indexer.GetByKey(key.(string))
	if err != nil {
		if clusterInformer.queue.NumRequeues(key) < c.Config.ProcessingRetries {
			clusterInformer.queue.AddRateLimited(key)
			return true
		}
	} else {
		if !exists {
			_, exists, err := c.finalizers[api.SensuClusterResourcePlural].GetByKey(key.(string))
			if exists && err != nil {
				c.finalizers[api.SensuClusterResourcePlural].Delete(key)
			}
		} else {
			if obj != nil {
				c.onUpdateSensuClus(obj)
				cluster := obj.(*api.SensuCluster)
				if cluster.DeletionTimestamp != nil {
					c.onDeleteSensuClus(obj)
				}
			}
		}
	}
	clusterInformer.queue.Forget(key)
	return true
}

func (c *Controller) processNextAssetItem() bool {
	var assetInformer = c.informers[api.SensuAssetResourcePlural]
	key, quit := assetInformer.queue.Get()
	if quit {
		return false
	}
	defer assetInformer.queue.Done(key)
	obj, exists, err := assetInformer.indexer.GetByKey(key.(string))
	if err != nil {
		if assetInformer.queue.NumRequeues(key) < c.Config.ProcessingRetries {
			assetInformer.queue.AddRateLimited(key)
			return true
		}
	} else {
		if !exists {
			_, exists, err := c.finalizers[api.SensuAssetResourcePlural].GetByKey(key.(string))
			if exists && err != nil {
				c.finalizers[api.SensuAssetResourcePlural].Delete(key)
			}
		} else {
			if obj != nil {
				c.onUpdateSensuAsset(obj)
				asset := obj.(*api.SensuAsset)
				// If asset deletion has been initiated, also delete asset from sensu cluster
				if asset.DeletionTimestamp != nil {
					c.onDeleteSensuAsset(obj)
				}
			}
		}
	}
	assetInformer.queue.Forget(key)
	return true
}

func (c *Controller) processNextCheckConfigItem() bool {
	var checkconfigInformer = c.informers[api.SensuCheckConfigResourcePlural]
	key, quit := checkconfigInformer.queue.Get()
	if quit {
		return false
	}
	defer checkconfigInformer.queue.Done(key)
	obj, exists, err := checkconfigInformer.indexer.GetByKey(key.(string))
	if err != nil {
		if checkconfigInformer.queue.NumRequeues(key) < c.Config.ProcessingRetries {
			checkconfigInformer.queue.AddRateLimited(key)
			return true
		}
	} else {
		if !exists {
			_, exists, err := c.finalizers[api.SensuCheckConfigResourcePlural].GetByKey(key.(string))
			if exists && err != nil {
				c.finalizers[api.SensuCheckConfigResourcePlural].Delete(key)
			}
		} else {
			if obj != nil {
				c.onUpdateSensuCheckConfig(obj)
				checkconfig := obj.(*api.SensuCheckConfig)
				// If checkconfig deletion has been initiated, also delete checkconfig from sensu cluster
				if checkconfig.DeletionTimestamp != nil {
					c.onDeleteSensuCheckConfig(obj)
				}
			}
		}
	}
	checkconfigInformer.queue.Forget(key)
	return true
}

func (c *Controller) processNextHandlerItem() bool {
	var handlerInformer = c.informers[api.SensuHandlerResourcePlural]
	key, quit := handlerInformer.queue.Get()
	if quit {
		return false
	}
	defer handlerInformer.queue.Done(key)
	obj, exists, err := handlerInformer.indexer.GetByKey(key.(string))
	if err != nil {
		if handlerInformer.queue.NumRequeues(key) < c.Config.ProcessingRetries {
			handlerInformer.queue.AddRateLimited(key)
			return true
		}
	} else {
		if !exists {
			_, exists, err := c.finalizers[api.SensuHandlerResourcePlural].GetByKey(key.(string))
			if exists && err != nil {
				c.finalizers[api.SensuHandlerResourcePlural].Delete(key)
			}
		} else {
			if obj != nil {
				c.onUpdateSensuHandler(obj.(*api.SensuHandler))
				handler := obj.(*api.SensuHandler)
				// If checkconfig deletion has been initiated, also delete checkconfig from sensu cluster
				if handler.DeletionTimestamp != nil {
					c.onDeleteSensuHandler(obj)
				}
			}
		}
	}
	handlerInformer.queue.Forget(key)
	return true
}

func (c *Controller) processNextEventFilterItem() bool {
	var eventFilterInformer = c.informers[api.SensuEventFilterResourcePlural]
	key, quit := eventFilterInformer.queue.Get()
	if quit {
		return false
	}
	defer eventFilterInformer.queue.Done(key)
	obj, exists, err := eventFilterInformer.indexer.GetByKey(key.(string))
	if err != nil {
		if eventFilterInformer.queue.NumRequeues(key) < c.Config.ProcessingRetries {
			eventFilterInformer.queue.AddRateLimited(key)
			return true
		}
	} else {
		if !exists {
			_, exists, err := c.finalizers[api.SensuEventFilterResourcePlural].GetByKey(key.(string))
			if exists && err != nil {
				c.finalizers[api.SensuEventFilterResourcePlural].Delete(key)
			}
		} else {
			if obj != nil {
				c.onUpdateSensuEventFilter(obj.(*api.SensuEventFilter))
				filter := obj.(*api.SensuEventFilter)
				// If filter deletion has been initiated, also delete filter from sensu cluster
				if filter.DeletionTimestamp != nil {
					c.onDeleteSensuEventFilter(obj)
				}
			}
		}
	}
	eventFilterInformer.queue.Forget(key)
	return true
}

func (c *Controller) processNextNodeItem() bool {
	var nodesInformer = c.informers["nodes"]
	c.logger.Debugf("in processNextNodeItem: getting key from queue")
	key, quit := nodesInformer.queue.Get()
	c.logger.Debugf("got key '%s', quit '%t' from queue", key, quit)
	if quit {
		c.logger.Debugf("got quit while processing next node item, returning")
		return false
	}
	defer nodesInformer.queue.Done(key)
	c.logger.Debugf("getting object by key from index informer")
	obj, exists, err := nodesInformer.indexer.GetByKey(key.(string))
	c.logger.Debugf("got object %+v, exists %t, err %+v from index informer", obj, exists, err)
	if err != nil {
		c.logger.Debugf("got error while getting object by key from index informer")
		if nodesInformer.queue.NumRequeues(key) < c.Config.ProcessingRetries {
			c.logger.Debugf("retries left, so re-queueing item")
			nodesInformer.queue.AddRateLimited(key)
			return true
		}
	} else {
		c.logger.Debugf("no error while getting object by key from index informer")
		if exists {
			c.logger.Debugf("object exists while getting object bby key from index informer")
			if obj != nil {
				c.logger.Debugf("object ! nil while getting object bby key from index informer")
				c.logger.Debugf("calling onUpdateNode for node %+v", obj)
				c.onUpdateNode(obj.(*corev1.Node))
			}
		} else {
			c.logger.Debugf("node %s appears to have been deleted, so calling onDeleteNode", key.(string))
			c.onDeleteNode(key.(string))
		}
	}
	nodesInformer.queue.Forget(key)
	return true
}

func (c *Controller) initResource() error {
	if c.Config.CreateCRD {
		err := c.initCRD()
		if err != nil {
			return fmt.Errorf("fail to init CRD: %v", err)
		}
	}
	return nil
}

func (c *Controller) onUpdateSensuClus(newObj interface{}) {
	c.syncSensuClus(newObj.(*api.SensuCluster))
}

func (c *Controller) onDeleteSensuClus(obj interface{}) {
	var err error
	clus, ok := obj.(*api.SensuCluster)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			panic(fmt.Sprintf("unknown object from SensuCluster delete event: %#v", obj))
		}
		clus, ok = tombstone.Obj.(*api.SensuCluster)
		if !ok {
			panic(fmt.Sprintf("Tombstone contained object that is not a SensuCluster: %#v", obj))
		}
	}

	clus.Finalizers = make([]string, 0)
	clus, err = c.SensuCRCli.ObjectrocketV1beta1().SensuClusters(clus.GetNamespace()).Update(clus)
	if err != nil {
		c.logger.Warningf("Failed to remove finalizers from cluster: %v", err)
	}
	ev := &Event{
		Type:   kwatch.Deleted,
		Object: clus,
	}
	pt.start()
	_, err = c.handleClusterEvent(ev)
	if err != nil {
		c.logger.Warningf("fail to handle event: %v", err)
	}
	pt.stop()
}

func (c *Controller) syncSensuClus(clus *api.SensuCluster) {
	var err error
	ev := &Event{
		Type:   kwatch.Added,
		Object: clus,
	}
	// re-watch or restart could give ADD event.
	// If for an ADD event the cluster spec is invalid then it is not added to the local cache
	// so modifying that cluster will result in another ADD event
	if _, ok := c.clusters[clus.Name]; ok {
		ev.Type = kwatch.Modified
	}
	// Ensure that the finalizer exists, failing if it can't be added at this time
	if len(clus.Finalizers) == 0 && clus.DeletionTimestamp == nil {
		clus.Finalizers = append(clus.Finalizers, "cluster.finalizer.objectrocket.com")
		if clus, err = c.SensuCRCli.ObjectrocketV1beta1().SensuClusters(clus.GetNamespace()).Update(clus); err != nil {
			msg := fmt.Sprintf("failed to update clusters's finalizer during sync event: %v", err)
			c.logger.Warningf(msg)
			return
		}
	}
	pt.start()
	_, err = c.handleClusterEvent(ev)
	if err != nil {
		c.logger.Warningf("fail to handle event: %v", err)
	}
	pt.stop()
}

func (c *Controller) managed(clus *api.SensuCluster) bool {
	if v, ok := clus.Annotations[k8sutil.AnnotationScope]; ok {
		if c.Config.ClusterWide {
			return v == k8sutil.AnnotationClusterWide
		}
	} else {
		if !c.Config.ClusterWide {
			return true
		}
	}
	return false
}
