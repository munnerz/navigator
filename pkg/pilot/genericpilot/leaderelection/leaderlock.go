package leaderelection

import (
	"encoding/json"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	navigator "github.com/jetstack/navigator/pkg/apis/navigator/v1alpha1"
	navigatorclient "github.com/jetstack/navigator/pkg/client/clientset/versioned/typed/navigator/v1alpha1"
)

type LeaderLock struct {
	// LeaderLockMeta should contain a Name and a Namespace of a
	// LeaderLock object that the LeaderElector will attempt to lead.
	LeaderLockMeta metav1.ObjectMeta
	Client         navigatorclient.LeaderLocksGetter
	LockConfig     resourcelock.ResourceLockConfig
	lock           *navigator.LeaderLock
}

// Get returns the cmlection record from a ConfigMap Annotation
func (cml *LeaderLock) Get() (*resourcelock.LeaderElectionRecord, error) {
	var record resourcelock.LeaderElectionRecord
	var err error
	cml.lock, err = cml.Client.LeaderLocks(cml.LeaderLockMeta.Namespace).Get(cml.LeaderLockMeta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if cml.lock.Annotations == nil {
		cml.lock.Annotations = make(map[string]string)
	}
	if recordBytes, found := cml.lock.Annotations[resourcelock.LeaderElectionRecordAnnotationKey]; found {
		if err := json.Unmarshal([]byte(recordBytes), &record); err != nil {
			return nil, err
		}
	}
	return &record, nil
}

// Create attempts to create a LeadercmlectionRecord annotation
func (cml *LeaderLock) Create(ler resourcelock.LeaderElectionRecord) error {
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return err
	}
	cml.lock, err = cml.Client.LeaderLocks(cml.LeaderLockMeta.Namespace).Create(&navigator.LeaderLock{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cml.LeaderLockMeta.Name,
			Namespace: cml.LeaderLockMeta.Namespace,
			Annotations: map[string]string{
				resourcelock.LeaderElectionRecordAnnotationKey: string(recordBytes),
			},
		},
	})
	return err
}

// Update will update and existing annotation on a given resource.
func (cml *LeaderLock) Update(ler resourcelock.LeaderElectionRecord) error {
	if cml.lock == nil {
		return errors.New("leaderlock not initialized, call get or create first")
	}
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return err
	}
	cml.lock.Annotations[resourcelock.LeaderElectionRecordAnnotationKey] = string(recordBytes)
	cml.lock, err = cml.Client.LeaderLocks(cml.LeaderLockMeta.Namespace).Update(cml.lock)
	return err
}

// RecordEvent in leader cmlection while adding meta-data
func (cml *LeaderLock) RecordEvent(s string) {
	events := fmt.Sprintf("%v %v", cml.LockConfig.Identity, s)
	cml.LockConfig.EventRecorder.Eventf(&navigator.LeaderLock{ObjectMeta: cml.lock.ObjectMeta}, corev1.EventTypeNormal, "LeaderElection", events)
}

// Describe is used to convert details on current resource lock
// into a string
func (cml *LeaderLock) Describe() string {
	return fmt.Sprintf("%v/%v", cml.LeaderLockMeta.Namespace, cml.LeaderLockMeta.Name)
}

// returns the Identity of the lock
func (cml *LeaderLock) Identity() string {
	return cml.LockConfig.Identity
}
