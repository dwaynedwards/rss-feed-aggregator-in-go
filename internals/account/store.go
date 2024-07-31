package account

import (
	"net/http"
)

type store struct {
	db inMemoryAccountDB
}

func NewAccountStore() *store {
	return &store{
		db: make(inMemoryAccountDB),
	}
}

func (s *store) Create(account *Account) error {
	_, ok := s.db[account.Email]
	if ok {
		return &AccountError{Status: http.StatusConflict, Msg: "User already exists"}
	}

	s.db[account.Email] = *account

	return nil
}
