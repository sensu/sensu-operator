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

func TestController_syncSensuAsset(t *testing.T) {
	assetInformer, checkInformer, handlerInformer := initInformers()
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
		asset    *api.SensuAsset
		initFunc func(*testing.T, *Controller, *api.SensuAsset)
		testFunc func(*Controller, *api.SensuAsset) error
	}{
		{
			"test asset status is updated when cluster is missing",
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
			&api.SensuAsset{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testAsset",
					Namespace: "sensu",
				},
				Spec: api.SensuAssetSpec{
					URL:          "fake",
					Sha512:       "fakeSHA",
					Organization: "fakeOrg",
					SensuMetadata: api.ObjectMeta{
						Name:        "testAsset",
						Namespace:   "default",
						ClusterName: "doesntExist",
					},
				},
			},
			func(t *testing.T, c *Controller, a *api.SensuAsset) {
				_, err := c.SensuCRCli.ObjectrocketV1beta1().SensuAssets("sensu").Create(a)
				if err != nil {
					t.Error(err)
				}
			},
			func(c *Controller, a *api.SensuAsset) error {
				c.syncSensuAsset(a)
				k8sAsset, err := c.SensuCRCli.ObjectrocketV1beta1().SensuAssets("sensu").Get("testAsset", metav1.GetOptions{})
				if err != nil {
					return errors.New("failed to find asset in k8s")
				}
				if k8sAsset.Status.Accepted {
					return errors.New("asset status should not be accepted=true")
				}
				if k8sAsset.Status.LastError == "" {
					return errors.New("asset status lastError should not be empty")
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
			tt.initFunc(t, c, tt.asset)
			c.syncSensuAsset(tt.asset)
			if err := tt.testFunc(c, tt.asset); err != nil {
				t.Errorf("Controller.syncSensuAsset() error = %v", err)
			}
		})
	}
}
