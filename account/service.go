package account

import (
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
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
	password, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, &common.AccountError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	account := &Account{
		Email:    req.Email,
		Password: password,
		Name:     req.Name,
	}

	if err := s.store.Create(account); err != nil {
		return nil, err
	}

	return &CreateAccountResponse{}, nil
}

func (s *service) SigninAccount(req *SigninAccountRequest) (*SigninAccountResponse, error) {
	account := &Account{
		Email: req.Email,
	}

	authenticatedAccount, err := s.store.Signin(account)
	if err != nil {
		return nil, err
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, authenticatedAccount.Password)
	if err != nil {
		return nil, &common.AccountError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	if !match {
		return nil, &common.AccountError{Status: http.StatusUnauthorized, Msg: "incorrect email or password provided"}
	}

	return &SigninAccountResponse{
		Token: authenticatedAccount.ID.String(),
	}, nil
}
