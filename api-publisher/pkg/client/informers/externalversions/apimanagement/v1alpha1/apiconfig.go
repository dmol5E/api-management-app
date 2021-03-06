/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	apimanagementv1alpha1 "github.com/dmol5e/api-management-app/api-publisher/pkg/apis/apimanagement/v1alpha1"
	versioned "github.com/dmol5e/api-management-app/api-publisher/pkg/client/clientset/versioned"
	internalinterfaces "github.com/dmol5e/api-management-app/api-publisher/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/dmol5e/api-management-app/api-publisher/pkg/client/listers/apimanagement/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// APIConfigInformer provides access to a shared informer and lister for
// APIConfigs.
type APIConfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.APIConfigLister
}

type aPIConfigInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewAPIConfigInformer constructs a new informer for APIConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewAPIConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredAPIConfigInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredAPIConfigInformer constructs a new informer for APIConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredAPIConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApimanagementV1alpha1().APIConfigs(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApimanagementV1alpha1().APIConfigs(namespace).Watch(context.TODO(), options)
			},
		},
		&apimanagementv1alpha1.APIConfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *aPIConfigInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredAPIConfigInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *aPIConfigInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apimanagementv1alpha1.APIConfig{}, f.defaultInformer)
}

func (f *aPIConfigInformer) Lister() v1alpha1.APIConfigLister {
	return v1alpha1.NewAPIConfigLister(f.Informer().GetIndexer())
}