package main

import (
	"log"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	clientset "github.com/kacejot/rep-controller/pkg/client/clientset/versioned"
	informers "github.com/kacejot/rep-controller/pkg/client/informers/externalversions"
)

type RepoController struct {
	Client          *clientset.Clientset
	InformerFactory informers.SharedInformerFactory
}

// NewRepoController creates controller for Repo resource
func NewRepoController() *RepoController {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed creating in-cluster config: %v\n", err)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed creating kubernetes client: %v\n", err)
	}

	informerFactory := informers.NewSharedInformerFactory(client, time.Second*30)
	reposInformer := informerFactory.Myproject().V1alpha1().Repos()

	controller := &RepoController{
		Client:          client,
		InformerFactory: informerFactory,
	}

	reposInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onCreate,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})

	return controller
}

func (rc *RepoController) Run(stopCh <-chan struct{}) {
	rc.InformerFactory.Start(stopCh)
}

func (rc *RepoController) onCreate(resource interface{}) {
	key := rc.getResourceKey(resource)
	log.Printf("Policy created: %s", key)
}

func (rc *RepoController) onUpdate(oldResource, newResource interface{}) {
	oldKey := rc.getResourceKey(oldResource)
	newKey := rc.getResourceKey(newResource)

	log.Printf("Policy %s updated to %s", oldKey, newKey)
}

func (rc *RepoController) onDelete(resource interface{}) {
	key := rc.getResourceKey(resource)
	log.Printf("Policy deleted: %s", key)
}

func (rc *RepoController) getResourceKey(resource interface{}) string {
	if key, err := cache.MetaNamespaceKeyFunc(resource); err != nil {
		log.Fatalf("Error retrieving policy key: %v", err)
	} else {
		return key
	}

	return ""
}
