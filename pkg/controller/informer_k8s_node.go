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
	c.logger.Debugf("in onDeleteNode, attempting to see if cluster %s exists", platformSensuClusterName)
	if c.clusterExists(platformSensuClusterName) {
		c.logger.Debugf("in onDeleteNode, cluster %s exists", platformSensuClusterName)
		c.logger.Debugf("getting client for cluster %s, k8s namespace %s, sensu namespace %s", platformSensuClusterName, platformKubernetesNamespace, platformSensuNamespace)
		sensuClient := sensu_client.New(platformSensuClusterName, platformKubernetesNamespace, platformSensuNamespace)
		c.logger.Debugf("calling sensuClient.DeleteNode")
		err := sensuClient.DeleteNode(nodeName)
		if err != nil {
			c.logger.Warningf("failed to handle node delete event: %v", err)
			return
		}
	}
	c.logger.Debugf("in onDeleteNode, end of func")
}

func (c *Controller) syncNode(*corev1.Node) {
	c.logger.Debugf("in syncNode, doing nothing")
}
