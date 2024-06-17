package main

import (
	"flag"
	"fmt"
	pb "github.com/kyma-project/networking-dev-tools/grpcbin/pkg/hello"
	"github.com/kyma-project/networking-dev-tools/grpcbin/pkg/hello_server"
	"google.golang.org/grpc"
	grpcreflection "google.golang.org/grpc/reflection"
	"log"
	"net"
)

var (
	port       string
	reflection bool
)

func init() {
	flag.StringVar(&port, "port", "50051", "The port to listen on")
	flag.BoolVar(&reflection, "reflection", false, "Enable gRPC reflection")

	flag.Parse()
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	if reflection {
		grpcreflection.Register(server)
	}

	pb.RegisterHelloServiceServer(server, &hello_server.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
