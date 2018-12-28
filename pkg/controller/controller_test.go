package controller

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	api "github.com/objectrocket/sensu-operator/pkg/apis/sensu/v1beta1"
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

type FakeInformer struct {
	controller cache.Controller
	hasSynced  bool
}

func (f *FakeInformer) Run(stopCh <-chan struct{}) {
	f.controller.Run(stopCh)
}

func (f *FakeInformer) HasSynced() bool {
	return f.hasSynced
}

func (f *FakeInformer) LastSyncResourceVersion() string {
	return f.controller.LastSyncResourceVersion()
}

func (s *InformerTestSuite) TestInformerWithNoEvents() {
	var (
		source *cache.ListWatch
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

	err := controller.initResource()
	s.Require().NoErrorf(err, "Failed to init resources: %v", err)
	probe.SetReady()

	controller.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	defer controller.queue.ShutDown()
	roundTripper := func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
		}, nil
	}
	controller.Config.SensuCRCli.SensuV1beta1()
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
	if controller.indexer == nil || controller.informer == nil {
		controller.indexer, controller.informer = cache.NewIndexerInformer(source, &api.SensuCluster{}, 0, cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err == nil {
					controller.queue.Add(key)
				}
			},
			UpdateFunc: func(old interface{}, new interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(new)
				if err == nil {
					controller.queue.Add(key)
				}
			},
			DeleteFunc: func(obj interface{}) {
				// IndexerInformer uses a delta queue, therefore for deletes we have to use this
				// key function.
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				if err == nil {
					controller.queue.Add(key)
				}
			},
		}, cache.Indexers{})
	}
	controller.informer = &FakeInformer{
		hasSynced:  true,
		controller: controller.informer,
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	controller.startProcessing(ctx)
	time.Sleep(5 * time.Second)
	cancelFunc()
}
