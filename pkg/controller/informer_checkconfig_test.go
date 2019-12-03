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

func TestController_syncSensuCheckConfig(t *testing.T) {
	assetInformer, checkInformer, handlerInformer, eventFilterInformer, nodeInformer := initInformers()
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
		check    *api.SensuCheckConfig
		initFunc func(*testing.T, *Controller, *api.SensuCheckConfig)
		testFunc func(*Controller, *api.SensuCheckConfig) error
	}{
		{
			"test checkconfig status is updated when cluster is missing",
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
			&api.SensuCheckConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testCheckConfig",
					Namespace: "sensu",
				},
				Spec: api.SensuCheckConfigSpec{
					SensuMetadata: api.ObjectMeta{
						Name:        "testCheckConfig",
						Namespace:   "default",
						ClusterName: "doesntExist",
					},
				},
			},
			func(t *testing.T, c *Controller, check *api.SensuCheckConfig) {
				_, err := c.SensuCRCli.ObjectrocketV1beta1().SensuCheckConfigs("sensu").Create(check)
				if err != nil {
					t.Error(err)
				}
			},
			func(c *Controller, check *api.SensuCheckConfig) error {
				c.syncSensuCheckConfig(check)
				k8sCheck, err := c.SensuCRCli.ObjectrocketV1beta1().SensuCheckConfigs("sensu").Get("testCheckConfig", metav1.GetOptions{})
				if err != nil {
					return errors.New("failed to find checkconfig in k8s")
				}
				if k8sCheck.Status.Accepted {
					return errors.New("check status should not be accepted=true")
				}
				if k8sCheck.Status.LastError == "" {
					return errors.New("check status lastError should not be empty")
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
			c.informers[api.SensuEventFilterResourcePlural] = &eventFilterInformer
			c.informers[CoreV1NodesPlural] = &nodeInformer
			tt.initFunc(t, c, tt.check)
			c.syncSensuCheckConfig(tt.check)
			if err := tt.testFunc(c, tt.check); err != nil {
				t.Errorf("Controller.syncSensuCheck() error = %v", err)
			}
		})
	}
}
