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

// SensuHandlerList is a list of sensu handlers.
type SensuHandlerList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuHandler `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// SensuHandler is the type of sensu handlers
type SensuHandler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SensuHandlerSpec   `json:"spec"`
	Status            SensuHandlerStatus `json:"status"`
}

// SensuHandlerSpec is the spec section of the custom object
// +k8s:openapi-gen=true
type SensuHandlerSpec struct {
	Type          string        `json:"type"`
	Mutator       string        `json:"mutator,omitempty"`
	Command       string        `json:"command,omitempty"`
	Timeout       uint32        `json:"timeout,omitempty"`
	Socket        HandlerSocket `json:"socket,omitempty"`
	Handlers      []string      `json:"handlers,omitempty"`
	Filters       []string      `json:"filters,omitempty"`
	EnvVars       []string      `json:"envVars,omitempty"`
	RuntimeAssets []string      `json:"runtimeAssets,omitempty"`
	// Metadata contains the sensu name, sensu namespace, sensu annotations, and sensu labels of the handler
	SensuMetadata ObjectMeta `json:"sensuMetadata"`
	// Validation is the OpenAPIV3Schema validation for sensu assets
	Validation k8s_api_extensions_v1beta1.CustomResourceValidation `json:"validation,omitempty"`
}

// HandlerSocket is the socket description of a sensu handler.
type HandlerSocket struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
}

// SensuHandlerStatus is the status of the sensu handler
type SensuHandlerStatus struct {
	Accepted  bool   `json:"accepted"`
	LastError string `json:"lastError"`
}

// ToSensuType returns a value of the Handler type from the Sensu API
func (a SensuHandler) ToSensuType() *sensutypes.Handler {
	return &sensutypes.Handler{
		ObjectMeta: sensutypes.ObjectMeta{
			Name:        a.ObjectMeta.Name,
			Namespace:   a.Spec.SensuMetadata.Namespace,
			Labels:      a.ObjectMeta.Labels,
			Annotations: a.ObjectMeta.Annotations,
		},
		Type:     a.Spec.Type,
		Mutator:  a.Spec.Mutator,
		Command:  a.Spec.Command,
		Timeout:  a.Spec.Timeout,
		Handlers: a.Spec.Handlers,
		Socket: &sensutypes.HandlerSocket{
			Host: a.Spec.Socket.Host,
			Port: a.Spec.Socket.Port,
		},
		Filters:       a.Spec.Filters,
		EnvVars:       a.Spec.EnvVars,
		RuntimeAssets: a.Spec.RuntimeAssets,
	}
}

// GetCustomResourceValidation rreturns the handlers's resource validation
func (a SensuHandler) GetCustomResourceValidation() *k8s_api_extensions_v1beta1.CustomResourceValidation {
	return crdutil.GetCustomResourceValidation("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuHandler", GetOpenAPIDefinitions)
}
