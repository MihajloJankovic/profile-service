package handlers

import (
	"context"
	"fmt"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
)

type myProfileServer struct {
	protos.UnimplementedProfileServer
}

func NewServer() *myProfileServer {
	return &myProfileServer{}
}

// add edit,create user ,delete user
func (s myProfileServer) GetProfile(ctx context.Context, in *protos.ProfileRequest) (*protos.ProfileResponse, error) {

	out := new(protos.ProfileResponse)

	//add db
	if in.Email == "pera@gmail.com" {
		out.Email = "pera@gmail.com"
		out.Firstname = "MIhajlo"
		out.Lastname = "Jankovic"
		out.Birthday = "23.04.2002"
		out.Gender = false
		return out, nil
	} else {
		return nil, fmt.Errorf("Greska")
	}
	return out, nil
}
