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
		s.logger.Println(err)
		return nil, err
	}
	return out, nil
}
func (s myProfileServer) SetProfile(kon context.Context, in *protos.ProfileResponse) (*protos.Empty, error) {

	out := new(protos.ProfileResponse)
	out.Email = in.GetEmail()
	out.Firstname = in.GetFirstname()
	out.Lastname = in.GetLastname()
	out.Birthday = in.GetBirthday()
	out.Gender = in.GetGender()

	err := s.repo.Create(out)
	if(err != nil) {
		s.logger.Println(err)
        return nil, err
	}
	return new(protos.Empty), nil
}
