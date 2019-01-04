/*
Copyright 2019 The sensu-operator Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// SensuAssetLister helps list SensuAssets.
type SensuAssetLister interface {
	// List lists all SensuAssets in the indexer.
	List(selector labels.Selector) (ret []*v1beta1.SensuAsset, err error)
	// SensuAssets returns an object that can list and get SensuAssets.
	SensuAssets(namespace string) SensuAssetNamespaceLister
	SensuAssetListerExpansion
}

// sensuAssetLister implements the SensuAssetLister interface.
type sensuAssetLister struct {
	indexer cache.Indexer
}

// NewSensuAssetLister returns a new SensuAssetLister.
func NewSensuAssetLister(indexer cache.Indexer) SensuAssetLister {
	return &sensuAssetLister{indexer: indexer}
}

// List lists all SensuAssets in the indexer.
func (s *sensuAssetLister) List(selector labels.Selector) (ret []*v1beta1.SensuAsset, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.SensuAsset))
	})
	return ret, err
}

// SensuAssets returns an object that can list and get SensuAssets.
func (s *sensuAssetLister) SensuAssets(namespace string) SensuAssetNamespaceLister {
	return sensuAssetNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// SensuAssetNamespaceLister helps list and get SensuAssets.
type SensuAssetNamespaceLister interface {
	// List lists all SensuAssets in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1beta1.SensuAsset, err error)
	// Get retrieves the SensuAsset from the indexer for a given namespace and name.
	Get(name string) (*v1beta1.SensuAsset, error)
	SensuAssetNamespaceListerExpansion
}

// sensuAssetNamespaceLister implements the SensuAssetNamespaceLister
// interface.
type sensuAssetNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all SensuAssets in the indexer for a given namespace.
func (s sensuAssetNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.SensuAsset, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.SensuAsset))
	})
	return ret, err
}

// Get retrieves the SensuAsset from the indexer for a given namespace and name.
func (s sensuAssetNamespaceLister) Get(name string) (*v1beta1.SensuAsset, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("sensuasset"), name)
	}
	return obj.(*v1beta1.SensuAsset), nil
}
