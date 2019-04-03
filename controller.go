package main

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/kacejot/ownership-controller/pkg/apis/owner/v1alpha1"
	clientset "github.com/kacejot/ownership-controller/pkg/client/clientset/versioned"
	informers "github.com/kacejot/ownership-controller/pkg/client/informers/externalversions"
	ownerinformer "github.com/kacejot/ownership-controller/pkg/client/informers/externalversions/owner/v1alpha1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

// OwnershipController check that all owned resoruces are created
// Otherwise it deletes all resources owned by owner resource
type OwnershipController struct {
	client              *clientset.Clientset
	informerFactory     informers.SharedInformerFactory
	kubeInformerFactory kubeinformers.SharedInformerFactory
	ownerInformer       ownerinformer.OwnerInformer
}

// NewOwnershipController creates controller for Owner resource
func NewOwnershipController() *OwnershipController {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed creating in-cluster config: %v\n", err)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed creating ownership client: %v\n", err)
	}

	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed creating kubernetes client: %v\n", err)
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeclient, time.Second*30)
	informerFactory := informers.NewSharedInformerFactory(client, time.Second*30)
	informer := informerFactory.Myproject().V1alpha1().Owners()

	controller := &OwnershipController{
		client:              client,
		informerFactory:     informerFactory,
		kubeInformerFactory: kubeInformerFactory,
		ownerInformer:       informer,
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
	rc.informerFactory.Start(stopCh)
	rc.kubeInformerFactory.Start(stopCh)
}

func (rc *OwnershipController) onCreate(resource interface{}) {
	owner, ok := resource.(v1alpha1.Owner)
	if !ok {
		log.Fatalf("Failed to cast to owner type")
	}

	for _, owned := range owner.Spec.OwnedResources {
		gvk := schema.FromAPIVersionAndKind(owned.APIVersion, owned.Kind)
		gvr, _ := meta.UnsafeGuessKindToResource(gvk)

		genericInformer, err := rc.kubeInformerFactory.ForResource(gvr)
		if err != nil {
			log.Fatalf("Failed to create generic informer: %v\n", err)
		}

		_, err = genericInformer.Lister().ByNamespace(owned.Namespace).Get(owned.Name)
		if err != nil {
			log.Printf("Failed to retrieve resource: %v\n", err)
			// Delete other resources
		}
	}
}

func (rc *OwnershipController) onUpdate(oldResource, newResource interface{}) {
	oldKey := rc.getResourceKey(oldResource)
	newKey := rc.getResourceKey(newResource)

	log.Printf("Owner %s updated to %s", oldKey, newKey)
}

func (rc *OwnershipController) onDelete(resource interface{}) {
	key := rc.getResourceKey(resource)
	log.Printf("Owner deleted: %s", key)
}

func (rc *OwnershipController) getResourceKey(resource interface{}) string {
	if key, err := cache.MetaNamespaceKeyFunc(resource); err != nil {
		log.Fatalf("Error retrieving owner key: %v", err)
	} else {
		return key
	}

	return ""
}
