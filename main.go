package main

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"k8s-global-view/pkg/controllers"
)

func main() {
	controller, err := controllers.NewController()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create controller")
	}

	stopCh := make(chan struct{})
	controller.Start(stopCh)

	sigCh := make(chan os.Signal, 0)
	signal.Notify(sigCh, os.Kill, os.Interrupt)

	<-sigCh
	close(stopCh)
}
