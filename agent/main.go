package main

import (
	"agent/server"
	pb "clarity/generated/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAgentServiceServer(grpcServer, server.NewAgentServer())
	reflection.Register(grpcServer)

	log.Printf("🚀 Clarity Agent is running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
