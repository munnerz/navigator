package navigator

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// In this file we define the outer containing types for the ElasticsearchCluster
// type. We could import these directly into message types defined in the types.proto
// file, but this is still TODO

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ElasticsearchCluster describes a specification for an Elasticsearch cluster
type ElasticsearchCluster struct {
	// we embed these types so the ElasticsearchCluster implements runtime.Object
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ElasticsearchClusterSpec
	Status ElasticsearchClusterStatus
}

type ElasticsearchClusterStatus struct {
	NodePools map[string]ElasticsearchClusterNodePoolStatus
}

type ElasticsearchClusterNodePoolStatus struct {
	ReadyReplicas int64
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ElasticsearchClusterList defines a List type for our custom ElasticsearchCluster type.
// This is needed in order to make List operations work.
type ElasticsearchClusterList struct {
	// we embed these types so that ElasticsearchClusterList implements runtime.Object and List interfaces
	metav1.TypeMeta
	metav1.ListMeta

	Items []ElasticsearchCluster
}

type ElasticsearchClusterSpec struct {
	Plugins   []string
	NodePools []ElasticsearchClusterNodePool
	Pilot     ElasticsearchPilotImage
	Image     ElasticsearchImage
	Sysctl    []string
}

type ElasticsearchClusterNodePool struct {
	Name         string
	Replicas     int64
	Roles        []ElasticsearchClusterRole
	NodeSelector map[string]string
	Resources    *v1.ResourceRequirements
	Persistence  ElasticsearchClusterPersistenceConfig
	// Config is a map of configuration files to be placed in the elasticsearch
	// config directory. Environment variables may be used in these files and
	// they will be automatically expanded by the Elasticsearch process.
	Config map[string]string
}

type ElasticsearchClusterRole string

const (
	ElasticsearchRoleData   ElasticsearchClusterRole = "data"
	ElasticsearchRoleMaster ElasticsearchClusterRole = "master"
	ElasticsearchRoleIngest ElasticsearchClusterRole = "ingest"
)

type ElasticsearchClusterPersistenceConfig struct {
	Enabled      bool
	Size         string
	StorageClass string
}

type ImageSpec struct {
	Repository string
	Tag        string
	PullPolicy string
}

type ElasticsearchPilotImage struct {
	ImageSpec
}

type ElasticsearchImage struct {
	ImageSpec
	FsGroup int64
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Pilot struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   PilotSpec
	Status PilotStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PilotList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Pilot
}

type PilotSpec struct {
	Phase         PilotPhase
	Elasticsearch *PilotElasticsearchSpec
}

type PilotPhase string

const (
	PilotPhaseStarted PilotPhase = "Started"
	PilotPhaseStopped PilotPhase = "Stopped"
)

type PilotElasticsearchSpec struct {
}

type PilotStatus struct {
	Conditions []PilotCondition
}

// PilotCondition contains condition information for a Pilot.
type PilotCondition struct {
	// Type of the condition, currently ('Ready').
	Type PilotConditionType

	// Status of the condition, one of ('True', 'False', 'Unknown').
	Status ConditionStatus

	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	LastTransitionTime metav1.Time

	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	Reason string

	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	Message string
}

// PilotConditionType represents a Pilot condition value.
type PilotConditionType string

const (
	// PilotConditionReady represents the fact that a given Pilot condition
	// is in ready state.
	PilotConditionReady PilotConditionType = "Ready"
	// PilotConditionStarted represents the fact that a given Pilot condition
	// is in started state.
	PilotConditionStarted PilotConditionType = "Started"
	// PilotConditionStopped represents the fact that a given Pilot
	// condition is in a stopped state.
	PilotConditionStopped PilotConditionType = "Stopped"
)

// ConditionStatus represents a condition's status.
type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in
// the condition; "ConditionFalse" means a resource is not in the condition;
// "ConditionUnknown" means kubernetes can't decide if a resource is in the
// condition or not. In the future, we could add other intermediate
// conditions, e.g. ConditionDegraded.
const (
	// ConditionTrue represents the fact that a given condition is true
	ConditionTrue ConditionStatus = "True"

	// ConditionFalse represents the fact that a given condition is false
	ConditionFalse ConditionStatus = "False"

	// ConditionUnknown represents the fact that a given condition is unknown
	ConditionUnknown ConditionStatus = "Unknown"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LeaderLock is used for leader election between Pilot resources
type LeaderLock struct {
	metav1.TypeMeta
	metav1.ObjectMeta
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// This is needed in order to make List operations work.
type LeaderLockList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []LeaderLock
}
