package handlers

import (
	"context"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"log"
)

type myProfileServer struct {
	protos.UnimplementedProfileServer
	logger *log.Logger
	// NoSQL: injecting product repository
	repo *ProfileRepo
}

func NewServer(l *log.Logger, r *ProfileRepo) *myProfileServer {
	return &myProfileServer{*new(protos.UnimplementedProfileServer), l, r}
}

// add edit,create user ,delete user
func (s myProfileServer) GetProfile(ctx context.Context, in *protos.ProfileRequest) (*protos.ProfileResponse, error) {

	out, err := s.repo.GetById(in.GetEmail())
	if err != nil {
		s.logger.Fatal(err)
		return nil, err
	}
	//
	////add db
	//if in.Email == "pera@gmail.com" {
	//	out.Email = "pera@gmail.com"
	//	out.Firstname = "MIhajlo"
	//	out.Lastname = "Jankovic"
	//	out.Birthday = "23.04.2002"
	//	out.Gender = false
	//	return out, nil
	//} else {
	//	return nil, fmt.Errorf("Greska")
	//}
	return out, nil
}
func (s myProfileServer) SetProfile(kon context.Context, in *protos.ProfileResponse) (*protos.Empty, error) {

	out := new(protos.ProfileResponse)
	out.Email = in.GetEmail()
	out.Firstname = in.GetFirstname()
	out.Lastname = in.GetLastname()
	out.Birthday = in.GetBirthday()
	out.Gender = in.GetGender()

	s.logger.Println(s.repo.Create(out))
	return new(protos.Empty), nil
}
