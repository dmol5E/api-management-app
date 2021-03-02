package discovery

import (
	"context"
	"flag"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	clientset "github.com/dmol5e/api-management-app/api-publisher/pkg/client/clientset/versioned"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	Debug      bool
	Kubeconfig *string
)

func init() {
	if home := homedir.HomeDir(); home != "" {
		Kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		Kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
}

func CreateRouteConfigClientSet() (*clientset.Clientset, error) {
	// creates the in-cluster config
	var config *rest.Config
	var err error
	if !Debug {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.WithError(err).Panic("Can't set up cluster config")
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *Kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	kubeClient, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

func CreateApiExtensionClientSet() (apiextension.Interface, error) {
	var config *rest.Config
	var err error
	if !Debug {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.WithError(err).Panic("Can't set up cluster config")
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *Kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	kubeClient, err := apiextension.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

func StartWatching(ctx context.Context, client *clientset.Clientset, namespace string, eventHandler func(event watch.Event) error) (chan int, error) {
	stop := make(chan int)
	watcher, err := client.ApimanagementV1alpha1().RouteConfigs(namespace).Watch(ctx, v1.ListOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("K8s watcher client created. Start watching...")
	go func() {

		ch := watcher.ResultChan()
		defer func() {
			watcher.Stop()
			log.Info("Wathcing stopped")
		}()
		for {
			select {
			case event := <-ch:
				log.Infof("Event type: %s, object: %v", event.Type, event.Object)
				if err := eventHandler(event); err != nil {
					log.Errorf("Failed to handle event %v: %v", event, err)
				}
			case <-stop:
				log.Info("Watching CR stopped")
				return
			}
		}
	}()
	return stop, nil
}
