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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kwatch "k8s.io/apimachinery/pkg/watch"
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
	var ns string
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
	c.startProcessing(ctx)
}

func (c *Controller) startProcessing(ctx context.Context) {
	var (
		clusterController hasSynced
		assetController   hasSynced
	)
	clusterController = c.informers[api.SensuClusterResourcePlural].controller
	assetController = c.informers[api.SensuAssetResourcePlural].controller
	go clusterController.Run(ctx.Done())
	go assetController.Run(ctx.Done())
	if !cache.WaitForCacheSync(ctx.Done(), clusterController.HasSynced) {
		c.logger.Fatal("Timed out waiting for cluster caches to sync")
	}
	if !cache.WaitForCacheSync(ctx.Done(), assetController.HasSynced) {
		c.logger.Fatal("Timed out waiting for asset caches to sync")
	}
	for i := 0; i < c.Config.WorkerThreads; i++ {
		go wait.Until(c.run, time.Second, ctx.Done())
	}
	select {
	case <-ctx.Done():
	}
}

func (c *Controller) addInformer(namespace string, resourcePlural string, objType runtime.Object) {
	var (
		informer Informer
		source   *cache.ListWatch
	)
	informer.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	source = cache.NewListWatchFromClient(
		c.Config.SensuCRCli.ObjectrocketV1beta1().RESTClient(),
		resourcePlural,
		namespace,
		fields.Everything())
	informer.indexer, informer.controller = cache.NewIndexerInformer(source, objType, 0, cache.ResourceEventHandlerFuncs{
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
	}, cache.Indexers{})
	c.informers[resourcePlural] = &informer
}

func (c *Controller) run() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer c.informers[api.SensuClusterResourcePlural].queue.ShutDown()
		for c.processNextClusterItem() {
		}
	}()
	go func() {
		defer wg.Done()
		defer c.informers[api.SensuAssetResourcePlural].queue.ShutDown()
		for c.processNextAssetItem() {
		}
	}()
	wg.Wait()
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
			c.onDeleteSensuClus(obj)
		} else {
			c.onUpdateSensuClus(obj)
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
			c.onDeleteSensuAsset(obj)
		} else {
			c.onUpdateSensuAsset(obj)
		}
	}
	assetInformer.queue.Forget(key)
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
	ev := &Event{
		Type:   kwatch.Deleted,
		Object: clus,
	}

	pt.start()
	_, err := c.handleClusterEvent(ev)
	if err != nil {
		c.logger.Warningf("fail to handle event: %v", err)
	}
	pt.stop()
}

func (c *Controller) syncSensuClus(clus *api.SensuCluster) {
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

	pt.start()
	_, err := c.handleClusterEvent(ev)
	if err != nil {
		c.logger.Warningf("fail to handle event: %v", err)
	}
	pt.stop()
}

func (c *Controller) onUpdateSensuAsset(newObj interface{}) {
	c.syncSensuAsset(newObj.(*api.SensuAsset))
}

func (c *Controller) onDeleteSensuAsset(obj interface{}) {
	//TODO: Implement
	c.logger.Warnf("Deleting SensuAssets not implemented.  Not deleting: %v", obj)
}

func (c *Controller) syncSensuAsset(obj *api.SensuAsset) {
	//TODO: Implement
	c.logger.Warnf("Syncing SensuAssets not implemented.  Not syncing: %v", obj)
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
