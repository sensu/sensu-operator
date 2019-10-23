package client

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/sensu/sensu-go/cli/client"
	"github.com/sensu/sensu-go/types"
)

const (
	platformSensuNamespace = "platform"
)

type fetchEntityResponse struct {
	entity *types.Entity
	err    error
}

// AddNode will do nothing on a k8s node being added/updated/reconciled, for now
func (s *SensuClient) AddNode(node *corev1.Node) error {
	return s.ensureNode(node)
}

// UpdateNode will do nothing on a k8s node being added/updated/reconciled, for now
func (s *SensuClient) UpdateNode(node *corev1.Node) error {
	return s.ensureNode(node)
}

// DeleteNode will ensure that sensu entities associated with this k8s node are cleaned up
func (s *SensuClient) DeleteNode(nodeName string) error {
	if err := s.ensureCredentials(); err != nil {
		return errors.Wrap(err, "failed to ensure credentials for sensu client")
	}
	return s.ensureDeleteNode(nodeName)
}

// ensureNode left here for future use, as we potentially want to cleanup any dangling entities
func (s *SensuClient) ensureNode(node *corev1.Node) error {
	return nil
}

func (s *SensuClient) ensureDeleteNode(nodeName string) error {
	entity, err := s.fetchEntity(nodeName)
	if err != nil {
		return errors.Wrapf(err, "failed to find entity from node name %s", nodeName)
	}
	if entity == nil {
		return errors.New(fmt.Sprintf("failed to find entity from node name %s; empty entity", nodeName))
	}
	err = s.sensuCli.Client.DeleteEntity(entity.GetNamespace(), entity.GetName())
	if err != nil {
		s.logger.Warnf("failed to delete entity %+v from namespace %s, err: %+v", entity, entity.GetNamespace(), err)
		return errors.Wrapf(err, "failed to delete entity %+v from namespace %s", entity, entity.GetNamespace())
	}
	return nil
}

func (s *SensuClient) fetchEntity(nodeName string) (*types.Entity, error) {
	var (
		entity   *types.Entity
		entities []types.Entity
		err      error
	)
	c1 := make(chan fetchEntityResponse, 1)
	go func() {
		if entities, err = s.sensuCli.Client.ListEntities(platformSensuNamespace, &client.ListOptions{
			LabelSelector: labels.FormatLabels(map[string]string{"k8s_node": nodeName}),
		}); err != nil {
			s.logger.Warnf("failed to retrieve entities from namespace %s, err: %+v", platformSensuNamespace, err)
			c1 <- fetchEntityResponse{nil, errors.Wrapf(err, "failed to retrieve entities from namespace %s", platformSensuNamespace)}
		}
		for _, entity := range entities {
			if val, ok := entity.GetLabels()["k8s_node"]; ok {
				if val == nodeName {
					s.logger.Debugf("found entity %+v", entity)
					c1 <- fetchEntityResponse{&entity, nil}
					return
				}
			}
		}
		c1 <- fetchEntityResponse{nil, nil}
	}()

	select {
	case response := <-c1:
		if response.err != nil {
			return nil, response.err
		}
		entity = response.entity
	case <-time.After(s.timeout):
		s.logger.Warnf("timeout from sensu server after 10 seconds")
		return nil, errors.New("timeout from sensu server after 10 seconds")
	}
	return entity, nil
}
