package client

import (
	"errors"
	"time"

	sensu_api_core_v2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-go/types"

	"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
)

type fetchEventFilterResponse struct {
	filter *types.EventFilter
	err    error
}

// AddEventFilter will add a new sensu EventFilter to the sensu server
func (s *SensuClient) AddEventFilter(filter *v1beta1.SensuEventFilter) error {
	return s.ensureEventFilter(filter)
}

// UpdateEventFilter will add a new sensu EventFilter to the sensu server
func (s *SensuClient) UpdateEventFilter(filter *v1beta1.SensuEventFilter) error {
	return s.ensureEventFilter(filter)
}

// DeleteEventFilter will delete an existing EventFilter from the sensu server
func (s *SensuClient) DeleteEventFilter(filter *v1beta1.SensuEventFilter) error {
	if err := s.ensureCredentials(); err != nil {
		return err
	}

	c1 := make(chan error, 1)

	go func() {
		var err error
		if err = s.sensuCli.Client.DeleteFilter(filter.Spec.SensuMetadata.Namespace, filter.Spec.SensuMetadata.Name); err != nil {
			s.logger.Errorf("failed to delete filter: %+v", err)
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

func (s *SensuClient) ensureEventFilter(filter *v1beta1.SensuEventFilter) error {
	var (
		sensuEventFilter *types.EventFilter
		err              error
	)

	if err := s.ensureCredentials(); err != nil {
		return err
	}

	if err := s.ensureNamespace(filter.Spec.SensuMetadata.Namespace); err != nil {
		return err
	}

	a1 := make(chan fetchEventFilterResponse, 1)
	go func() {
		var err error

		if sensuEventFilter, err = s.sensuCli.Client.FetchFilter(filter.GetName()); err != nil {
			if err.Error() == errSensuClusterObjectNotFound.Error() {
				if err = s.sensuCli.Client.CreateFilter(filter.ToSensuType()); err != nil {
					s.logger.Errorf("Failed to create new filter %s: %s", filter.GetName(), err)
					a1 <- fetchEventFilterResponse{filter.ToSensuType(), err}
					return
				}
			}
			s.logger.Warnf("failed to retrieve filter name %s from namespace %s, err: %+v", filter.GetName(), s.sensuCli.Config.Namespace(), err)
		}
		a1 <- fetchEventFilterResponse{sensuEventFilter, nil}
	}()

	select {
	case response := <-a1:
		if response.err != nil {
			return response.err
		}
		sensuEventFilter = response.filter
	case <-time.After(s.timeout):
		s.logger.Warnf("timeout from sensu server after 10 seconds")
		return errors.New("timeout from sensu server after 10 seconds")
	}

	// Check to see if EventFilter needs updated?
	if !filterEqual(sensuEventFilter, filter.ToSensuType()) {
		s.logger.Infof("current filter wasn't equal to new filter, so updating...")
		a2 := make(chan error, 1)
		go func() {
			if err = s.sensuCli.Client.UpdateFilter(filter.ToSensuType()); err != nil {
				s.logger.Errorf("Failed to update filter %s: %+v", filter.GetName(), err)
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

func filterEqual(f1, f2 *sensu_api_core_v2.EventFilter) bool {
	if f1 == nil || f2 == nil {
		return false
	}

	if f1.Action != f2.Action {
		return false
	}

	if len(f1.Expressions) != len(f2.Expressions) {
		return false
	}

	if len(f1.RuntimeAssets) != len(f2.RuntimeAssets) {
		return false
	}

	for i, expression := range f1.Expressions {
		if expression != f2.Expressions[i] {
			return false
		}
	}

	for i, assets := range f1.RuntimeAssets {
		if assets != f2.RuntimeAssets[i] {
			return false
		}
	}

	return true
}
