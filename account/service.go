package account

import (
	"net/http"

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
		return nil, &common.AccountError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, &common.AccountError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	accountToInsert := &Account{
		ID:       id,
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	if !s.store.Insert(accountToInsert) {
		return nil, &common.AccountError{Status: http.StatusConflict, Msg: "account already exists"}
	}

	return &CreateAccountResponse{}, nil
}

func (s *service) SigninAccount(req *SigninAccountRequest) (*SigninAccountResponse, error) {
	accountFound := s.store.GetByEmail(req.Email)
	if accountFound == nil {
		return nil, &common.AccountError{Status: http.StatusUnauthorized, Msg: "incorrect email or password provided"}
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, accountFound.Password)
	if err != nil {
		return nil, &common.AccountError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	if !match {
		return nil, &common.AccountError{Status: http.StatusUnauthorized, Msg: "incorrect email or password provided"}
	}

	return &SigninAccountResponse{}, nil
}
