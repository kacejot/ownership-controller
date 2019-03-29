package main

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type RepoController struct {
}

func main() {
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed creating in-cluster config: %v\n", err)
	}

	client * kubernetes.Clientset
}
