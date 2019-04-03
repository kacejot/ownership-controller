package v1alpha1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Owner struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            OwnerSpec   `json:"spec"`
	Status          OwnerStatus `json:"status"`
}

type OwnerSpec struct {
	OwnedResources   []OwnedResource   `json:"ownedResources"`
}

type OwnedResource struct {
    Name string `json:"name"`
    Namespace string `json:"namespace"`
}

type OwnerStatus struct {
	Log []string `json:"log"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type OwnerList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata"`
	Items         []Owner `json:"items"`
}
