package controller

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	"github.com/objectrocket/sensu-operator/pkg/cluster"
	fakesensu "github.com/objectrocket/sensu-operator/pkg/generated/clientset/versioned/fake"
	fakeapiextensionsapiserver "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

func TestController_syncNodeHandler(t *testing.T) {
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
		node     *corev1.Node
		initFunc func(*testing.T, *Controller, *corev1.Node)
		testFunc func(*Controller, *corev1.Node) error
	}{
		{
			"test node still exists after syncing nodes",
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
			&corev1.Node{},
			func(t *testing.T, c *Controller, node *corev1.Node) {
				_, err := c.KubeCli.CoreV1().Nodes().Create(node)
				if err != nil {
					t.Error(err)
				}
			},
			func(c *Controller, node *corev1.Node) error {
				c.syncNode(node)
				node, err := c.KubeCli.CoreV1().Nodes().Get(node.GetName(), metav1.GetOptions{})
				if err != nil {
					return errors.New("failed to find node in k8s")
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
			c.informers["nodes"] = &nodeInformer
			tt.initFunc(t, c, tt.node)
			c.syncNode(tt.node)
			if err := tt.testFunc(c, tt.node); err != nil {
				t.Errorf("Controller.syncNode() error = %v", err)
			}
		})
	}
}
