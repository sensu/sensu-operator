package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sensu_client "github.com/objectrocket/sensu-operator/pkg/sensu_client"
)

// Function to reconcile sensu entities
func (c *Controller) reconcileSensuEntities() {
	var (
		nodes    []string
		entities []string
		err      error
	)

	if nodes, err = c.getK8sNodes(); err != nil {
		c.logger.Debugf("At reconcileSensuEntities, could not get k8s nodes error: %v", err)
		return
	}
	if entities, err = c.getSensuEntities(); err != nil {
		c.logger.Debugf("At reconcileSensuEntities, could not get sensu entities error: %v", err)
		return
	}

	c.logger.Debugf("At reconcileSensuEntities, starting reconcile process")
	difference := c.calculateSensuEntitiesForRemoval(entities, nodes)
	if len(difference) > 0 {
		c.logger.Debugf("At reconcileSensuEntities, k8s nodes: %v", nodes)
		c.logger.Debugf("At reconcileSensuEntities, sensu entities: %v", entities)
		c.logger.Debugf("At reconcileSensuEntities, sensu entities for removal: %v", difference)
		c.deleteSensuEntities(difference)
	} else {
		c.logger.Debugf("At reconcileSensuEntities, finish with no differences from k8s and sensu entities")
	}

	return
}

// getK8sNodes returns a slice of node names from kubernetes
func (c *Controller) getK8sNodes() ([]string, error) {
	nodes := []string{}
	k8sNodes, err := c.Config.KubeCli.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, item := range k8sNodes.Items {
		nodes = append(nodes, item.Name)
	}
	return nodes, nil

}

// getSensuEntities returns a slice of sensu entities with matching name <-> k8s_node label
func (c *Controller) getSensuEntities() ([]string, error) {
	entities := []string{}
	for clusterName := range c.clusters {
		if c.clusterExists(clusterName) {
			c.logger.Debugf("Listing entities for sensu cluster: %v", clusterName)
			sensuClient := sensu_client.New(clusterName, platformKubernetesNamespace, platformSensuNamespace)
			sensuEntities, err := sensuClient.ListEntities(platformSensuNamespace)
			if err != nil {
				c.logger.Debugf("Failed to list entities for sensu namespace: %v", platformSensuNamespace)
				return nil, err
			}
			if len(sensuEntities) == 0 {
				c.logger.Debugf("List of entities from sensu is empty")
				return nil, nil
			}
			for _, sensuEntity := range sensuEntities {
				// We only return sensu entities with the label "k8s_node" assigned and matching its entity name. See PLAT-9090 for reference.
				if nodeName, ok := sensuEntity.GetLabels()["k8s_node"]; ok {
					if sensuEntity.GetName() == nodeName {
						c.logger.Debugf("Found Sensu entity: %v with k8s_node label: %v", sensuEntity.GetName(), nodeName)
						entities = append(entities, nodeName)
					}
				}
			}
		}
	}
	return entities, nil

}

// calculateSensuEntitiesForRemoval is used to return a slice of sensu entities not seen as kubernetes nodes,
// and not in the list of exclusions.
func (c *Controller) calculateSensuEntitiesForRemoval(entities []string, nodes []string) []string {
	var diff []string

	// find sensu entities not in kubernetes,
	for i := 0; i < 1; i++ {
		for _, s1 := range entities {
			found := false
			for _, s2 := range nodes {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
	}
	return diff
}

// deleteSensuEntities takes an slice of sensu entities and removes them from sensu
func (c *Controller) deleteSensuEntities(entities []string) {
	for clusterName := range c.clusters {
		if c.clusterExists(clusterName) {
			sensuClient := sensu_client.New(clusterName, platformKubernetesNamespace, platformSensuNamespace)
			for _, entitiy := range entities {
				c.logger.Debugf("calling sensuClient.DeleteNode for %v", entitiy)
				err := sensuClient.DeleteNode(entitiy)
				if err != nil {
					c.logger.Warningf("failed to handle node delete event: %v", err)
					return
				}
			}

		}
	}
}
