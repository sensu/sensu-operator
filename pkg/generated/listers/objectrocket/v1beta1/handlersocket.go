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

// HandlerSocketLister helps list HandlerSockets.
type HandlerSocketLister interface {
	// List lists all HandlerSockets in the indexer.
	List(selector labels.Selector) (ret []*v1beta1.HandlerSocket, err error)
	// HandlerSockets returns an object that can list and get HandlerSockets.
	HandlerSockets(namespace string) HandlerSocketNamespaceLister
	HandlerSocketListerExpansion
}

// handlerSocketLister implements the HandlerSocketLister interface.
type handlerSocketLister struct {
	indexer cache.Indexer
}

// NewHandlerSocketLister returns a new HandlerSocketLister.
func NewHandlerSocketLister(indexer cache.Indexer) HandlerSocketLister {
	return &handlerSocketLister{indexer: indexer}
}

// List lists all HandlerSockets in the indexer.
func (s *handlerSocketLister) List(selector labels.Selector) (ret []*v1beta1.HandlerSocket, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.HandlerSocket))
	})
	return ret, err
}

// HandlerSockets returns an object that can list and get HandlerSockets.
func (s *handlerSocketLister) HandlerSockets(namespace string) HandlerSocketNamespaceLister {
	return handlerSocketNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// HandlerSocketNamespaceLister helps list and get HandlerSockets.
type HandlerSocketNamespaceLister interface {
	// List lists all HandlerSockets in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1beta1.HandlerSocket, err error)
	// Get retrieves the HandlerSocket from the indexer for a given namespace and name.
	Get(name string) (*v1beta1.HandlerSocket, error)
	HandlerSocketNamespaceListerExpansion
}

// handlerSocketNamespaceLister implements the HandlerSocketNamespaceLister
// interface.
type handlerSocketNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all HandlerSockets in the indexer for a given namespace.
func (s handlerSocketNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.HandlerSocket, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.HandlerSocket))
	})
	return ret, err
}

// Get retrieves the HandlerSocket from the indexer for a given namespace and name.
func (s handlerSocketNamespaceLister) Get(name string) (*v1beta1.HandlerSocket, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("handlersocket"), name)
	}
	return obj.(*v1beta1.HandlerSocket), nil
}
