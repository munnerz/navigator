/*
Copyright 2017 Jetstack Ltd.

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
package fake

import (
	v1alpha1 "github.com/jetstack/navigator/pkg/apis/navigator/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeLeaderLocks implements LeaderLockInterface
type FakeLeaderLocks struct {
	Fake *FakeNavigatorV1alpha1
	ns   string
}

var leaderlocksResource = schema.GroupVersionResource{Group: "navigator.jetstack.io", Version: "v1alpha1", Resource: "leaderlocks"}

var leaderlocksKind = schema.GroupVersionKind{Group: "navigator.jetstack.io", Version: "v1alpha1", Kind: "LeaderLock"}

// Get takes name of the leaderLock, and returns the corresponding leaderLock object, and an error if there is any.
func (c *FakeLeaderLocks) Get(name string, options v1.GetOptions) (result *v1alpha1.LeaderLock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(leaderlocksResource, c.ns, name), &v1alpha1.LeaderLock{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LeaderLock), err
}

// List takes label and field selectors, and returns the list of LeaderLocks that match those selectors.
func (c *FakeLeaderLocks) List(opts v1.ListOptions) (result *v1alpha1.LeaderLockList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(leaderlocksResource, leaderlocksKind, c.ns, opts), &v1alpha1.LeaderLockList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.LeaderLockList{}
	for _, item := range obj.(*v1alpha1.LeaderLockList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested leaderLocks.
func (c *FakeLeaderLocks) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(leaderlocksResource, c.ns, opts))

}

// Create takes the representation of a leaderLock and creates it.  Returns the server's representation of the leaderLock, and an error, if there is any.
func (c *FakeLeaderLocks) Create(leaderLock *v1alpha1.LeaderLock) (result *v1alpha1.LeaderLock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(leaderlocksResource, c.ns, leaderLock), &v1alpha1.LeaderLock{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LeaderLock), err
}

// Update takes the representation of a leaderLock and updates it. Returns the server's representation of the leaderLock, and an error, if there is any.
func (c *FakeLeaderLocks) Update(leaderLock *v1alpha1.LeaderLock) (result *v1alpha1.LeaderLock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(leaderlocksResource, c.ns, leaderLock), &v1alpha1.LeaderLock{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LeaderLock), err
}

// Delete takes name of the leaderLock and deletes it. Returns an error if one occurs.
func (c *FakeLeaderLocks) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(leaderlocksResource, c.ns, name), &v1alpha1.LeaderLock{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeLeaderLocks) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(leaderlocksResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.LeaderLockList{})
	return err
}

// Patch applies the patch and returns the patched leaderLock.
func (c *FakeLeaderLocks) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.LeaderLock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(leaderlocksResource, c.ns, name, data, subresources...), &v1alpha1.LeaderLock{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LeaderLock), err
}
