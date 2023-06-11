package controllers

import (
	"os"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"k8s-global-view/pkg/handlers"
)

// Controller is the controller for the operator
type Controller struct {
	Clientset kubernetes.Interface

	DiscoveryClient discovery.DiscoveryInterface

	DynamicClient dynamic.Interface

	Log *logrus.Entry

	Group string
}

// NewController creates a new controller
func NewController(log *logrus.Logger, cfg *rest.Config) (*Controller, error) {
	l := log.WithFields(logrus.Fields{
		"region": cfg.Host,
	})
	l.Info("Creating controller")
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		l.WithError(err).Error("Failed to create clientset")
		return nil, err
	}

	discoveryClient := clientset.Discovery()

	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		l.WithError(err).Error("Failed to create dynamic client")
		return nil, err
	}

	group := os.Getenv("GROUP")

	return &Controller{
		Log:             l,
		Clientset:       clientset,
		DiscoveryClient: discoveryClient,
		DynamicClient:   dc,
		Group:           group,
	}, nil
}

// Start starts the controller
func (c *Controller) Start(stopCh <-chan struct{}) error {
	l := c.Log
	apiResources, err := c.DiscoveryClient.ServerPreferredResources()
	if err != nil {
		l.WithError(err).Error("Failed to get server preferred resources")
		return err
	}

	groupVersionMap := mapAPIResourcesByGroup(apiResources, c.Group)

	// TODO: probably allow for filtering via labels
	f := dynamicinformer.NewFilteredDynamicSharedInformerFactory(c.DynamicClient, 0, v1.NamespaceAll, nil)

	resourceArgs := resourceArgList(groupVersionMap)
	informers := setupInformers(f, resourceArgs, handlers.Handlers(l))

	startInformers(l, informers, stopCh)

	return nil
}
