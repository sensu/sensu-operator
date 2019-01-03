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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	internalinterfaces "github.com/objectrocket/sensu-operator/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// SensuBackups returns a SensuBackupInformer.
	SensuBackups() SensuBackupInformer
	// SensuClusters returns a SensuClusterInformer.
	SensuClusters() SensuClusterInformer
	// SensuRestores returns a SensuRestoreInformer.
	SensuRestores() SensuRestoreInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// SensuBackups returns a SensuBackupInformer.
func (v *version) SensuBackups() SensuBackupInformer {
	return &sensuBackupInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// SensuClusters returns a SensuClusterInformer.
func (v *version) SensuClusters() SensuClusterInformer {
	return &sensuClusterInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// SensuRestores returns a SensuRestoreInformer.
func (v *version) SensuRestores() SensuRestoreInformer {
	return &sensuRestoreInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
