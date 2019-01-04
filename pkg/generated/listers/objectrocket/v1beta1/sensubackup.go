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

// SensuBackupLister helps list SensuBackups.
type SensuBackupLister interface {
	// List lists all SensuBackups in the indexer.
	List(selector labels.Selector) (ret []*v1beta1.SensuBackup, err error)
	// SensuBackups returns an object that can list and get SensuBackups.
	SensuBackups(namespace string) SensuBackupNamespaceLister
	SensuBackupListerExpansion
}

// sensuBackupLister implements the SensuBackupLister interface.
type sensuBackupLister struct {
	indexer cache.Indexer
}

// NewSensuBackupLister returns a new SensuBackupLister.
func NewSensuBackupLister(indexer cache.Indexer) SensuBackupLister {
	return &sensuBackupLister{indexer: indexer}
}

// List lists all SensuBackups in the indexer.
func (s *sensuBackupLister) List(selector labels.Selector) (ret []*v1beta1.SensuBackup, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.SensuBackup))
	})
	return ret, err
}

// SensuBackups returns an object that can list and get SensuBackups.
func (s *sensuBackupLister) SensuBackups(namespace string) SensuBackupNamespaceLister {
	return sensuBackupNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// SensuBackupNamespaceLister helps list and get SensuBackups.
type SensuBackupNamespaceLister interface {
	// List lists all SensuBackups in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1beta1.SensuBackup, err error)
	// Get retrieves the SensuBackup from the indexer for a given namespace and name.
	Get(name string) (*v1beta1.SensuBackup, error)
	SensuBackupNamespaceListerExpansion
}

// sensuBackupNamespaceLister implements the SensuBackupNamespaceLister
// interface.
type sensuBackupNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all SensuBackups in the indexer for a given namespace.
func (s sensuBackupNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.SensuBackup, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.SensuBackup))
	})
	return ret, err
}

// Get retrieves the SensuBackup from the indexer for a given namespace and name.
func (s sensuBackupNamespaceLister) Get(name string) (*v1beta1.SensuBackup, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("sensubackup"), name)
	}
	return obj.(*v1beta1.SensuBackup), nil
}