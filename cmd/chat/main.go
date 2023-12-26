package main

import (
	"flag"
	"log"
	"net"

	"github.com/manujelko/grpc-example/internal/chat"
	pb "github.com/manujelko/grpc-example/pkg/api/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var address string
	flag.StringVar(&address, "address", "", "Address for the service to listen on")
	flag.Parse()

	if address == "" {
		log.Panic("Address flag is required")
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	server := chat.NewServer()

	pb.RegisterSimpleChatServer(grpcServer, server)

	log.Printf("Starting service on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
