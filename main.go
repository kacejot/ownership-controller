package main

import (
	"log"

	"github.com/kacejot/ownership-controller/pkg/signals"
)

func main() {
	stopCh := signals.SetupSignalHandler()
	controller := NewOwnershipController()
	controller.Run(stopCh)

	log.Println("Controller has started")
	<-stopCh
	log.Println("Controller has stopped")
}
