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
	"time"

	api "github.com/objectrocket/sensu-operator/pkg/apis/sensu/v1beta1"
	"github.com/objectrocket/sensu-operator/pkg/util/k8sutil"
	"github.com/objectrocket/sensu-operator/pkg/util/probe"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
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
	var (
		ns     string
		source *cache.ListWatch
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
	c.logger.Infof("Created resources")
	probe.SetReady()

	if c.Config.ClusterWide {
		ns = metav1.NamespaceAll
	} else {
		ns = c.Config.Namespace
	}
	c.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	defer c.queue.ShutDown()
	source = cache.NewListWatchFromClient(
		c.Config.SensuCRCli.SensuV1beta1().RESTClient(),
		api.SensuClusterResourcePlural,
		ns,
		fields.Everything())
	c.indexer, c.informer = cache.NewIndexerInformer(source, &api.SensuCluster{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				c.queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				c.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				c.queue.Add(key)
			}
		},
	}, cache.Indexers{})
	c.logger.Infof("Created index and informer")
	c.startProcessing(ctx)
}

func (c *Controller) startProcessing(ctx context.Context) {
	var (
		cacheStopCh chan struct{}
		untilStopCh chan struct{}
	)
	cacheStopCh = make(chan struct{})
	untilStopCh = make(chan struct{})
	infCtx, infCancel := context.WithCancel(ctx)
	go c.informer.Run(infCtx.Done())
	c.logger.Infof("Started informer")
	if !cache.WaitForCacheSync(infCtx.Done(), c.informer.HasSynced) {
		c.logger.Fatal("Timed out waiting for caches to sync")
	}
	for i := 0; i < c.Config.WorkerThreads; i++ {
		go wait.Until(c.run, time.Second, untilStopCh)
	}
	c.logger.Infof("Started worker %d threads", c.Config.WorkerThreads)
	select {
	case <-ctx.Done():
		c.logger.Errorf("Stopping channels")
		close(cacheStopCh)
		close(untilStopCh)
		infCancel()
	}
	c.logger.Infof("Finished startProcessing")
}

func (c *Controller) run() {
	for c.processNextItem() {
		c.logger.Infof("Processed next item")
	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)
	obj, exists, err := c.indexer.GetByKey(key.(string))
	if err != nil {
		if c.queue.NumRequeues(key) < c.Config.ProcessingRetries {
			c.logger.Errorf("Error fetching object with key %s from store: %v. Retrying later.", key, err)
			c.queue.AddRateLimited(key)
			return true
		}
		c.logger.Errorf("Fatal failure fetching object with key %s from store: %v", key, err)
	} else {
		c.logger.Infof("Object: %v, exists: %t", obj, exists)
		if !exists {
			c.onDeleteSensuClus(obj)
		} else {
			c.onUpdateSensuClus(obj)
		}
	}
	c.queue.Forget(key)
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
