package main

import (
	"context"
	"fmt"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"google.golang.org/grpc"
	"log"
	"net"
)

type myInvoiceServer struct {
	protos.UnimplementedProfileServer
}

func GetProfile(ctx context.Context, in *protos.ProfileRequest) (*protos.ProfileResponse, error) {

	out := new(protos.ProfileResponse)

	fmt.Println(in.GetEmail())
	fmt.Println(in.Email)
	if in.Email == "pera@gmail.com" {
		out.Email = "pera@gmail.com"
		out.Firstname = "MIhajlo"
		out.Lastname = "Jankovic"
		out.Birthday = "23.04.2002"
		out.Gender = false
		return out, nil
	} else {
		return nil, fmt.Errorf("Greskaaa")
	}
	return out, nil
}
func mustEmbedUnimplementedProfileServer() {

}
func main() {
	
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatal(err)
	}
	serverRegistar := grpc.NewServer()
	service := &myInvoiceServer{}

	protos.RegisterProfileServer(serverRegistar, service)
	err = serverRegistar.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
