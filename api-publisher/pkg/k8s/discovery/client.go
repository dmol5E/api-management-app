package discovery

import (
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
)

func InitClient() error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	kubeClient, err := apiextension.NewForConfig(config)
	if err != nil {
		return err
	}
}
