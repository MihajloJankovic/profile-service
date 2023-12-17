package handlers

import (
	"context"
	"errors"
	protosAuth "github.com/MihajloJankovic/Auth-Service/protos/main"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"log"
	"strings"
)

type myProfileServer struct {
	protos.UnimplementedProfileServer
	logger *log.Logger
	// NoSQL: injecting product repository
	repo *ProfileRepo
	cc   protosAuth.AuthClient
}

func NewServer(l *log.Logger, r *ProfileRepo, cc protosAuth.AuthClient) *myProfileServer {
	return &myProfileServer{*new(protos.UnimplementedProfileServer), l, r, cc}
}

// add edit,create user ,delete user
func (s myProfileServer) GetProfile(ctx context.Context, in *protos.ProfileRequest) (*protos.ProfileResponse, error) {
	// Validate required field
	if in.GetEmail() == "" {
		return nil, errors.New("Invalid input. Email is required.")
	}

	out, err := s.repo.GetById(strings.TrimSpace(in.GetEmail())) // Added trim here
	s.logger.Println(out)
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}
	return out, nil
}

func (s myProfileServer) SetProfile(ctx context.Context, in *protos.ProfileSet) (*protos.Empty, error) {
	// Validate required fields
	if strings.TrimSpace(in.GetEmail()) == "" || strings.TrimSpace(in.GetFirstname()) == "" || strings.TrimSpace(in.GetLastname()) == "" || strings.TrimSpace(in.GetUsername()) == "" {
		return nil, errors.New("Invalid input. Email, firstname, username, and lastname are required.")
	}

	// Validate email format
	if !isValidEmailFormat(strings.TrimSpace(in.GetEmail())) { // Added trim here
		return nil, errors.New("Invalid email format.")
	}

	// Additional validation for other fields if needed

	out := new(protos.ProfileResponse)
	out.Email = strings.TrimSpace(in.GetEmail())         // Added trim here
	out.Firstname = strings.TrimSpace(in.GetFirstname()) // Added trim here
	out.Lastname = strings.TrimSpace(in.GetLastname())   // Added trim here
	out.Birthday = strings.TrimSpace(in.GetBirthday())   // Added trim here
	out.Gender = strings.TrimSpace(in.GetGender())       // Added trim here
	out.Role = strings.TrimSpace(in.GetRole())           // Added trim here
	out.Username = strings.TrimSpace(in.GetUsername())   // Added trim here

	err := s.repo.Create(out)
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}
	out2 := new(protosAuth.AuthRequest)
	out2.Email = strings.TrimSpace(in.GetEmail())
	out2.Password = strings.TrimSpace(in.GetPassword())
	_, err = s.cc.Register(context.Background(), out2)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		return nil, err
	}
	return new(protos.Empty), nil
}

func (s myProfileServer) UpdateProfile(ctx context.Context, in *protos.ProfileResponse) (*protos.Empty, error) {
	// Validate required fields
	if strings.TrimSpace(in.GetEmail()) == "" || strings.TrimSpace(in.GetFirstname()) == "" || strings.TrimSpace(in.GetLastname()) == "" {
		return nil, errors.New("Invalid input. Email, firstname, and lastname are required.")
	}

	// Validate email format
	if !isValidEmailFormat(strings.TrimSpace(in.GetEmail())) { // Added trim here
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
	email = strings.TrimSpace(email) // Trim leading and trailing whitespaces
	// Perform a simple check for '@' and '.com'
	return strings.Contains(email, "@") && strings.HasSuffix(email, ".com")
}

func (s myProfileServer) DeleteProfile(ctx context.Context, in *protos.ProfileRequest) (*protos.Empty, error) {
	err := s.repo.DeleteByEmail(strings.TrimSpace(in.GetEmail())) // Added trim here
	if err != nil {
		return nil, err
	}
	return new(protos.Empty), nil
}
