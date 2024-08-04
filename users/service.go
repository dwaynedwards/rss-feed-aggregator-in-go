package users

import (
	"github.com/alexedwards/argon2id"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
	"github.com/google/uuid"
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
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	userToInsert := &User{
		ID:       id,
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	if !s.store.InsertUser(userToInsert) {
		// InvalidUserExists mornally this workflow would be handled with a status 201 and a message saying an email was sent to
		// verify the user. When this error is hit, an email would be sent saying if you're trying to create
		// an you can trying executing the forgot password workflow instead of leaking internal info to the user
		// that an user already exists with the email provided, but this is outside of the scope of this project
		return nil, common.InvalidUserExists()
	}

	return &SignUpUserResponse{}, nil
}

func (s *service) SignInUser(req *SignInUserRequest) (*SignInUserResponse, error) {
	userFound := s.store.GetUserByEmail(req.Email)
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
