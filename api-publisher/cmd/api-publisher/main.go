package main

import (
	"context"
	"fmt"
	logger "log"
	"net"
	"net/http"
	"os"

	"github.com/dmol5e/api-management-app/api-publisher/pkg/transport/xds"

	"google.golang.org/grpc"
)

const (
	grpcMaxConcurrentStreams = 1000000
	xdsPort                  = 15010
)

var (
	log *logger.Logger
)

func init() {
	log = &logger.Logger{}
	log.SetOutput(os.Stdout)
}

func main() {
	ctx := context.Background()

	log.Println("Starting application")

	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	grpcServer := grpc.NewServer(grpcOptions...)
	xds.RunServer(ctx, grpcServer, log)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", xdsPort))
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("management server listening on %d\n", xdsPort)
		if err = grpcServer.Serve(lis); err != nil {
			log.Println(err)
		}
	}()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("/health")
		fmt.Fprint(w, "{\"status\":\"UP\"}")
	})

	log.Println("App has started on port 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
