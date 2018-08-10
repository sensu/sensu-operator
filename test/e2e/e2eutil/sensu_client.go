/*
Copyright 2018 The sensu-operator Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package e2eutil

import (
	"github.com/sensu/sensu-go/cli/client"
	"github.com/sensu/sensu-go/cli/client/config/basic"
)

func NewSensuClient(apiURL string) (*client.RestClient, error) {
	sensuClientConfig := &basic.Config{
		Cluster: basic.Cluster{
			APIUrl: apiURL,
		},
		Profile: basic.Profile{
			Environment:  "default",
			Format:       "none",
			Organization: "default",
		},
	}
	sensuClient := client.New(sensuClientConfig)
	tokens, _, err := sensuClient.CreateAccessToken(apiURL, "admin", "P@ssw0rd!")
	if err != nil {
		return nil, err
	}

	if sensuClientConfig.SaveTokens(tokens); err != nil {
		return nil, err
	}

	return sensuClient, nil
}
