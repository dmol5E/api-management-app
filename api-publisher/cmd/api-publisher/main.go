package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

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
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"
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
		log.Panicf("Failed to create CRD APIConfig: %v", err)
	}

	apiConfigClient, err := discovery.CreateAPIConfigClientSet()
	if err != nil {
		log.Panicf("Failed to create ClientSet for handling APIConfig CR: %v", err)
	}
	handler := MakeEventHandler(ctx, xdsServer, apiConfigClient, namespace)
	stopCh, err := discovery.StartWatching(ctx, apiConfigClient, namespace, handler)
	if err != nil {
		log.Panicf("Failed to start Watching APIConfig: %v", err)
	}
	defer func() {
		stopCh <- 0
	}()

	log.Info("App has started on port 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func MakeEventHandler(ctx context.Context, xdsServer *xds.Server, apiConfigClient *versioned.Clientset, namespace string) func(watch.Event) error {
	return func(e watch.Event) error {
		apiConfigList, err := apiConfigClient.ApimanagementV1alpha1().APIConfigs(namespace).List(ctx, v1.ListOptions{})
		if err != nil {
			log.WithError(err).Errorf("Failed to get list of APIConfig")
		}

		gatewayConfigs := MakeGatewayConfiguration(apiConfigList.Items)
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
	Endpoints    map[string]*endpoint.Endpoint
	Clusters     map[string]*cluster.Cluster
	RouteConfigs map[string]*route.RouteConfiguration
	Listeners    map[string]*listener.Listener
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
	clusters := make([]types.Resource, 0)
	for _, v := range r.Clusters {
		clusters = append(clusters, v)
	}
	return clusters
}

func (r SnapshotResources) GetRouteConfigs() []types.Resource {
	routeConfigs := make([]types.Resource, 0)
	for _, v := range r.RouteConfigs {
		routeConfigs = append(routeConfigs, v)
	}
	return routeConfigs
}

func (r SnapshotResources) GetListeners() []types.Resource {
	listeners := make([]types.Resource, 0)
	for _, v := range r.Listeners {
		listeners = append(listeners, v)
	}
	return listeners
}

func (r SnapshotResources) GetRuntimes() []types.Resource {
	runtimes := make([]types.Resource, 0)
	for _, v := range r.Runtimes {
		runtimes = append(runtimes, v)
	}
	return runtimes
}

func (r SnapshotResources) GetSecrets() []types.Resource {
	secrets := make([]types.Resource, 0)
	for _, v := range r.Secrets {
		secrets = append(secrets, v)
	}
	return secrets
}

func (r SnapshotResources) Merge(newRes SnapshotResources) SnapshotResources {
	for k, v := range newRes.Endpoints {
		r.Endpoints[k] = v
	}
	for k, v := range newRes.Clusters {
		r.Clusters[k] = v
	}
	for k, v := range newRes.RouteConfigs {
		r.RouteConfigs[k] = mergeRouteConfigs(r.RouteConfigs[k], v)
	}
	for k, v := range newRes.Listeners {
		r.Listeners[k] = v
	}
	for k, v := range newRes.Runtimes {
		r.Runtimes[k] = v
	}
	for k, v := range newRes.Secrets {
		r.Secrets[k] = v
	}
	return r
}

func MakeGatewayConfiguration(routeConfigs []v1alpha1.APIConfig) map[string]SnapshotResources {
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

func MakeSnapshotResources(apiConfig v1alpha1.APIConfigSpec) SnapshotResources {
	result := SnapshotResources{
		Endpoints:    make(map[string]*endpoint.Endpoint),
		Clusters:     make(map[string]*cluster.Cluster),
		RouteConfigs: make(map[string]*route.RouteConfiguration),
		Listeners:    make(map[string]*listener.Listener),
	}

	allRoutes := make([]*route.Route, 0)
	for _, apiConfigRoute := range apiConfig.Routes {
		cluster := makeCluster(apiConfigRoute.Destination)
		routes := make([]*route.Route, len(apiConfigRoute.Rules))
		for i, apiConfigRule := range apiConfigRoute.Rules {
			routes[i] = makeRoute(cluster.Name, apiConfigRule)
		}
		result.Clusters[cluster.Name] = cluster
		allRoutes = append(allRoutes, routes...)
	}
	envoyRouteConfig := makeRouteConfig("main-route-config", allRoutes)
	listener := makeHTTPListener("main-listener", "main-route-config")
	result.RouteConfigs[envoyRouteConfig.Name] = envoyRouteConfig
	result.Listeners[listener.Name] = listener
	return result
}

func makeHTTPListener(listenerName string, routeConfigName string) *listener.Listener {
	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: routeConfigName,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: wellknown.Router,
		}},
	}
	pbst, err := ptypes.MarshalAny(manager)
	if err != nil {
		panic(err)
	}

	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: uint32(8080),
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}
}

func makeRouteConfig(routeName string, routes []*route.Route) *route.RouteConfiguration {
	routeConfig := &route.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
			Routes:  routes,
		}},
	}
	return routeConfig
}

func makeRoute(clusterName string, apiConfigRule v1alpha1.Rule) *route.Route {
	return &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Prefix{
				Prefix: apiConfigRule.Match.Path,
			},
			Headers: makeHeaders(apiConfigRule.Match.Headers),
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				PrefixRewrite: apiConfigRule.PathRewrite,
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterName,
				},
				HostRewriteSpecifier: &route.RouteAction_AutoHostRewrite{},
			},
		},
	}
}

func makeHeaders(apiConfigHeaders []v1alpha1.MatchHeader) []*route.HeaderMatcher {
	headerMatchers := make([]*route.HeaderMatcher, len(apiConfigHeaders))
	for i, apiConfigHeader := range apiConfigHeaders {
		headerMatchers[i] = &route.HeaderMatcher{
			Name: apiConfigHeader.Name,
			HeaderMatchSpecifier: &route.HeaderMatcher_ExactMatch{
				ExactMatch: apiConfigHeader.Value,
			},
		}
	}
	return headerMatchers
}

func makeCluster(destination v1alpha1.Destination) *cluster.Cluster {
	host := destination.Address.Host
	port := destination.Address.Port
	clusterName := host + "||" + strconv.FormatUint(uint64(port), 10)
	return &cluster.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       ptypes.DurationProto(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeEndpoint(clusterName, host, port),
		DnsLookupFamily:      cluster.Cluster_V4_ONLY,
	}
}

func makeEndpoint(clusterName string, host string, port uint32) *endpoint.ClusterLoadAssignment {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  host,
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: port,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}

func mergeRouteConfigs(a *route.RouteConfiguration, b *route.RouteConfiguration) *route.RouteConfiguration {
	b.VirtualHosts[0].Routes = append(b.VirtualHosts[0].Routes, a.VirtualHosts[0].Routes...)
	return b
}
