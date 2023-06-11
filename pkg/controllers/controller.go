package controllers

import (
	"fmt"
	"k8s-global-view/pkg/handlers"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// Controller is the controller for the operator
type Controller struct {
	Clientset kubernetes.Interface

	DiscoveryClient discovery.DiscoveryInterface

	DynamicClient dynamic.Interface

	Log *logrus.Logger

	Group string
}

// NewController creates a new controller
func NewController() (*Controller, error) {
	l := SetupLogging()

	cfg, err := restConfig()
	if err != nil {
		l.WithError(err).Error("Failed to create rest config")
		return nil, err
	}

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
	informers := setupInformers(f, resourceArgs, handlers.Handlers())

	go startInformers(c.Log, informers, stopCh)

	return nil
}

func startInformers(l *logrus.Logger, informers []informers.GenericInformer, stopCh <-chan struct{}) {
	for _, informer := range informers {
		go informer.Informer().Run(stopCh)
	}
}

func setupInformers(f dynamicinformer.DynamicSharedInformerFactory, resourceArgs []string, handlers cache.ResourceEventHandlerFuncs) []informers.GenericInformer {
	informers := make([]informers.GenericInformer, 0, len(resourceArgs))
	for _, resourceArg := range resourceArgs {
		gvr, _ := schema.ParseResourceArg(resourceArg)
		i := f.ForResource(*gvr)
		i.Informer().AddEventHandler(handlers)
		informers = append(informers, i)
	}

	return informers
}

func resourceArgList(groupVersionMap map[string][]*v1.APIResourceList) []string {
	argList := make([]string, 0)
	for groupVersion, apiResourceLists := range groupVersionMap {
		for _, apiResourceList := range apiResourceLists {
			for _, apiResource := range apiResourceList.APIResources {
				groupVersionSplit := strings.Split(groupVersion, "/")
				group := groupVersionSplit[0]
				version := groupVersionSplit[1]
				argList = append(argList, fmt.Sprintf("%s.%s.%s", apiResource.Name, version, group))
			}
		}
	}
	return argList
}

func mapAPIResourcesByGroup(apiResources []*v1.APIResourceList, group string) map[string][]*v1.APIResourceList {
	groupVersionMap := make(map[string][]*v1.APIResourceList, 0)
	for _, apiResourceList := range apiResources {
		if strings.Contains(apiResourceList.GroupVersion, group) {
			groupVersionMap[apiResourceList.GroupVersion] = append(groupVersionMap[apiResourceList.GroupVersion], apiResourceList)
		}
	}

	return groupVersionMap

}
