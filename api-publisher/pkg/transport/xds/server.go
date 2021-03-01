package xds

import (
	"context"

	log "github.com/sirupsen/logrus"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
)

type Server struct {
	cache cachev3.SnapshotCache
}

//RunServer Start server to serve Envoy xDS requests
func RunServer(ctx context.Context, grpcServer *grpc.Server) *Server {

	cache := cachev3.NewSnapshotCache(false, cachev3.IDHash{}, logger)
	server := serverv3.NewServer(ctx, cache, nil)

	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, server)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)

	return &Server{
		cache: cache,
	}
}

func (s *Server) Update(nodeID string, snapshot cachev3.Snapshot) error {
	if err := snapshot.Consistent(); err != nil {
		log.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
	}
	log.Printf("will serve snapshot %+v", snapshot)
	if err := s.cache.SetSnapshot(nodeID, snapshot); err != nil {
		log.Fatalf("snapshot error %q for %+v", err, snapshot)
		return err
	}
	return nil
}
