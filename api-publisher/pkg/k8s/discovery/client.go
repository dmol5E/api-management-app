package discovery

import (
	"context"

	log "github.com/sirupsen/logrus"

	clientset "github.com/dmol5e/api-management-app/api-publisher/pkg/client/clientset/versioned"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func CreateRouteConfigClientSet() (*clientset.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.WithError(err).Panic("Can't set up cluster config")
	}

	kubeClient, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

func CreateApiExtensionClientSet() (apiextension.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.WithError(err).Panic("Can't set up cluster config")
	}

	kubeClient, err := apiextension.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

func StartWatching(ctx context.Context, client *clientset.Clientset, namespace string) (chan int, error) {
	stop := make(chan int)
	watcher, err := client.ApimanagementV1alpha1().RouteConfigs(namespace).Watch(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	go func() {
		ch := watcher.ResultChan()
		defer watcher.Stop()
		for {
			select {
			case event := <-ch:
				log.Infof("Event type: %s, object: %v", event.Type, event.Object)
			case <-stop:
				log.Info("Watching CR stopped")
				return
			}
		}
	}()
	return stop, nil
}
