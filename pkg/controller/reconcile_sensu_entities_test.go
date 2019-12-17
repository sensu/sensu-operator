package controller

import (
	"github.com/google/go-cmp/cmp"
	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	"github.com/objectrocket/sensu-operator/pkg/cluster"
	fakesensu "github.com/objectrocket/sensu-operator/pkg/generated/clientset/versioned/fake"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	fakeapiextensionsapiserver "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
	"testing"
)

func Test_calculateSensuEntitiesForRemoval(t *testing.T) {
	assetInformer, checkInformer, handlerInformer, eventFilterInformer, nodeInformer := initInformers()
	type fields struct {
		logger     *logrus.Entry
		Config     Config
		informers  map[string]*Informer
		finalizers map[string]cache.Indexer
		clusters   map[string]*cluster.Cluster
		nodes      []string
		entities   []string
	}
	tests := []struct {
		name   string
		fields fields
		result []string
		want   bool
	}{
		{
			"test calculateSensuEntitiesForRemoval, returns entities not seen as k8s nodes",
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
				[]string{"Node1", "Node2"},
				[]string{"Entity1", "Entity2", "Node1", "Node2"},
			},
			[]string{"Entity1", "Entity2"},
			true,
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
			response := c.calculateSensuEntitiesForRemoval(tt.fields.entities, tt.fields.nodes)
			t.Logf("response: %v", response)
			if got := cmp.Equal(tt.result, response); got != tt.want {
				t.Errorf("equal() = calculateSensuEntitiesForRemoval equality issue, wanted %t, got %t: diff %s", tt.want, got, cmp.Diff(tt.result, response))
			}
		})
	}
}

func Test_getK8sNodes(t *testing.T) {
	assetInformer, checkInformer, handlerInformer, eventFilterInformer, nodeInformer := initInformers()
	type fields struct {
		logger     *logrus.Entry
		Config     Config
		informers  map[string]*Informer
		finalizers map[string]cache.Indexer
		clusters   map[string]*cluster.Cluster
	}
	tests := []struct {
		name   string
		fields fields
		result []string
		nodes  *corev1.Node
		want   bool
	}{
		{
			"test getK8sNodes, returns nodes from k8s as an slice of strings with only the name",
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
			[]string{"testinstance.localhost"},
			&corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testinstance.localhost",
				},
			},
			true,
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
			c.Config.KubeCli.CoreV1().Nodes().Create(tt.nodes)
			response, _ := c.getK8sNodes()
			t.Logf("response: %v", response)
			if got := cmp.Equal(response, tt.result); got != tt.want {
				t.Errorf("equal() = getK8sNodes equality issue, wanted %t, got %t: diff %s", tt.want, got, cmp.Diff(response, tt.result))
			}
		})
	}
}
