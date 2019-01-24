// Copyright 2017 The etcd-operator Authors
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

package v1beta1

import (
	crdutil "github.com/objectrocket/sensu-operator/pkg/util/k8sutil/conversionutil"
	sensutypes "github.com/sensu/sensu-go/types"
	k8s_api_extensions_v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuEventFilterList is a list of sensu assets.
type SensuEventFilterList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuAsset `json:"items"`
}

// SensuEventFilter is the type of sensu assets
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type SensuEventFilter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SensuEventFilterSpec `json:"spec"`
	// Status is the sensu asset's status
	Status SensuEventFilterStatus `json:"status"`
}

// SensuEventFilterSpec is the specification for a sensu event filter
// +k8s:openapi-gen=true
type SensuEventFilterSpec struct {
	// Action specifies to allow/deny events to continue through the pipeline
	Action string `json:"action"`
	// Expressions is an array of boolean expressions that are &&'d together
	// to determine if the event matches this filter.
	Expressions []string `json:"expressions"`
	// Runtime assets are Sensu assets that contain javascript libraries. They
	// are evaluated within the execution context.
	RuntimeAssets []string `json:"runtime_assets,omitempty"`

	// Organization indicates to which org an asset belongs to
	Organization string `json:"organization,omitempty"`
	// Metadata contains the sensu name, sensu namespace, sensu labels and sensu annotations of the check
	SensuMetadata ObjectMeta `json:"sensuMetadata"`
}

// SensuEventFilterStatus is the status of the sensu event filter
type SensuEventFilterStatus struct {
	Accepted  bool   `json:"accepted"`
	LastError string `json:"lastError"`
}

// ToSensuType returns a value of the SensuAsset type from the Sensu API
func (f SensuEventFilter) ToSensuType() *sensutypes.EventFilter {
	return &sensutypes.EventFilter{
		ObjectMeta: sensutypes.ObjectMeta{
			Name:        f.ObjectMeta.Name,
			Namespace:   f.Spec.SensuMetadata.Namespace,
			Labels:      f.ObjectMeta.Labels,
			Annotations: f.ObjectMeta.Annotations,
		},
		Action:        f.Spec.Action,
		Expressions:   f.Spec.Expressions,
		RuntimeAssets: f.Spec.RuntimeAssets,
	}
}

// GetCustomResourceValidation returns the asset's resource validation
func (f SensuEventFilter) GetCustomResourceValidation() *k8s_api_extensions_v1beta1.CustomResourceValidation {
	return crdutil.GetCustomResourceValidation("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuEventHandler", GetOpenAPIDefinitions)

}
