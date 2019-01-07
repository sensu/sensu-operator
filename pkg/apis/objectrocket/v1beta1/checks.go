package v1beta1

import (
	sensu_go_v2 "github.com/sensu/sensu-go/api/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/runtime"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CheckConfigList is a list of CheckConfigs.
type CheckConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CheckConfig `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CheckConfig is the k8s object associated with a sensu check
type CheckConfig struct {
	metav1.TypeMeta         `json:",inline"`
	metav1.ObjectMeta       `json:"metadata,omitempty"`
	sensu_go_v2.CheckConfig `json:"v2CheckConfig"`
	ClusterName             string `json:"clusterName"`
	SensuNamespace          string `json:"sensuNamespace"`
}

func (c *CheckConfig) DeepCopyInto(out *CheckConfig) {
	// TODO
}

func (c *CheckConfig) DeepCopy() *CheckConfig {
	out := &CheckConfig{}
	c.DeepCopyInto(out)
	return out
}
