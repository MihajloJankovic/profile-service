package main

import (
	"context"
	"github.com/MihajloJankovic/profile-service/handlers"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func main() {

	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatal(err)
	}
	serverRegistar := grpc.NewServer()

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[profile-main] ", log.LstdFlags)
	profilelog := log.New(os.Stdout, "[profile-repo-log] ", log.LstdFlags)

	profileRepo, err := handlers.New(timeoutContext, profilelog)
	if err != nil {
		logger.Fatal(err)
	}
	defer profileRepo.Disconnect(timeoutContext)

	// NoSQL: Checking if the connection was established
	profileRepo.Ping()

	//Initialize the handler and inject said logger
	service := handlers.NewServer(logger, profileRepo)

	protos.RegisterProfileServer(serverRegistar, service)
	err = serverRegistar.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
