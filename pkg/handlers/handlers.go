package handlers

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

// Handlers returns the handlers for the informer
func Handlers() cache.ResourceEventHandlerFuncs {
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)

			logrus.WithFields(logrus.Fields{
				"name":      u.GetName(),
				"namespace": u.GetNamespace(),
				"kind":      u.GroupVersionKind().GroupKind().String(),
				"labels":    u.GetLabels(),
			}).Info("received add event!")
		},
		UpdateFunc: func(oldObj, obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			logrus.WithFields(logrus.Fields{
				"kind": u.GroupVersionKind().GroupKind().String(),
			}).Info("received update event!")
		},
		DeleteFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			logrus.WithFields(logrus.Fields{
				"kind": u.GroupVersionKind().GroupKind().String(),
			}).Info("received delete event!")
		},
	}

	return handlers
}