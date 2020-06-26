package v1alpha1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzureDiskDriver is a specification for a AzureDiskDriver resource
type AzureDiskDriver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureDiskDriverSpec   `json:"spec"`
	Status AzureDiskDriverStatus `json:"status"`
}

// AzureDiskDriverSpec is the spec for a AzureDiskDriver resource
type AzureDiskDriverSpec struct {
	operatorv1.OperatorSpec `json:",inline"`
}

// AzureDiskDriverStatus is the status for a AzureDiskDriver resource
type AzureDiskDriverStatus struct {
	operatorv1.OperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzureDiskDriverList is a list of AzureDiskDriver resources
type AzureDiskDriverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []AzureDiskDriver `json:"items"`
}
