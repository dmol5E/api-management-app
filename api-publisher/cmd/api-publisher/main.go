package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	applog "github.com/dmol5e/api-management-app/api-publisher/pkg/log"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog/v2"

	"github.com/dmol5e/api-management-app/api-publisher/pkg/apis/apimanagement/v1alpha1"
	"github.com/dmol5e/api-management-app/api-publisher/pkg/client/clientset/versioned"
	"github.com/dmol5e/api-management-app/api-publisher/pkg/k8s/discovery"
	"github.com/dmol5e/api-management-app/api-publisher/pkg/transport/xds"
	"google.golang.org/grpc"
)

const (
	grpcMaxConcurrentStreams = 1000000
	xdsPort                  = 15010
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{DisableColors: false})
	klog.SetLogger(&applog.StaticLogger{})
}

func main() {
	ctx := context.Background()

	var namespace string
	namespace, found := os.LookupEnv("CLOUD_NAMESPACE")
	if !found {
		namespace = "default"
	}

	log.Infof("Starting application. Namespace: %s", namespace)

	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	grpcServer := grpc.NewServer(grpcOptions...)
	xdsServer := xds.PrepareServer(ctx, grpcServer)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", xdsPort))
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("management server listening on %d", xdsPort)
		if err = grpcServer.Serve(lis); err != nil {
			log.Println(err)
		}
	}()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("/health")
		fmt.Fprint(w, "{\"status\":\"UP\"}")
	})

	apiExtClient, err := discovery.CreateApiExtensionClientSet()
	if err != nil {
		log.Panicf("Failed to create ClientSet for k8s: %v", err)
	}
	_, err = v1alpha1.CreateCRD(ctx, apiExtClient)
	if err != nil {
		log.Panicf("Failed to create CRD RouteConfig: %v", err)
	}

	routeConfigClient, err := discovery.CreateRouteConfigClientSet()
	if err != nil {
		log.Panicf("Failed to create ClientSet for handling RouteConfig CR: %v", err)
	}
	handler := MakeEventHandler(xdsServer, routeConfigClient, namespace)
	stopCh, err := discovery.StartWatching(ctx, routeConfigClient, namespace, handler)
	if err != nil {
		log.Panicf("Failed to start Watching RouteConfig: %v", err)
	}
	defer func() {
		stopCh <- 0
	}()

	log.Info("App has started on port 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func MakeEventHandler(ctx context.Context, xdsServer *xds.Server, routeConfigClient *versioned.Clientset, namespace string) func(watch.Event) error {
	return func(e watch.Event) error {
		routeConfigList, err := routeConfigClient.ApimanagementV1alpha1().RouteConfigs(namespace).List(ctx, v1.ListOptions{})
		if err != nil {
			log.WithError(err).Errorf("Failed to get list of RouteConfig")
		}

		gatewayConfigs := MakeGatewayConfiguration(routeConfigList.Items)
		for k, v := range gatewayConfigs {
			snapshot := cache.NewSnapshot(
				uuid.New().String(),
				v.GetEndpoints(),
				v.GetClusters(),
				v.GetRouteConfigs(),
				v.GetListeners(),
				v.GetRuntimes(),
				v.GetSecrets(),
			)
			xdsServer.Update(k, snapshot)
		}
		return nil
	}
}

type SnapshotResources struct {
	Endpoints    map[string]types.Resource
	Clusters     map[string]types.Resource
	RouteConfigs map[string]types.Resource
	Listeners    map[string]types.Resource
	Runtimes     map[string]types.Resource
	Secrets      map[string]types.Resource
}

func (r SnapshotResources) GetEndpoints() []types.Resource {
	endpoints := make([]types.Resource, len(r.Endpoints))
	for _, v := range r.Endpoints {
		endpoints = append(endpoints, v)
	}
	return endpoints
}

func (r SnapshotResources) GetClusters() []types.Resource {
	clusters := make([]types.Resource, len(r.Clusters))
	for _, v := range r.Clusters {
		clusters = append(clusters, v)
	}
	return clusters
}

func (r SnapshotResources) GetRouteConfigs() []types.Resource {
	routeConfigs := make([]types.Resource, len(r.RouteConfigs))
	for _, v := range r.RouteConfigs {
		routeConfigs = append(routeConfigs, v)
	}
	return routeConfigs
}

func (r SnapshotResources) GetListeners() []types.Resource {
	listeners := make([]types.Resource, len(r.Listeners))
	for _, v := range r.Listeners {
		listeners = append(listeners, v)
	}
	return listeners
}

func (r SnapshotResources) GetRuntimes() []types.Resource {
	runtimes := make([]types.Resource, len(r.Runtimes))
	for _, v := range r.Runtimes {
		runtimes = append(runtimes, v)
	}
	return runtimes
}

func (r SnapshotResources) GetSecrets() []types.Resource {
	secrets := make([]types.Resource, len(r.Secrets))
	for _, v := range r.Secrets {
		secrets = append(secrets, v)
	}
	return secrets
}

func (r SnapshotResources) Merge(newRes SnapshotResources) SnapshotResources {
	return SnapshotResources{}
}

func MakeGatewayConfiguration(routeConfigs []v1alpha1.RouteConfig) map[string]SnapshotResources {
	result := make(map[string]SnapshotResources)
	for _, routeConfig := range routeConfigs {
		gateway := routeConfig.Spec.Gateway
		if snapshotRes, found := result[gateway]; !found {
			result[gateway] = MakeSnapshotResources(routeConfig.Spec)
		} else {
			result[gateway] = snapshotRes.Merge(MakeSnapshotResources(routeConfig.Spec))
		}
	}
	return result
}

func MakeSnapshotResources(routeConfig v1alpha1.RouteConfigSpec) SnapshotResources {
	return SnapshotResources{}
}
