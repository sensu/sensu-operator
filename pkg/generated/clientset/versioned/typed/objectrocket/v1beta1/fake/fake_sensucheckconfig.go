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
	v1beta1 "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSensuCheckConfigs implements SensuCheckConfigInterface
type FakeSensuCheckConfigs struct {
	Fake *FakeObjectrocketV1beta1
	ns   string
}

var sensucheckconfigsResource = schema.GroupVersionResource{Group: "objectrocket.com", Version: "v1beta1", Resource: "sensucheckconfigs"}

var sensucheckconfigsKind = schema.GroupVersionKind{Group: "objectrocket.com", Version: "v1beta1", Kind: "SensuCheckConfig"}

// Get takes name of the sensuCheckConfig, and returns the corresponding sensuCheckConfig object, and an error if there is any.
func (c *FakeSensuCheckConfigs) Get(name string, options v1.GetOptions) (result *v1beta1.SensuCheckConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(sensucheckconfigsResource, c.ns, name), &v1beta1.SensuCheckConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SensuCheckConfig), err
}

// List takes label and field selectors, and returns the list of SensuCheckConfigs that match those selectors.
func (c *FakeSensuCheckConfigs) List(opts v1.ListOptions) (result *v1beta1.SensuCheckConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(sensucheckconfigsResource, sensucheckconfigsKind, c.ns, opts), &v1beta1.SensuCheckConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.SensuCheckConfigList{ListMeta: obj.(*v1beta1.SensuCheckConfigList).ListMeta}
	for _, item := range obj.(*v1beta1.SensuCheckConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested sensuCheckConfigs.
func (c *FakeSensuCheckConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(sensucheckconfigsResource, c.ns, opts))

}

// Create takes the representation of a sensuCheckConfig and creates it.  Returns the server's representation of the sensuCheckConfig, and an error, if there is any.
func (c *FakeSensuCheckConfigs) Create(sensuCheckConfig *v1beta1.SensuCheckConfig) (result *v1beta1.SensuCheckConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(sensucheckconfigsResource, c.ns, sensuCheckConfig), &v1beta1.SensuCheckConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SensuCheckConfig), err
}

// Update takes the representation of a sensuCheckConfig and updates it. Returns the server's representation of the sensuCheckConfig, and an error, if there is any.
func (c *FakeSensuCheckConfigs) Update(sensuCheckConfig *v1beta1.SensuCheckConfig) (result *v1beta1.SensuCheckConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(sensucheckconfigsResource, c.ns, sensuCheckConfig), &v1beta1.SensuCheckConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SensuCheckConfig), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSensuCheckConfigs) UpdateStatus(sensuCheckConfig *v1beta1.SensuCheckConfig) (*v1beta1.SensuCheckConfig, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(sensucheckconfigsResource, "status", c.ns, sensuCheckConfig), &v1beta1.SensuCheckConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SensuCheckConfig), err
}

// Delete takes name of the sensuCheckConfig and deletes it. Returns an error if one occurs.
func (c *FakeSensuCheckConfigs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(sensucheckconfigsResource, c.ns, name), &v1beta1.SensuCheckConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSensuCheckConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(sensucheckconfigsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.SensuCheckConfigList{})
	return err
}

// Patch applies the patch and returns the patched sensuCheckConfig.
func (c *FakeSensuCheckConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.SensuCheckConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(sensucheckconfigsResource, c.ns, name, data, subresources...), &v1beta1.SensuCheckConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SensuCheckConfig), err
}
