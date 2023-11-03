package main

import (
	"github.com/MihajloJankovic/profile-service/handlers"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatal(err)
	}
	serverRegistar := grpc.NewServer()
	service := handlers.NewServer()

	protos.RegisterProfileServer(serverRegistar, service)
	err = serverRegistar.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
