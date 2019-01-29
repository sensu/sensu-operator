// Copyright 2016 The etcd-operator Authors
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

package cluster

import (
	"errors"
	"fmt"

	"github.com/objectrocket/sensu-operator/pkg/util/k8sutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ErrLostQuorum indicates that the etcd cluster lost its quorum.
var ErrLostQuorum = errors.New("lost quorum")

// reconcile reconciles cluster current state to desired state specified by spec.
// - it tries to reconcile the cluster to desired size.
// - if the cluster needs for upgrade, it tries to upgrade old member one by one.
func (c *Cluster) reconcile(pods []*v1.Pod) error {
	if c.statefulSet.Spec.Replicas == nil {
		c.logger.Infof("StatefulSet for cluster %s has nil Replicas.  Fetching new StatefulSet", c.name())
		set, err := c.config.KubeCli.AppsV1beta1().StatefulSets(c.cluster.Namespace).Get(c.cluster.Name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("Failed to fetch new StatefulSet: %v", err)
		}
		c.statefulSet = set
		// Return here since reconcile will be called again and we can do this check again if necessary
		return nil
	}
	if c.cluster.Spec.Size != int(*c.statefulSet.Spec.Replicas) {
		set, err := c.config.KubeCli.AppsV1beta1().StatefulSets(c.cluster.Namespace).Get(c.cluster.Name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("Error getting StatefulSet %s for size update: %v", c.statefulSet.GetName(), err)
		}
		set, err = c.config.KubeCli.AppsV1beta1().StatefulSets(c.cluster.Namespace).Update(set)
		if err != nil {
			return fmt.Errorf("Error updating StatefulSet %s size: %v", c.statefulSet.GetName(), err)
		}
		c.statefulSet = set
		c.logger.Infof("Update StatefulSet %s size from %d to %d", c.statefulSet.GetName(), *c.statefulSet.Spec.Replicas, c.cluster.Spec.Size)
		return nil
	}
	var oldPod *v1.Pod
	oldPod = pickOneOldMember(pods, c.cluster.Spec.Version)
	if oldPod != nil {
		// This needs to be handled once the etcd cluster is either external or has multiple nodes
		c.logger.Warnf("Pod %s needs upgraded from version %s to %s", oldPod.GetName(), k8sutil.GetSensuVersion(oldPod), c.cluster.Spec.Version)
		return nil
	}
	return nil
}

func pickOneOldMember(pods []*v1.Pod, newVersion string) *v1.Pod {
	for _, pod := range pods {
		if k8sutil.GetSensuVersion(pod) == newVersion {
			continue
		}
		return pod
	}
	return nil
}
