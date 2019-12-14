package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sensu_client "github.com/objectrocket/sensu-operator/pkg/sensu_client"
)

// Function to reconcile sensu entities
func (c *Controller) reconcileSensuEntities() {
	c.logger.Debugf("At reconcileSensuEntities, starting reconcile process")
	nodes, err := c.Config.KubeCli.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return
	}
	// items := nodes.Items
	for _, item := range nodes.Items {
		c.logger.Debugf("At processReconcileItems, kubernetes node name: %v", item.Name)
	}
	for clusterName := range c.clusters {
		if c.clusterExists(clusterName) {
			c.logger.Debugf("Listing entities for sensu cluster: %v", clusterName)
			// Being that the cluster nodes are *not* a CR, and don't have the metadata for the sensu namespace, and k8s namespace,
			// we can only inspect the entities in the platform sensu namespace
			sensuClient := sensu_client.New(clusterName, platformKubernetesNamespace, platformSensuNamespace)
			entities, err := sensuClient.ListEntities(platformKubernetesNamespace)
			if err != nil {
				c.logger.Debugf("failed to list entities for sensu namespace: %v", platformKubernetesNamespace)
				return
			}
			c.logger.Debugf("Found sensu entities: %v", entities)
			return
		}
	}

}
