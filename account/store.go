package account

import (
	"net/http"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
)

type store struct {
	db inMemoryAccountDB
}

func NewStore() *store {
	return &store{
		db: make(inMemoryAccountDB),
	}
}

func (s *store) Create(account *Account) error {
	_, ok := s.db[account.Email]
	if ok {
		return &common.AccountError{Status: http.StatusConflict, Msg: "account already exists"}
	}

	s.db[account.Email] = *account

	return nil
}

func (s *store) Signin(account *Account) (*Account, error) {
	foundAccount, ok := s.db[account.Email]
	if !ok {
		return nil, &common.AccountError{Status: http.StatusUnauthorized, Msg: "incorrect email or password provided"}
	}

	return &Account{
		ID:       foundAccount.ID,
		Password: foundAccount.Password,
	}, nil
}
