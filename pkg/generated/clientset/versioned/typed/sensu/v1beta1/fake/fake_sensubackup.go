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

// FakeSensuBackups implements SensuBackupInterface
type FakeSensuBackups struct {
	Fake *FakeSensuV1beta1
	ns   string
}

var sensubackupsResource = schema.GroupVersionResource{Group: "sensu.io", Version: "v1beta1", Resource: "sensubackups"}

var sensubackupsKind = schema.GroupVersionKind{Group: "sensu.io", Version: "v1beta1", Kind: "SensuBackup"}

// Get takes name of the sensuBackup, and returns the corresponding sensuBackup object, and an error if there is any.
func (c *FakeSensuBackups) Get(name string, options v1.GetOptions) (result *v1beta1.SensuBackup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(sensubackupsResource, c.ns, name), &v1beta1.SensuBackup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SensuBackup), err
}

// List takes label and field selectors, and returns the list of SensuBackups that match those selectors.
func (c *FakeSensuBackups) List(opts v1.ListOptions) (result *v1beta1.SensuBackupList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(sensubackupsResource, sensubackupsKind, c.ns, opts), &v1beta1.SensuBackupList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.SensuBackupList{}
	for _, item := range obj.(*v1beta1.SensuBackupList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested sensuBackups.
func (c *FakeSensuBackups) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(sensubackupsResource, c.ns, opts))

}

// Create takes the representation of a sensuBackup and creates it.  Returns the server's representation of the sensuBackup, and an error, if there is any.
func (c *FakeSensuBackups) Create(sensuBackup *v1beta1.SensuBackup) (result *v1beta1.SensuBackup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(sensubackupsResource, c.ns, sensuBackup), &v1beta1.SensuBackup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SensuBackup), err
}

// Update takes the representation of a sensuBackup and updates it. Returns the server's representation of the sensuBackup, and an error, if there is any.
func (c *FakeSensuBackups) Update(sensuBackup *v1beta1.SensuBackup) (result *v1beta1.SensuBackup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(sensubackupsResource, c.ns, sensuBackup), &v1beta1.SensuBackup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SensuBackup), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSensuBackups) UpdateStatus(sensuBackup *v1beta1.SensuBackup) (*v1beta1.SensuBackup, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(sensubackupsResource, "status", c.ns, sensuBackup), &v1beta1.SensuBackup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SensuBackup), err
}

// Delete takes name of the sensuBackup and deletes it. Returns an error if one occurs.
func (c *FakeSensuBackups) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(sensubackupsResource, c.ns, name), &v1beta1.SensuBackup{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSensuBackups) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(sensubackupsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.SensuBackupList{})
	return err
}

// Patch applies the patch and returns the patched sensuBackup.
func (c *FakeSensuBackups) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.SensuBackup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(sensubackupsResource, c.ns, name, data, subresources...), &v1beta1.SensuBackup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SensuBackup), err
}
