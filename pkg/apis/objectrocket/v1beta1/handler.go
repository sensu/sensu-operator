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
	sensutypes "github.com/sensu/sensu-go/types"
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

// SensuHandler is the type of sensu handlers
type SensuHandler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SensuHandlerSpec   `json:"spec"`
	Status            SensuHandlerStatus `json:"status"`
}

// SensuHandlerSpec is the spec section of the custom object
type SensuHandlerSpec struct {
	Type          string        `json:"type"`
	Mutator       string        `json:"mutator"`
	Command       string        `json:"command"`
	Timeout       uint32        `json:"timeout"`
	Socket        HandlerSocket `json:"socket"`
	Handlers      []string      `json:"handlers"`
	Filters       []string      `json:"filters"`
	EnvVars       []string      `json:"envVars"`
	RuntimeAssets []string      `json:"runtimeAssets"`
	// Metadata contains the name, namespace, and labels of the handler
	SensuMetadata ObjectMeta `json:"sensu_metadata,omitempty"`
}

// HandlerSocket is the socket description of a sensu handler.
type HandlerSocket struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
}

// SensuHandlerStatus is the status of the sensu handler
type SensuHandlerStatus struct {
	Accepted bool `json:"accepted"`
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
