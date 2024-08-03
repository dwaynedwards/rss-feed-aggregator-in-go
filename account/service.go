package account

import (
	"github.com/alexedwards/argon2id"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
	"github.com/google/uuid"
)

type service struct {
	store AccountStore
}

func NewService(store AccountStore) *service {
	return &service{
		store: store,
	}
}

func (s *service) CreateAccount(req *CreateAccountRequest) (*CreateAccountResponse, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	accountToInsert := &Account{
		ID:       id,
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	if !s.store.Insert(accountToInsert) {
		// InvalidAccountExists mornally this workflow would be handled with a status 201 and a message saying an email was sent to
		// verify the account. When this error is hit, an email would be sent saying if you're trying to create
		// an you can trying executing the forgot password workflow instead of leaking internal info to the user
		// that an account already exists with the email provided, but this is outside of the scope of this project
		return nil, common.InvalidAccountExists()
	}

	return &CreateAccountResponse{}, nil
}

func (s *service) SigninAccount(req *SigninAccountRequest) (*SigninAccountResponse, error) {
	accountFound := s.store.GetByEmail(req.Email)
	if accountFound == nil {
		return nil, common.InvalidCredentials()
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, accountFound.Password)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, common.InvalidCredentials()
	}

	return &SigninAccountResponse{}, nil
}
