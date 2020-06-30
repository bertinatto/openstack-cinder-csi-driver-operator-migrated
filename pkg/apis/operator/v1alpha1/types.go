package v1alpha1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenStackCinderDriver is a specification for a OpenStackCinderDriver resource
type OpenStackCinderDriver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackCinderDriverSpec   `json:"spec"`
	Status OpenStackCinderDriverStatus `json:"status"`
}

// OpenStackCinderDriverSpec is the spec for a OpenStackCinderDriver resource
type OpenStackCinderDriverSpec struct {
	operatorv1.OperatorSpec `json:",inline"`
}

// OpenStackCinderDriverStatus is the status for a OpenStackCinderDriver resource
type OpenStackCinderDriverStatus struct {
	operatorv1.OperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenStackCinderDriverList is a list of OpenStackCinderDriver resources
type OpenStackCinderDriverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []OpenStackCinderDriver `json:"items"`
}
