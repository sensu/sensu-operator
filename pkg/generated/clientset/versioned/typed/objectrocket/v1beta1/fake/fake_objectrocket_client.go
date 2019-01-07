/*
Copyright 2019 The sensu-operator Authors

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/objectrocket/sensu-operator/pkg/generated/clientset/versioned/typed/objectrocket/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeObjectrocketV1beta1 struct {
	*testing.Fake
}

func (c *FakeObjectrocketV1beta1) CheckConfigs(namespace string) v1beta1.CheckConfigInterface {
	return &FakeCheckConfigs{c, namespace}
}

func (c *FakeObjectrocketV1beta1) HandlerSockets(namespace string) v1beta1.HandlerSocketInterface {
	return &FakeHandlerSockets{c, namespace}
}

func (c *FakeObjectrocketV1beta1) SensuAssets(namespace string) v1beta1.SensuAssetInterface {
	return &FakeSensuAssets{c, namespace}
}

func (c *FakeObjectrocketV1beta1) SensuBackups(namespace string) v1beta1.SensuBackupInterface {
	return &FakeSensuBackups{c, namespace}
}

func (c *FakeObjectrocketV1beta1) SensuClusters(namespace string) v1beta1.SensuClusterInterface {
	return &FakeSensuClusters{c, namespace}
}

func (c *FakeObjectrocketV1beta1) SensuHandlers(namespace string) v1beta1.SensuHandlerInterface {
	return &FakeSensuHandlers{c, namespace}
}

func (c *FakeObjectrocketV1beta1) SensuRestores(namespace string) v1beta1.SensuRestoreInterface {
	return &FakeSensuRestores{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeObjectrocketV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
