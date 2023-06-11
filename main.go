package main

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/RealHarshThakur/supercontroller/pkg/controllers"
	"github.com/RealHarshThakur/supercontroller/pkg/setup"
)

func main() {
	l := setup.Logging()

	config, err := setup.LoadConfig()
	if err != nil {
		l.Fatal(err)
	}

	kcs := setup.BuildKubernetesConfigs(l.WithField("setup", "kubeconfigs"), config)
	stopCh := make(chan struct{})
	for _, kc := range kcs.Configs {
		controller, err := controllers.NewController(l, kc)
		if err != nil {
			logrus.WithError(err).Error("Failed to create controller for", kc.Host)
		}

		controller.Start(stopCh)
	}

	sigCh := make(chan os.Signal, 0)
	signal.Notify(sigCh, os.Kill, os.Interrupt)

	<-sigCh
	close(stopCh)
}
