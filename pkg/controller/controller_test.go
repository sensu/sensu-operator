package controller

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	"github.com/objectrocket/sensu-operator/pkg/cluster"
	fakesensu "github.com/objectrocket/sensu-operator/pkg/generated/clientset/versioned/fake"
	sensuscheme "github.com/objectrocket/sensu-operator/pkg/generated/clientset/versioned/scheme"
	"github.com/objectrocket/sensu-operator/pkg/util/probe"
	"github.com/stretchr/testify/suite"
	fakeapiextensionsapiserver "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	testclient "k8s.io/client-go/kubernetes/fake"
	fakerest "k8s.io/client-go/rest/fake"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type fakeIndexer struct{}
type fakeController struct{}
type fakeQueue struct{}

func (f fakeIndexer) GetByKey(string) (interface{}, bool, error) {
	return nil, false, nil
}

func (f fakeController) Run(stopCh <-chan struct{}) {}
func (f fakeController) HasSynced() bool {
	return true
}

func (f fakeQueue) Add(item interface{})  {}
func (f fakeQueue) Done(item interface{}) {}
func (f fakeQueue) ShutDown()             {}
func (f fakeQueue) Get() (item interface{}, shutdown bool) {
	return nil, true
}
func (f fakeQueue) Forget(interface{}) {}
func (f fakeQueue) NumRequeues(item interface{}) int {
	return 0
}
func (f fakeQueue) AddRateLimited(item interface{}) {}

type InformerTestSuite struct {
	suite.Suite
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func (s *InformerTestSuite) SetupSuite() {

}

func (s *InformerTestSuite) SetupTest() {
	s.ctx, s.cancelCtx = context.WithCancel(context.Background())
}

func (s *InformerTestSuite) TearDownTest() {
	s.cancelCtx()
	select {
	case <-s.ctx.Done():
		return
	case <-time.After(2 * time.Minute):
		s.Fail("Timed out waiting for test to tear down")
	}
}

func TestRunSuite(t *testing.T) {
	suiteTester := new(InformerTestSuite)
	suite.Run(t, suiteTester)
}

func (s *InformerTestSuite) TestInformerWithNoEvents() {
	var (
		source          *cache.ListWatch
		clusterInformer Informer
		assetInformer   Informer
		checkInformer   Informer
		handlerInformer Informer
	)

	controller := New(Config{
		Namespace:         "testns",
		ClusterWide:       true,
		ServiceAccount:    "testsa",
		KubeCli:           testclient.NewSimpleClientset(),
		KubeExtCli:        fakeapiextensionsapiserver.NewSimpleClientset(),
		SensuCRCli:        fakesensu.NewSimpleClientset(),
		CreateCRD:         false,
		WorkerThreads:     1,
		ProcessingRetries: 0,
	})
	assetInformer.indexer = fakeIndexer{}
	assetInformer.controller = fakeController{}
	assetInformer.queue = fakeQueue{}
	checkInformer.indexer = fakeIndexer{}
	checkInformer.controller = fakeController{}
	checkInformer.queue = fakeQueue{}
	handlerInformer.indexer = fakeIndexer{}
	handlerInformer.controller = fakeController{}
	handlerInformer.queue = fakeQueue{}
	controller.informers[api.SensuAssetResourcePlural] = &assetInformer
	controller.informers[api.SensuCheckConfigResourcePlural] = &checkInformer
	controller.informers[api.SensuHandlerResourcePlural] = &handlerInformer

	err := controller.initResource()
	s.Require().NoErrorf(err, "Failed to init resources: %v", err)
	probe.SetReady()

	clusterInformer.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	roundTripper := func(req *http.Request) (*http.Response, error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`
			{
				"apiVersion": "objectrocket.com/v1beta1",
				"items": [],
				"kind": "SensuClusterList",
				"metadata": {
				  "continue": "",
				  "resourceVersion": "3570",
				  "selfLink": "/apis/objectrocket.com/v1beta1/namespaces/default/sensuclusters"
				}
			  }
`)),
			StatusCode: 200,
		}
		response.Header = http.Header{"Content-Type": []string{"application/json"}}
		return response, nil
	}
	controller.Config.SensuCRCli.ObjectrocketV1beta1()
	source = cache.NewListWatchFromClient(
		&fakerest.RESTClient{
			Client: fakerest.CreateHTTPClient(roundTripper),
			NegotiatedSerializer: serializer.DirectCodecFactory{
				CodecFactory: serializer.NewCodecFactory(sensuscheme.Scheme),
			},
			GroupVersion:     schema.GroupVersion{},
			VersionedAPIPath: "/not/a/real/path",
		},
		api.SensuClusterResourcePlural,
		controller.Config.Namespace,
		fields.Everything())
	clusterInformer.indexer, clusterInformer.controller = cache.NewIndexerInformer(source, &api.SensuCluster{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				clusterInformer.queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				clusterInformer.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				clusterInformer.queue.Add(key)
			}
		},
	}, cache.Indexers{})
	controller.informers[api.SensuClusterResourcePlural] = &clusterInformer
	ctx, cancelFunc := context.WithCancel(context.Background())
	go controller.startProcessing(ctx)
	time.Sleep(2 * time.Second)
	cancelFunc()
}

func (s *InformerTestSuite) TestInformerWithOneCluster() {
	var (
		source          *cache.ListWatch
		clusterInformer Informer
		assetInformer   Informer
		checkInformer   Informer
		handlerInformer Informer
	)

	controller := New(Config{
		Namespace:         "testns",
		ClusterWide:       true,
		ServiceAccount:    "testsa",
		KubeCli:           testclient.NewSimpleClientset(),
		KubeExtCli:        fakeapiextensionsapiserver.NewSimpleClientset(),
		SensuCRCli:        fakesensu.NewSimpleClientset(),
		CreateCRD:         false,
		WorkerThreads:     1,
		ProcessingRetries: 0,
	})
	assetInformer.indexer = fakeIndexer{}
	assetInformer.controller = fakeController{}
	assetInformer.queue = fakeQueue{}
	checkInformer.indexer = fakeIndexer{}
	checkInformer.controller = fakeController{}
	checkInformer.queue = fakeQueue{}
	handlerInformer.indexer = fakeIndexer{}
	handlerInformer.controller = fakeController{}
	handlerInformer.queue = fakeQueue{}
	controller.informers[api.SensuAssetResourcePlural] = &assetInformer
	controller.informers[api.SensuCheckConfigResourcePlural] = &checkInformer
	controller.informers[api.SensuHandlerResourcePlural] = &handlerInformer
	err := controller.initResource()
	s.Require().NoErrorf(err, "Failed to init resources: %v", err)
	probe.SetReady()

	clusterInformer.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	roundTripper := func(req *http.Request) (*http.Response, error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`
			{
				"apiVersion": "objectrocket.com/v1beta1",
				"items": [
				  {
					"apiVersion": "objectrocket.com/v1beta1",
					"kind": "SensuCluster",
					"metadata": {
					  "annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"objectrocket.com/v1beta1\",\"kind\":\"SensuCluster\",\"metadata\":{\"annotations\":{},\"name\":\"example-sensu-cluster\",\"namespace\":\"default\"},\"spec\":{\"size\":3,\"version\":\"2.0.0-beta.8\"}}\n"
					  },
					  "clusterName": "",
					  "creationTimestamp": "2019-01-02T23:14:52Z",
					  "generation": 1,
					  "name": "example-sensu-cluster",
					  "namespace": "default",
					  "resourceVersion": "3570",
					  "selfLink": "/apis/objectrocket.com/v1beta1/namespaces/default/sensuclusters/example-sensu-cluster",
					  "uid": "358db0b6-0ee4-11e9-a33b-0800272dcccb"
					},
					"spec": {
					  "repository": "sensu/sensu",
					  "size": 3,
					  "version": "2.0.0-beta.8"
					},
					"status": {
					  "agentPort": 8081,
					  "agentServiceName": "example-sensu-cluster-agent",
					  "apiPort": 8080,
					  "apiServiceName": "example-sensu-cluster-api",
					  "conditions": [
						{
						  "lastTransitionTime": "2019-01-02T23:15:48Z",
						  "lastUpdateTime": "2019-01-02T23:15:48Z",
						  "reason": "Cluster available",
						  "status": "True",
						  "type": "Available"
						}
					  ],
					  "currentVersion": "2.0.0-beta.8",
					  "dashboardPort": 3000,
					  "dashboardServiceName": "example-sensu-cluster-dashboard",
					  "members": {
						"ready": [
						  "example-sensu-cluster-6h5wp5t264",
						  "example-sensu-cluster-8ldr4vhlz5",
						  "example-sensu-cluster-b4cf6wcnpc"
						]
					  },
					  "phase": "Running",
					  "size": 3,
					  "targetVersion": ""
					}
				  }
				],
				"kind": "SensuClusterList",
				"metadata": {
				  "continue": "",
				  "resourceVersion": "3570",
				  "selfLink": "/apis/objectrocket.com/v1beta1/namespaces/default/sensuclusters"
				}
			  }
`)),
			StatusCode: 200,
		}
		response.Header = http.Header{"Content-Type": []string{"application/json"}}
		return response, nil
	}
	controller.Config.SensuCRCli.ObjectrocketV1beta1()
	source = cache.NewListWatchFromClient(
		&fakerest.RESTClient{
			Client: fakerest.CreateHTTPClient(roundTripper),
			NegotiatedSerializer: serializer.DirectCodecFactory{
				CodecFactory: serializer.NewCodecFactory(sensuscheme.Scheme),
			},
			GroupVersion:     schema.GroupVersion{},
			VersionedAPIPath: "/not/a/real/path",
		},
		api.SensuClusterResourcePlural,
		controller.Config.Namespace,
		fields.Everything())
	clusterInformer.indexer, clusterInformer.controller = cache.NewIndexerInformer(source, &api.SensuCluster{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			//s.Failf("Failed with obj:", " %v", obj)
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				clusterInformer.queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				clusterInformer.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				clusterInformer.queue.Add(key)
			}
		},
	}, cache.Indexers{})
	ctx, cancelFunc := context.WithCancel(context.Background())
	controller.informers[api.SensuClusterResourcePlural] = &clusterInformer
	go controller.startProcessing(ctx)
	time.Sleep(2 * time.Second)
	cancelFunc()
}

func initInformers() (Informer, Informer, Informer) {
	var (
		assetInformer   Informer
		checkInformer   Informer
		handlerInformer Informer
	)
	assetInformer.indexer = fakeIndexer{}
	assetInformer.controller = fakeController{}
	assetInformer.queue = fakeQueue{}
	checkInformer.indexer = fakeIndexer{}
	checkInformer.controller = fakeController{}
	checkInformer.queue = fakeQueue{}
	handlerInformer.indexer = fakeIndexer{}
	handlerInformer.controller = fakeController{}
	handlerInformer.queue = fakeQueue{}
	return assetInformer, checkInformer, handlerInformer
}

func TestController_initCRD(t *testing.T) {
	assetInformer, checkInformer, handlerInformer := initInformers()
	type fields struct {
		logger     *logrus.Entry
		Config     Config
		informers  map[string]*Informer
		finalizers map[string]cache.Indexer
		clusters   map[string]*cluster.Cluster
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"test both cluster, and checkconfig crds are created, and become valid",
			fields{
				logrus.WithField("pkg", "test"),
				Config{
					Namespace:         "testns",
					ClusterWide:       true,
					ServiceAccount:    "testsa",
					KubeCli:           testclient.NewSimpleClientset(),
					KubeExtCli:        fakeapiextensionsapiserver.NewSimpleClientset(),
					SensuCRCli:        fakesensu.NewSimpleClientset(),
					CreateCRD:         false,
					WorkerThreads:     1,
					ProcessingRetries: 0,
				},
				map[string]*Informer{},
				map[string]cache.Indexer{},
				map[string]*cluster.Cluster{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				logger:     tt.fields.logger,
				Config:     tt.fields.Config,
				informers:  tt.fields.informers,
				finalizers: tt.fields.finalizers,
				clusters:   tt.fields.clusters,
			}
			c.informers[api.SensuAssetResourcePlural] = &assetInformer
			c.informers[api.SensuCheckConfigResourcePlural] = &checkInformer
			c.informers[api.SensuHandlerResourcePlural] = &handlerInformer
			if err := c.initCRD(); (err != nil) != tt.wantErr {
				t.Errorf("Controller.initCRD() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
