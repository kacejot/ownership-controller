package main

import (
	"log"

	signals "github.com/kacejot/rep-controller/pkg/signals"
)

func main() {

	stopCh := signals.SetupSignalHandler()
	controller := NewRepoController()
	controller.Run(stopCh)

	<-stopCh

	log.Println("Controller has stopped")
}
