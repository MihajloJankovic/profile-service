package main

import (
	"context"
	protosAuth "github.com/MihajloJankovic/Auth-Service/protos/main"
	"github.com/MihajloJankovic/profile-service/handlers"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	connAuth, err := grpc.Dial("auth-service:9094", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(connAuth)
	ccAuth := protosAuth.NewAuthClient(connAuth)
	//Initialize the handler and inject said logger
	service := handlers.NewServer(logger, profileRepo, ccAuth)

	protos.RegisterProfileServer(serverRegistar, service)
	err = serverRegistar.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
