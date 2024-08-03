package store

import (
	"sync"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/account"
	"github.com/google/uuid"
)

type optionsFunc func(c *mapAccountStore) error

type mapAccountStore struct {
	mu sync.RWMutex
	db map[uuid.UUID]*account.Account
}

func NewMapAccountStore(opts ...optionsFunc) (*mapAccountStore, error) {
	store := &mapAccountStore{
		db: map[uuid.UUID]*account.Account{},
	}

	for _, opt := range opts {
		if err := opt(store); err != nil {
			return nil, err
		}
	}

	return store, nil
}

func WithDB(db map[uuid.UUID]*account.Account) optionsFunc {
	return func(c *mapAccountStore) error {
		c.db = db
		return nil
	}
}

func (s *mapAccountStore) Insert(acc *account.Account) bool {
	foundAccount := s.GetByEmail(acc.Email)

	if foundAccount != nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[acc.ID] = acc

	return true
}

func (s *mapAccountStore) GetByID(id uuid.UUID) *account.Account {
	s.mu.RLock()
	defer s.mu.RUnlock()

	foundAccount, ok := s.db[id]
	if !ok {
		return nil
	}

	return foundAccount
}

func (s *mapAccountStore) GetByEmail(email string) *account.Account {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, acc := range s.db {
		if acc.Email != email {
			continue
		}
		return acc
	}

	return nil
}
