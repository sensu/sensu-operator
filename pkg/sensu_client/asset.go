package client

import (
	"errors"
	"time"

	sensu_api_core_v2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-go/types"

	"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
)

type fetchAssetResponse struct {
	asset *types.Asset
	err   error
}

// AddAsset will add a new sensu Asset to the sensu server
func (s *SensuClient) AddAsset(asset *v1beta1.SensuAsset) error {
	return s.ensureAsset(asset)
}

// UpdateAsset will add a new sensu Asset to the sensu server
func (s *SensuClient) UpdateAsset(asset *v1beta1.SensuAsset) error {
	return s.ensureAsset(asset)
}

// DeleteAsset will delete an existing Asset from the sensu server
func (s *SensuClient) DeleteAsset(asset *v1beta1.SensuAsset) error {
	// Delete Asset is not a thing, apparantly..
	return nil
}

func (s *SensuClient) ensureAsset(asset *v1beta1.SensuAsset) error {
	var (
		sensuAsset *types.Asset
		err        error
	)

	if err := s.ensureCredentials(); err != nil {
		return err
	}

	a1 := make(chan fetchAssetResponse, 1)
	go func() {
		var err error

		if sensuAsset, err = s.sensuCli.Client.FetchAsset(asset.GetName()); err != nil {
			if err.Error() == errSensuClusterObjectNotFound.Error() {
				if err = s.sensuCli.Client.CreateAsset(asset.ToAPISensuAsset()); err != nil {
					s.logger.Errorf("Failed to create new asset %s: %s", asset.GetName(), err)
					a1 <- fetchAssetResponse{asset.ToAPISensuAsset(), err}
					return
				}
			}
			s.logger.Warnf("failed to retrieve asset name %s from namespace %s, err: %+v", asset.GetName(), s.sensuCli.Config.Namespace(), err)
		}
		a1 <- fetchAssetResponse{sensuAsset, nil}
	}()

	select {
	case response := <-a1:
		if response.err != nil {
			return response.err
		}
		sensuAsset = response.asset
	case <-time.After(s.timeout):
		s.logger.Warnf("timeout from sensu server after 10 seconds")
		return errors.New("timeout from sensu server after 10 seconds")
	}

	// Check to see if Asset needs updated?
	if !assetEqual(sensuAsset, asset.ToAPISensuAsset()) {
		s.logger.Infof("current asset wasn't equal to new asset, so updating...")
		a2 := make(chan error, 1)
		go func() {
			if err = s.sensuCli.Client.UpdateAsset(asset.ToAPISensuAsset()); err != nil {
				s.logger.Errorf("Failed to update asset %s: %+v", asset.GetName(), err)
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

func assetEqual(a1, a2 *sensu_api_core_v2.Asset) bool {
	if a1 == nil || a2 == nil {
		return false
	}

	if a1.URL != a2.URL ||
		a1.Sha512 != a2.Sha512 {
		return false
	}

	for i, filter := range a1.Filters {
		if filter != a2.Filters[i] {
			return false
		}
	}

	return true
}
