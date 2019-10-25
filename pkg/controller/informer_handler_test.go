package controller

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	"github.com/objectrocket/sensu-operator/pkg/cluster"
	fakesensu "github.com/objectrocket/sensu-operator/pkg/generated/clientset/versioned/fake"
	fakeapiextensionsapiserver "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

func TestController_syncSensuHandler(t *testing.T) {
	assetInformer, checkInformer, handlerInformer, eventFilterHandler, nodeInformer := initInformers()
	type fields struct {
		logger     *logrus.Entry
		Config     Config
		informers  map[string]*Informer
		finalizers map[string]cache.Indexer
		clusters   map[string]*cluster.Cluster
	}
	tests := []struct {
		name     string
		fields   fields
		handler  *api.SensuHandler
		initFunc func(*testing.T, *Controller, *api.SensuHandler)
		testFunc func(*Controller, *api.SensuHandler) error
	}{
		{
			"test handler status is updated when cluster is missing",
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
			&api.SensuHandler{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testHandler",
					Namespace: "sensu",
				},
				Spec: api.SensuHandlerSpec{
					SensuMetadata: api.ObjectMeta{
						Name:        "testHandler",
						Namespace:   "default",
						ClusterName: "doesntExist",
					},
				},
			},
			func(t *testing.T, c *Controller, handler *api.SensuHandler) {
				_, err := c.SensuCRCli.ObjectrocketV1beta1().SensuHandlers("sensu").Create(handler)
				if err != nil {
					t.Error(err)
				}
			},
			func(c *Controller, handler *api.SensuHandler) error {
				c.syncSensuHandler(handler)
				k8sHandler, err := c.SensuCRCli.ObjectrocketV1beta1().SensuHandlers("sensu").Get("testHandler", metav1.GetOptions{})
				if err != nil {
					return errors.New("failed to find handler in k8s")
				}
				if k8sHandler.Status.Accepted {
					return errors.New("handler status should not be accepted=true")
				}
				if k8sHandler.Status.LastError == "" {
					return errors.New("handler status lastError should not be empty")
				}
				return nil
			},
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
			c.informers[api.SensuEventFilterResourcePlural] = &eventFilterHandler
			c.informers[CoreV1NodesPlural] = &nodeInformer
			tt.initFunc(t, c, tt.handler)
			c.syncSensuHandler(tt.handler)
			if err := tt.testFunc(c, tt.handler); err != nil {
				t.Errorf("Controller.syncSensuHandler() error = %v", err)
			}
		})
	}
}
