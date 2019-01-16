// Copyright 2016 The etcd-operator Authors
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

package k8sutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	"github.com/objectrocket/sensu-operator/pkg/util/retryutil"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	fakediscovery "k8s.io/client-go/discovery/fake"
)

// TODO: replace this package with Operator client

// SensuClusterCRUpdateFunc is a function to be used when atomically
// updating a Cluster CR.
type SensuClusterCRUpdateFunc func(*api.SensuCluster)

func GetClusterList(restcli rest.Interface, ns string) (*api.SensuClusterList, error) {
	b, err := restcli.Get().RequestURI(listClustersURI(ns)).DoRaw()
	if err != nil {
		return nil, err
	}

	clusters := &api.SensuClusterList{}
	if err := json.Unmarshal(b, clusters); err != nil {
		return nil, err
	}
	return clusters, nil
}

func listClustersURI(ns string) string {
	return fmt.Sprintf("/apis/%s/namespaces/%s/%s", api.SchemeGroupVersion.String(), ns, api.SensuClusterResourcePlural)
}

func CreateCRD(clientset apiextensionsclient.Interface,
	crdName,
	rkind,
	rplural,
	shortName string,
	validation *apiextensionsv1beta1.CustomResourceValidation) error {
	crd := &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: crdName,
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   api.SchemeGroupVersion.Group,
			Version: api.SchemeGroupVersion.Version,
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: rplural,
				Kind:   rkind,
			},
			Validation: validation,
		},
	}
	if len(shortName) != 0 {
		crd.Spec.Names.ShortNames = []string{shortName}
	}
	existingCRD, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && IsKubernetesResourceAlreadyExistError(err) {
		// Get the version from k8s, as the above version doesn't seem to have resourceversion
		if existingCRD, err = clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.GetName(), metav1.GetOptions{}); err != nil {
			return err
		}
		if !crdEqual(existingCRD, crd) {
			crd.ResourceVersion = existingCRD.ResourceVersion
			if _, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Update(crd); err != nil {
				return err
			}
		}
		return nil
	} else if err != nil && !IsKubernetesResourceAlreadyExistError(err) {
		return err
	}
	return nil
}

func crdEqual(crd1, crd2 *apiextensionsv1beta1.CustomResourceDefinition) (equal bool) {
	if crd1.Spec.Group != crd2.Spec.Group ||
		crd1.Spec.Version != crd2.Spec.Version ||
		crd1.Spec.Scope != crd2.Spec.Scope ||
		!reflect.DeepEqual(crd1.Spec.Names, crd2.Spec.Names) ||
		!reflect.DeepEqual(crd1.Spec.Validation, crd2.Spec.Validation) {
		return false
	}
	return true
}

func WaitCRDReady(clientset apiextensionsclient.Interface, crdName string) error {
	// If we're testing, then just assume the CRDs are ready,
	// as they will never get status conditions in testing
	//
	// TODO: is there a better way to do this - mmontgomery
	switch clientset.Discovery().(type) {
	case *fakediscovery.FakeDiscovery:
		return nil
	}

	err := retryutil.Retry(5*time.Second, 20, func() (bool, error) {
		crd, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crdName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1beta1.Established:
				if cond.Status == apiextensionsv1beta1.ConditionTrue {
					return true, nil
				}
			case apiextensionsv1beta1.NamesAccepted:
				if cond.Status == apiextensionsv1beta1.ConditionFalse {
					return false, fmt.Errorf("Name conflict: %v", cond.Reason)
				}
			}
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("wait CRD created failed: %v", err)
	}
	return nil
}

func MustNewKubeExtClient() apiextensionsclient.Interface {
	cfg, err := InClusterConfig()
	if err != nil {
		panic(err)
	}
	return apiextensionsclient.NewForConfigOrDie(cfg)
}
