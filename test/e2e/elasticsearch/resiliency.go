package elasticsearch

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes"

	"github.com/jetstack/navigator/internal/test/util/generate"
	"github.com/jetstack/navigator/pkg/apis/navigator/v1alpha1"
	clientset "github.com/jetstack/navigator/pkg/client/clientset/versioned"
	"github.com/jetstack/navigator/test/e2e/framework"
)

var _ = Describe("Resiliency tests", func() {
	f := framework.NewDefaultFramework("elasticsearch-resiliency")
	var ns string
	var kubeClient kubernetes.Interface
	var navClient clientset.Interface

	BeforeEach(func() {
		kubeClient = f.KubeClientset
		navClient = f.NavigatorClientset
		ns = f.Namespace.Name
	})

	framework.NavigatorDescribe("Elasticsearch resiliency tests [ElasticsearchResiliency]", func() {
		clusterName := "test"

		AfterEach(func() {
			if CurrentGinkgoTestDescription().Failed {
				framework.DumpDebugInfo(kubeClient, ns)
			}
			framework.Logf("Deleting all elasticsearchClusters in ns %v", ns)
			framework.DeleteAllElasticsearchClusters(navClient, ns)
			framework.DeleteAllStatefulSets(kubeClient, ns)
			framework.WaitForNoPodsInNamespace(kubeClient, ns, framework.NamespaceCleanupTimeout)
		})

		It("should continue to serve reads when a single node in a two node cluster fails", func() {
			nodePoolName := "mixed"
			cluster := generate.Cluster(generate.ClusterConfig{
				Name:      clusterName,
				Namespace: ns,
				Version:   "5.6.2",
				ClusterConfig: v1alpha1.NavigatorClusterConfig{
					PilotImage: framework.DefaultElasticsearchPilotImageSpec(),
					Sysctls:    framework.DefaultElasticsearchSysctls(),
				},
				NodePools: []v1alpha1.ElasticsearchClusterNodePool{
					{
						Name:      nodePoolName,
						Replicas:  2,
						Resources: framework.DefaultElasticsearchNodeResources(),
						Roles: []v1alpha1.ElasticsearchClusterRole{
							v1alpha1.ElasticsearchRoleData,
							v1alpha1.ElasticsearchRoleIngest,
							v1alpha1.ElasticsearchRoleMaster,
						},
					},
				},
			})
			tester := framework.NewElasticsearchTester(kubeClient, navClient)
			cluster = tester.CreateClusterAndWaitForReady(cluster)
			tester.WaitForHealth(cluster, v1alpha1.ElasticsearchClusterHealthGreen)
			By("Sending a GET / request to the cluster")
			out, err := tester.QueryCluster(cluster, "/", "GET", "")
			framework.Logf(out)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
