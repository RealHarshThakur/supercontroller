package controllers

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

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
