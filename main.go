package main

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Repo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RepoSpec   `json:"spec"`
	Status            RepoStatus `json:"status"`
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

type RepoController struct {
}

func main() {
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed creating in-cluster config: %v\n", err)
	}

	client * kubernetes.Clientset
}
