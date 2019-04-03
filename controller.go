package main

import (
	"log"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	clientset "github.com/kacejot/rep-controller/pkg/client/clientset/versioned"
	informers "github.com/kacejot/rep-controller/pkg/client/informers/externalversions"
)

// OwnershipController check that all owned resoruces are created
// Otherwise it deletes all resources owned by owner resource
type OwnershipController struct {
	Client          *clientset.Clientset
	InformerFactory informers.SharedInformerFactory
}

// NewCreationController creates controller for Owner resource
func NewOwnershipController() *OwnershipController {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed creating in-cluster config: %v\n", err)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed creating kubernetes client: %v\n", err)
	}

	informerFactory := informers.NewSharedInformerFactory(client, time.Second*30)
	informer := informerFactory.Myproject().V1alpha1().Owners()

	controller := &OwnershipController{
		Client:          client,
		InformerFactory: informerFactory,
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.onCreate,
		UpdateFunc: controller.onUpdate,
		DeleteFunc: controller.onDelete,
	})

	return controller
}

// Run starts informer that monitors cluster for resource events
func (rc *OwnershipController) Run(stopCh <-chan struct{}) {
	rc.InformerFactory.Start(stopCh)
}

func (rc *OwnershipController) onCreate(resource interface{}) {
	key := rc.getResourceKey(resource)
	log.Printf("Policy created: %s", key)
}

func (rc *OwnershipController) onUpdate(oldResource, newResource interface{}) {
	oldKey := rc.getResourceKey(oldResource)
	newKey := rc.getResourceKey(newResource)

	log.Printf("Policy %s updated to %s", oldKey, newKey)
}

func (rc *OwnershipController) onDelete(resource interface{}) {
	key := rc.getResourceKey(resource)
	log.Printf("Policy deleted: %s", key)
}

func (rc *OwnershipController) getResourceKey(resource interface{}) string {
	if key, err := cache.MetaNamespaceKeyFunc(resource); err != nil {
		log.Fatalf("Error retrieving policy key: %v", err)
	} else {
		return key
	}

	return ""
}
