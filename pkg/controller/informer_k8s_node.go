package controller

import (
	corev1 "k8s.io/api/core/v1"

	sensu_client "github.com/objectrocket/sensu-operator/pkg/sensu_client"
)

const (
	platformSensuClusterName    = "sensu"
	platformSensuNamespace      = "platform"
	platformKubernetesNamespace = "sensu"
)

func (c *Controller) onUpdateNode(newObj interface{}) {
	c.logger.Debugf("in onUpdateNode, calling syncNode")
	c.syncNode(newObj.(*corev1.Node))
}

func (c *Controller) onDeleteNode(nodeName string) {
	for clusterName := range c.clusters {
		c.logger.Debugf("in onDeleteNode, attempting to see if cluster %s exists", clusterName)
		if c.clusterExists(clusterName) {
			c.logger.Debugf("in onDeleteNode, cluster %s exists", clusterName)
			c.logger.Debugf("getting client for cluster %s, k8s namespace %s, sensu namespace %s", clusterName, platformKubernetesNamespace, platformSensuNamespace)
			// Being that the cluster nodes are *not* a CR, and don't have the metadata for the sensu namespace, and k8s namespace,
			// we can only inspect the entities in the platform sensu namespace
			sensuClient := sensu_client.New(clusterName, platformKubernetesNamespace, platformSensuNamespace)
			c.logger.Debugf("calling sensuClient.DeleteNode")
			err := sensuClient.DeleteNode(nodeName)
			if err != nil {
				c.logger.Warningf("failed to handle node delete event: %v", err)
				return
			}
		}
	}
	c.logger.Debugf("in onDeleteNode, end of func")
}

func (c *Controller) syncNode(*corev1.Node) {
	c.logger.Debugf("in syncNode, doing nothing")
}
