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
package internalversion

import (
	navigator "github.com/jetstack/navigator/pkg/apis/navigator"
	scheme "github.com/jetstack/navigator/pkg/client/clientset/internalversion/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// LeaderLocksGetter has a method to return a LeaderLockInterface.
// A group's client should implement this interface.
type LeaderLocksGetter interface {
	LeaderLocks(namespace string) LeaderLockInterface
}

// LeaderLockInterface has methods to work with LeaderLock resources.
type LeaderLockInterface interface {
	Create(*navigator.LeaderLock) (*navigator.LeaderLock, error)
	Update(*navigator.LeaderLock) (*navigator.LeaderLock, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*navigator.LeaderLock, error)
	List(opts v1.ListOptions) (*navigator.LeaderLockList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *navigator.LeaderLock, err error)
	LeaderLockExpansion
}

// leaderLocks implements LeaderLockInterface
type leaderLocks struct {
	client rest.Interface
	ns     string
}

// newLeaderLocks returns a LeaderLocks
func newLeaderLocks(c *NavigatorClient, namespace string) *leaderLocks {
	return &leaderLocks{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the leaderLock, and returns the corresponding leaderLock object, and an error if there is any.
func (c *leaderLocks) Get(name string, options v1.GetOptions) (result *navigator.LeaderLock, err error) {
	result = &navigator.LeaderLock{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("leaderlocks").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of LeaderLocks that match those selectors.
func (c *leaderLocks) List(opts v1.ListOptions) (result *navigator.LeaderLockList, err error) {
	result = &navigator.LeaderLockList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("leaderlocks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested leaderLocks.
func (c *leaderLocks) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("leaderlocks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a leaderLock and creates it.  Returns the server's representation of the leaderLock, and an error, if there is any.
func (c *leaderLocks) Create(leaderLock *navigator.LeaderLock) (result *navigator.LeaderLock, err error) {
	result = &navigator.LeaderLock{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("leaderlocks").
		Body(leaderLock).
		Do().
		Into(result)
	return
}

// Update takes the representation of a leaderLock and updates it. Returns the server's representation of the leaderLock, and an error, if there is any.
func (c *leaderLocks) Update(leaderLock *navigator.LeaderLock) (result *navigator.LeaderLock, err error) {
	result = &navigator.LeaderLock{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("leaderlocks").
		Name(leaderLock.Name).
		Body(leaderLock).
		Do().
		Into(result)
	return
}

// Delete takes name of the leaderLock and deletes it. Returns an error if one occurs.
func (c *leaderLocks) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("leaderlocks").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *leaderLocks) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("leaderlocks").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched leaderLock.
func (c *leaderLocks) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *navigator.LeaderLock, err error) {
	result = &navigator.LeaderLock{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("leaderlocks").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
