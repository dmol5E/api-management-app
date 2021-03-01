package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	applog "github.com/dmol5e/api-management-app/api-publisher/pkg/log"
	log "github.com/sirupsen/logrus"
	"k8s.io/klog/v2"

	"github.com/dmol5e/api-management-app/api-publisher/pkg/apis/apimanagement/v1alpha1"
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
	xds.RunServer(ctx, grpcServer)

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
	stopCh, err := discovery.StartWatching(ctx, routeConfigClient, namespace)
	if err != nil {
		log.Panicf("Failed to start Watching RouteConfig: %v", err)
	}
	defer func() {
		stopCh <- 0
	}()

	log.Info("App has started on port 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
