/*
Copyright 2018 The sensu-operator Authors

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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	time "time"

	sensu_v1beta1 "github.com/objectrocket/sensu-operator/pkg/apis/sensu/v1beta1"
	versioned "github.com/objectrocket/sensu-operator/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/objectrocket/sensu-operator/pkg/generated/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/objectrocket/sensu-operator/pkg/generated/listers/sensu/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SensuRestoreInformer provides access to a shared informer and lister for
// SensuRestores.
type SensuRestoreInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.SensuRestoreLister
}

type sensuRestoreInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSensuRestoreInformer constructs a new informer for SensuRestore type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSensuRestoreInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSensuRestoreInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSensuRestoreInformer constructs a new informer for SensuRestore type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSensuRestoreInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SensuV1beta1().SensuRestores(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SensuV1beta1().SensuRestores(namespace).Watch(options)
			},
		},
		&sensu_v1beta1.SensuRestore{},
		resyncPeriod,
		indexers,
	)
}

func (f *sensuRestoreInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSensuRestoreInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *sensuRestoreInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&sensu_v1beta1.SensuRestore{}, f.defaultInformer)
}

func (f *sensuRestoreInformer) Lister() v1beta1.SensuRestoreLister {
	return v1beta1.NewSensuRestoreLister(f.Informer().GetIndexer())
}
