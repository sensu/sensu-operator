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
	"errors"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultRepository = "sensu/sensu"
	// DefaultSensuVersion is the default sensu version
	DefaultSensuVersion = "5.1.0"
	// DefaultClusterAdminUsername is the default cluster admin username
	DefaultClusterAdminUsername = "admin"
	// DefaultClusterAdminPassword is the default cluster admin password
	DefaultClusterAdminPassword = "P@ssw0rd!"
)

var (
	// ErrBackupUnsetRestoreSet is an error specific to backup policy not being set
	// TODO: move validation code into separate package.
	ErrBackupUnsetRestoreSet = errors.New("spec: backup policy must be set if restore policy is set")
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuClusterList is a list of sensu clusters.
type SensuClusterList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuCluster `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuCluster is a sensu cluster
type SensuCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterSpec   `json:"spec"`
	Status            ClusterStatus `json:"status"`
}

// AsOwner returns an owner ref
func (c *SensuCluster) AsOwner() metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: SchemeGroupVersion.String(),
		Kind:       SensuClusterResourceKind,
		Name:       c.Name,
		UID:        c.UID,
		Controller: &trueVar,
	}
}

// ClusterSpec is the cluster specification details
type ClusterSpec struct {
	// Size is the expected size of the sensu cluster.
	// The sensu-operator will eventually make the size of the running
	// cluster equal to the expected size.
	// It is currently limited to 1.
	Size int `json:"size"`
	// Repository is the name of the repository that hosts
	// sensu container images. It should be direct clone of the repository in official
	// release:
	//   https://github.com/sensu/sensu-go/releases
	// That means, it should have exact same tags and the same meaning for the tags.
	//
	// By default, it is `sensuapp/sensu-go`.
	Repository string `json:"repository,omitempty"`

	// Version is the expected version of the sensu cluster.
	// The sensu-operator will eventually make the sensu cluster version
	// equal to the expected version.
	//
	// Available versions can be found here: https://hub.docker.com/r/sensuapp/sensu-go/tags/
	//
	// If version is not set, default is "2.0.0-alpha".
	Version string `json:"version,omitempty"`

	// Paused is to pause the control of the operator for the sensu cluster.
	Paused bool `json:"paused,omitempty"`

	// Pod defines the policy to create pod for the sensu pod.
	//
	// Updating Pod does not take effect on any existing sensu pods.
	Pod *PodPolicy `json:"pod,omitempty"`

	// sensu cluster TLS configuration
	TLS *TLSPolicy `json:"TLS,omitempty"`

	// sensu cluster admin username
	ClusterAdminUsername string `json:"clusteradminusername,omitempty"`

	// sensu cluster admin password
	ClusterAdminPassword string `json:"clusteradminpassword,omitempty"`
}

// PodPolicy defines the policy to create pod for the sensu container.
type PodPolicy struct {
	// Labels specifies the labels to attach to pods the operator creates for the
	// sensu cluster.
	// "app" and "sensu_*" labels are reserved for the internal use of the sensu operator.
	// Do not overwrite them.
	Labels map[string]string `json:"labels,omitempty"`

	// NodeSelector specifies a map of key-value pairs. For the pod to be eligible
	// to run on a node, the node must have each of the indicated key-value pairs as
	// labels.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// The scheduling constraints on sensu pods.
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// **DEPRECATED**. Use Affinity instead.
	AntiAffinity bool `json:"antiAffinity,omitempty"`

	// Resources is the resource requirements for the sensu container.
	// This field cannot be updated once the cluster is created.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`

	// Tolerations specifies the pod's tolerations.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// List of environment variables to set in the sensu container.
	// This is used to configure etcd in the sensu process. The sensu cluster cannot be created, when
	// bad environement variables are provided. Do not overwrite any flags used to
	// bootstrap the cluster (for example `--initial-cluster` flag).
	// This field cannot be updated.
	SensuEnv []v1.EnvVar `json:"sensuEnv,omitempty"`

	// PersistentVolumeClaimSpec is the spec to describe PVC for the sensu container
	// This field is optional. If no PVC spec, sensu container will use emptyDir as volume
	// Note. This feature is in alpha stage. It is currently only used as non-stable storage,
	// not the stable storage. Future work need to make it used as stable storage.
	PersistentVolumeClaimSpec *v1.PersistentVolumeClaimSpec `json:"persistentVolumeClaimSpec,omitempty"`

	// Annotations specifies the annotations to attach to pods the operator creates for the
	// sensu cluster.
	// The "sensu.version" annotation is reserved for the internal use of the sensu operator.
	Annotations map[string]string `json:"annotations,omitempty"`

	// busybox init container image. default is busybox:1.28.0-glibc
	// busybox:latest uses uclibc which contains a bug that sometimes prevents name resolution
	// More info: https://github.com/docker-library/busybox/issues/27
	BusyboxImage string `json:"busyboxImage,omitempty"`

	// SecurityContext specifies the security context for the entire pod
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context
	SecurityContext *v1.PodSecurityContext `json:"securityContext,omitempty"`

	// DNSTimeoutInSecond is the maximum allowed time for the init container of the sensu pod to
	// reverse DNS lookup its IP given the hostname.
	// The default is to wait indefinitely and has a vaule of 0.
	DNSTimeoutInSecond int64 `json:"DNSTimeoutInSecond,omitempty"`
}

// Validate validates the sensu cluster
// TODO: move this to initializer
func (c *ClusterSpec) Validate() error {
	if c.TLS != nil {
		if err := c.TLS.Validate(); err != nil {
			return err
		}
	}

	if c.Pod != nil {
		for k := range c.Pod.Labels {
			if k == "app" || strings.HasPrefix(k, "sensu_") {
				return errors.New("spec: pod labels contains reserved label")
			}
		}
	}
	return nil
}

// SetDefaults cleans up user passed spec, e.g. defaulting, transforming fields.
// TODO: move this to initializer
func (c *SensuCluster) SetDefaults() {
	spec := &c.Spec
	if len(spec.Repository) == 0 {
		spec.Repository = defaultRepository
	}

	if len(spec.Version) == 0 {
		spec.Version = DefaultSensuVersion
	}

	if len(spec.ClusterAdminUsername) == 0 {
		spec.ClusterAdminUsername = DefaultClusterAdminUsername
	}

	if len(spec.ClusterAdminPassword) == 0 {
		spec.ClusterAdminPassword = DefaultClusterAdminPassword
	}

	spec.Version = strings.TrimLeft(spec.Version, "v")

	// convert PodPolicy.AntiAffinity to Pod.Affinity.PodAntiAffinity
	// TODO: Remove this once PodPolicy.AntiAffinity is removed
	if spec.Pod != nil && spec.Pod.AntiAffinity && spec.Pod.Affinity == nil {
		spec.Pod.Affinity = &v1.Affinity{
			PodAntiAffinity: &v1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []v1.PodAffinityTerm{
					{
						// set anti-affinity to the sensu pods that belongs to the same cluster
						LabelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{
							"sensu_cluster": c.Name,
						}},
						TopologyKey: "kubernetes.io/hostname",
					},
				},
			},
		}
	}
}
