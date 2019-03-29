package v1alpha1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Repo struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            RepoSpec   `json:"spec"`
	Status          RepoStatus `json:"status"`
}

type RepoSpec struct {
	RepoAddress   string   `json:"repoAddress"`
	DockerImage   string   `json:"dockerImage"`
	BuildCommands []string `json:"buildCommands"`
	Artifacts     []string `json:"artifacts"`
}

type RepoStatus struct {
	Log []string `json:"log"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RepoList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata"`
	Items         []Repo `json:"items"`
}
