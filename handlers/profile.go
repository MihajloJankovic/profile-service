package handlers

import (
	"context"
	"errors"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"log"
	"strings"
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
	// Validate required field
	if in.GetEmail() == "" {
		return nil, errors.New("Invalid input. Email is required.")
	}

	out, err := s.repo.GetById(in.GetEmail())
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}
	return out, nil
}

func (s myProfileServer) SetProfile(ctx context.Context, in *protos.ProfileResponse) (*protos.Empty, error) {
	// Validate required fields
	if in.GetEmail() == "" || in.GetFirstname() == "" || in.GetLastname() == "" || in.GetUsername() == ""{
		return nil, errors.New("Invalid input. Email, firstname, username, and lastname are required.")
	}

	// Validate email format
	if !isValidEmailFormat(in.GetEmail()) {
		return nil, errors.New("Invalid email format.")
	}

	// Additional validation for other fields if needed

	out := new(protos.ProfileResponse)
	out.Email = in.GetEmail()
	out.Firstname = in.GetFirstname()
	out.Lastname = in.GetLastname()
	out.Birthday = in.GetBirthday()
	out.Gender = in.GetGender()
	out.Role = in.GetRole()
	out.Username= in.Username

	err := s.repo.Create(out)
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}
	return new(protos.Empty), nil
}

func (s myProfileServer) UpdateProfile(ctx context.Context, in *protos.ProfileResponse) (*protos.Empty, error) {
	// Validate required fields
	if in.GetEmail() == "" || in.GetFirstname() == "" || in.GetLastname() == "" {
		return nil, errors.New("Invalid input. Email, firstname, and lastname are required.")
	}

	// Validate email format
	if !isValidEmailFormat(in.GetEmail()) {
		return nil, errors.New("Invalid email format.")
	}

	// Additional validation for other fields if needed

	err := s.repo.Update(in)
	if err != nil {
		return nil, err
	}
	return new(protos.Empty), nil
}

// isValidEmailFormat checks if the given email is in a valid format.
func isValidEmailFormat(email string) bool {
	// Perform a simple check for '@' and '.com'
	return strings.Contains(email, "@") && strings.HasSuffix(email, ".com")
}
