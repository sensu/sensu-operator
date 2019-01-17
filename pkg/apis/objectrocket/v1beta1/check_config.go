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
	sensu_go_v2 "github.com/sensu/sensu-go/api/core/v2"
	k8s_api_extensions_v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuCheckConfigList is a list of CheckConfigs.
type SensuCheckConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuCheckConfig `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuCheckConfig is the k8s object associated with a sensu check
type SensuCheckConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SensuCheckConfigSpec   `json:"spec"`
	Status            SensuCheckConfigStatus `json:"status"`
}

// SensuCheckConfigSpec is the specification for a sensu check config
type SensuCheckConfigSpec struct {
	// Command is the command to be executed.
	Command string `json:"command,omitempty"`
	// Handlers are the event handler for the check (incidents and/or metrics).
	Handlers []string `json:"handlers"`
	// HighFlapThreshold is the flap detection high threshold (% state change) for
	// the check. Sensu uses the same flap detection algorithm as Nagios.
	HighFlapThreshold uint32 `json:"high_flap_threshold"`
	// Interval is the interval, in seconds, at which the check should be run.
	Interval uint32 `json:"interval"`
	// LowFlapThreshold is the flap detection low threshold (% state change) for
	// the check. Sensu uses the same flap detection algorithm as Nagios.
	LowFlapThreshold uint32 `json:"low_flap_threshold"`
	// Publish indicates if check requests are published for the check
	Publish bool `json:"publish"`
	// RuntimeAssets are a list of assets required to execute check.
	RuntimeAssets []string `json:"runtime_assets"`
	// Subscriptions is the list of subscribers for the check.
	Subscriptions []string `json:"subscriptions"`
	// ExtendedAttributes store serialized arbitrary JSON-encoded data
	ExtendedAttributes []byte `json:"-"`
	// Sources indicates the name of the entity representing an external resource
	ProxyEntityName string `json:"proxy_entity_name"`
	// CheckHooks is the list of check hooks for the check
	CheckHooks []HookList `json:"check_hooks"`
	// STDIN indicates if the check command accepts JSON via stdin from the agent
	Stdin bool `json:"stdin"`
	// Subdue represents one or more time windows when the check should be subdued.
	Subdue *TimeWindowWhen `json:"subdue"`
	// Cron is the cron string at which the check should be run.
	Cron string `json:"cron,omitempty"`
	// TTL represents the length of time in seconds for which a check result is valid.
	Ttl int64 `json:"ttl"`
	// Timeout is the timeout, in seconds, at which the check has to run
	Timeout uint32 `json:"timeout"`
	// ProxyRequests represents a request to execute a proxy check
	ProxyRequests *ProxyRequests `json:"proxy_requests,omitempty"`
	// RoundRobin enables round-robin scheduling if set true.
	RoundRobin bool `json:"round_robin"`
	// OutputOutputMetricFormat is the metric protocol that the check's output will be
	// expected to follow in order to be extracted.
	OutputMetricFormat string `json:"output_metric_format"`
	// OutputOutputMetricHandlers is the list of event handlers that will respond to metrics
	// that have been extracted from the check.
	OutputMetricHandlers []string `json:"output_metric_handlers"`
	// EnvVars is the list of environment variables to set for the check's
	// execution environment.
	EnvVars []string `json:"env_vars"`
	// Metadata contains the sensu name, sensu clusterName, sensu namespace, sensu labels and sensu annotations of the check
	SensuMetadata ObjectMeta `json:"sensuMetadata,omitempty"`
	// Validation is the OpenAPIV3Schema validation for sensu checks
	Validation k8s_api_extensions_v1beta1.CustomResourceValidation `json:"validation,omitempty"`
}

// SensuCheckConfigStatus is the status of the sensu check config
type SensuCheckConfigStatus struct {
	Accepted  bool   `json:"accepted"`
	LastError string `json:"lastError"`
}

// A ProxyRequests represents a request to execute a proxy check
type ProxyRequests struct {
	// EntityAttributes store serialized arbitrary JSON-encoded data to match
	// entities in the registry.
	EntityAttributes []string `json:"entity_attributes"`
	// Splay indicates if proxy check requests should be splayed, published evenly
	// over a window of time.
	Splay bool `json:"splay"`
	// SplayCoverage is the percentage used for proxy check request splay
	// calculation.
	SplayCoverage uint32 `json:"splay_coverage"`
}

// ToSensuType will convert an objectrocket/v1betata checkconfig into
// a sensu-go/api/core/v2 checkconfig type
func (c SensuCheckConfig) ToSensuType() *sensu_go_v2.CheckConfig {
	checkConfig := &sensu_go_v2.CheckConfig{
		Command:              c.Spec.Command,
		Handlers:             c.Spec.Handlers,
		HighFlapThreshold:    c.Spec.HighFlapThreshold,
		Interval:             c.Spec.Interval,
		LowFlapThreshold:     c.Spec.LowFlapThreshold,
		Publish:              c.Spec.Publish,
		RuntimeAssets:        c.Spec.RuntimeAssets,
		Subscriptions:        c.Spec.Subscriptions,
		ExtendedAttributes:   c.Spec.ExtendedAttributes,
		ProxyEntityName:      c.Spec.ProxyEntityName,
		Stdin:                c.Spec.Stdin,
		Cron:                 c.Spec.Cron,
		Ttl:                  c.Spec.Ttl,
		Timeout:              c.Spec.Timeout,
		RoundRobin:           c.Spec.RoundRobin,
		OutputMetricFormat:   c.Spec.OutputMetricFormat,
		OutputMetricHandlers: c.Spec.OutputMetricHandlers,
		EnvVars:              c.Spec.EnvVars,
		ObjectMeta: sensu_go_v2.ObjectMeta{
			Name:        c.Spec.SensuMetadata.Name,
			Namespace:   c.Spec.SensuMetadata.Namespace,
			Labels:      c.Spec.SensuMetadata.Labels,
			Annotations: c.Spec.SensuMetadata.Annotations,
		},
	}

	checkHookList := make([]sensu_go_v2.HookList, len(c.Spec.CheckHooks))
	for i, hookList := range c.Spec.CheckHooks {
		hList := sensu_go_v2.HookList{
			Hooks: make([]string, len(hookList.Hooks)),
		}
		for j, hook := range hookList.Hooks {
			hList.Hooks[j] = hook
		}
		checkHookList[i] = hList
	}

	checkConfig.Subdue = &sensu_go_v2.TimeWindowWhen{
		Days: sensu_go_v2.TimeWindowDays{},
	}

	if c.Spec.Subdue != nil {
		checkConfig.Subdue.Days.All = make([]*sensu_go_v2.TimeWindowTimeRange, len(c.Spec.Subdue.Days.All))
		for i, a := range c.Spec.Subdue.Days.All {
			checkConfig.Subdue.Days.All[i] = &sensu_go_v2.TimeWindowTimeRange{
				Begin: a.Begin,
				End:   a.End,
			}
		}

		checkConfig.Subdue.Days.Sunday = make([]*sensu_go_v2.TimeWindowTimeRange, len(c.Spec.Subdue.Days.Sunday))
		for i, a := range c.Spec.Subdue.Days.Sunday {
			checkConfig.Subdue.Days.Sunday[i] = &sensu_go_v2.TimeWindowTimeRange{
				Begin: a.Begin,
				End:   a.End,
			}
		}

		checkConfig.Subdue.Days.Monday = make([]*sensu_go_v2.TimeWindowTimeRange, len(c.Spec.Subdue.Days.Monday))
		for i, a := range c.Spec.Subdue.Days.Monday {
			checkConfig.Subdue.Days.Monday[i] = &sensu_go_v2.TimeWindowTimeRange{
				Begin: a.Begin,
				End:   a.End,
			}
		}

		checkConfig.Subdue.Days.Tuesday = make([]*sensu_go_v2.TimeWindowTimeRange, len(c.Spec.Subdue.Days.Tuesday))
		for i, a := range c.Spec.Subdue.Days.Tuesday {
			checkConfig.Subdue.Days.Tuesday[i] = &sensu_go_v2.TimeWindowTimeRange{
				Begin: a.Begin,
				End:   a.End,
			}
		}

		checkConfig.Subdue.Days.Wednesday = make([]*sensu_go_v2.TimeWindowTimeRange, len(c.Spec.Subdue.Days.Wednesday))
		for i, a := range c.Spec.Subdue.Days.Wednesday {
			checkConfig.Subdue.Days.Wednesday[i] = &sensu_go_v2.TimeWindowTimeRange{
				Begin: a.Begin,
				End:   a.End,
			}
		}

		checkConfig.Subdue.Days.Thursday = make([]*sensu_go_v2.TimeWindowTimeRange, len(c.Spec.Subdue.Days.Thursday))
		for i, a := range c.Spec.Subdue.Days.Thursday {
			checkConfig.Subdue.Days.Thursday[i] = &sensu_go_v2.TimeWindowTimeRange{
				Begin: a.Begin,
				End:   a.End,
			}
		}

		checkConfig.Subdue.Days.Friday = make([]*sensu_go_v2.TimeWindowTimeRange, len(c.Spec.Subdue.Days.Friday))
		for i, a := range c.Spec.Subdue.Days.Friday {
			checkConfig.Subdue.Days.Friday[i] = &sensu_go_v2.TimeWindowTimeRange{
				Begin: a.Begin,
				End:   a.End,
			}
		}

		checkConfig.Subdue.Days.Saturday = make([]*sensu_go_v2.TimeWindowTimeRange, len(c.Spec.Subdue.Days.Saturday))
		for i, a := range c.Spec.Subdue.Days.Saturday {
			checkConfig.Subdue.Days.Saturday[i] = &sensu_go_v2.TimeWindowTimeRange{
				Begin: a.Begin,
				End:   a.End,
			}
		}

		checkConfig.Subdue.Days.Sunday = make([]*sensu_go_v2.TimeWindowTimeRange, len(c.Spec.Subdue.Days.Sunday))
		for i, a := range c.Spec.Subdue.Days.Sunday {
			checkConfig.Subdue.Days.Sunday[i] = &sensu_go_v2.TimeWindowTimeRange{
				Begin: a.Begin,
				End:   a.End,
			}
		}

	}

	if checkConfig.ProxyRequests != nil {
		checkConfig.ProxyRequests = &sensu_go_v2.ProxyRequests{
			Splay:         c.Spec.ProxyRequests.Splay,
			SplayCoverage: c.Spec.ProxyRequests.SplayCoverage,
		}
	}

	if c.Spec.ProxyRequests != nil {
		checkConfig.ProxyRequests.EntityAttributes = make([]string, len(c.Spec.ProxyRequests.EntityAttributes))
		for i, e := range c.Spec.ProxyRequests.EntityAttributes {
			checkConfig.ProxyRequests.EntityAttributes[i] = e
		}
	}

	return checkConfig
}

// GetCustomResourceValidation returns the asset's resource validation
func (c SensuCheckConfig) GetCustomResourceValidation() *k8s_api_extensions_v1beta1.CustomResourceValidation {
	minItems := int64(1)
	return &k8s_api_extensions_v1beta1.CustomResourceValidation{
		OpenAPIV3Schema: &k8s_api_extensions_v1beta1.JSONSchemaProps{
			Properties: map[string]k8s_api_extensions_v1beta1.JSONSchemaProps{
				"metadata": k8s_api_extensions_v1beta1.JSONSchemaProps{
					Required: []string{"finalizers", "name", "namespace"},
					Properties: map[string]k8s_api_extensions_v1beta1.JSONSchemaProps{
						"finalizers": k8s_api_extensions_v1beta1.JSONSchemaProps{
							Type:     "array",
							MinItems: &minItems,
							// This is required to be set to false, or you get error
							// 'uniqueItems cannot be set to true since the runtime complexity becomes quadratic'
							UniqueItems: false,
							// MinItems by itself doesn't seem to work.
							Required: []string{"check.finalizer.objectrocket.com"},
						},
					},
				},
			},
		},
	}
}
