package client

import (
	"errors"
	"time"

	sensu_api_core_v2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-go/types"

	"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
)

type fetchHandlerResponse struct {
	handler *types.Handler
	err     error
}

// AddHandler will add a new sensu Handler to the sensu server
func (s *SensuClient) AddHandler(handler *v1beta1.SensuHandler) error {
	return s.ensureHandler(handler)
}

// UpdateHandler will add a new sensu Handler to the sensu server
func (s *SensuClient) UpdateHandler(handler *v1beta1.SensuHandler) error {
	return s.ensureHandler(handler)
}

// DeleteHandler will delete an existing Handler from the sensu server
func (s *SensuClient) DeleteHandler(handler *v1beta1.SensuHandler) error {
	if err := s.ensureCredentials(); err != nil {
		return err
	}

	c1 := make(chan error, 1)

	go func() {
		var err error
		if err = s.sensuCli.Client.DeleteHandler(handler.ToSensuType()); err != nil {
			s.logger.Errorf("failed to delete handler: %+v", err)
		}
		c1 <- err
	}()

	select {
	case err := <-c1:
		return err
	case <-time.After(s.timeout):
		s.logger.Warnf("timeout from sensu server after 10 seconds")
		return errors.New("timeout from sensu server after 10 seconds")
	}
}

func (s *SensuClient) ensureHandler(handler *v1beta1.SensuHandler) error {
	var (
		sensuHandler *types.Handler
		err          error
	)

	if err := s.ensureCredentials(); err != nil {
		return err
	}

	a1 := make(chan fetchHandlerResponse, 1)
	go func() {
		var err error

		if sensuHandler, err = s.sensuCli.Client.FetchHandler(handler.GetName()); err != nil {
			if err.Error() == errSensuClusterObjectNotFound.Error() {
				if err = s.sensuCli.Client.CreateHandler(handler.ToSensuType()); err != nil {
					s.logger.Errorf("Failed to create new handler %s: %s", handler.GetName(), err)
					a1 <- fetchHandlerResponse{handler.ToSensuType(), err}
					return
				}
			}
			s.logger.Warnf("failed to retrieve handler name %s from namespace %s, err: %+v", handler.GetName(), s.sensuCli.Config.Namespace(), err)
		}
		a1 <- fetchHandlerResponse{sensuHandler, nil}
	}()

	select {
	case response := <-a1:
		if response.err != nil {
			return response.err
		}
		sensuHandler = response.handler
	case <-time.After(s.timeout):
		s.logger.Warnf("timeout from sensu server after 10 seconds")
		return errors.New("timeout from sensu server after 10 seconds")
	}

	// Check to see if Handler needs updated?
	if !handlerEqual(sensuHandler, handler.ToSensuType()) {
		s.logger.Infof("current handler wasn't equal to new handler, so updating...")
		a2 := make(chan error, 1)
		go func() {
			if err = s.sensuCli.Client.UpdateHandler(handler.ToSensuType()); err != nil {
				s.logger.Errorf("Failed to update handler %s: %+v", handler.GetName(), err)
				a2 <- err
			}
			a2 <- nil
		}()

		select {
		case err = <-a2:
			return err
		case <-time.After(s.timeout):
			s.logger.Warnf("timeout from sensu server after 10 seconds")
			return errors.New("timeout from sensu server after 10 seconds")
		}
	}

	return nil
}

func handlerEqual(a1, a2 *sensu_api_core_v2.Handler) bool {
	if a1 == nil || a2 == nil {
		return false
	}

	if a1.Type != a2.Type ||
		a1.Mutator != a2.Mutator ||
		a1.Command != a2.Command ||
		a1.Timeout != a2.Timeout {
		return false
	}

	for i, handler := range a1.Handlers {
		if handler != a2.Handlers[i] {
			return false
		}
	}

	for i, filter := range a1.Filters {
		if filter != a2.Filters[i] {
			return false
		}
	}

	for i, envVar := range a1.EnvVars {
		if envVar != a2.EnvVars[i] {
			return false
		}
	}

	return true
}
