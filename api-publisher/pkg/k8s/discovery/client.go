package discovery

import (
	"context"
	"fmt"

	clientset "github.com/dmol5e/api-management-app/api-publisher/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func InitClient(ctx context.Context, namespace string) (chan int, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	stop := make(chan int)

	kubeClient, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	watcher, err := kubeClient.ApimanagementV1alpha1().RouteConfigs("").Watch(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	go func() {
		ch := watcher.ResultChan()
		defer watcher.Stop()
		for {
			select {
			case event := <-ch:
				fmt.Printf("Event type: %s, object: %v", event.Type, event.Object)
			case <-stop:
				return
			}
		}
	}()
	return stop, nil
}
