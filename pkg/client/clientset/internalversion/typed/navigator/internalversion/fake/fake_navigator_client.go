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
	internalversion "github.com/jetstack/navigator/pkg/client/clientset/internalversion/typed/navigator/internalversion"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeNavigator struct {
	*testing.Fake
}

func (c *FakeNavigator) CassandraClusters(namespace string) internalversion.CassandraClusterInterface {
	return &FakeCassandraClusters{c, namespace}
}

func (c *FakeNavigator) ElasticsearchClusters(namespace string) internalversion.ElasticsearchClusterInterface {
	return &FakeElasticsearchClusters{c, namespace}
}

func (c *FakeNavigator) Pilots(namespace string) internalversion.PilotInterface {
	return &FakePilots{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeNavigator) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
