package users

import (
	"github.com/alexedwards/argon2id"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
)

type service struct {
	store UsersStore
}

func NewService(store UsersStore) *service {
	return &service{
		store: store,
	}
}

func (s *service) SignUpUser(req *SignUpUserRequest) (*SignUpUserResponse, error) {
	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	userToInsert := &User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	if err := s.store.InsertUser(userToInsert); err != nil {
		return nil, err
	}

	return &SignUpUserResponse{}, nil
}

func (s *service) SignInUser(req *SignInUserRequest) (*SignInUserResponse, error) {
	userFound, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if userFound == nil {
		return nil, common.InvalidCredentials()
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, userFound.Password)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, common.InvalidCredentials()
	}

	return &SignInUserResponse{}, nil
}
