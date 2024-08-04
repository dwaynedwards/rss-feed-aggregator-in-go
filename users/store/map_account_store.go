package store

import (
	"sync"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/users"
	"github.com/google/uuid"
)

type optionsFunc func(c *mapUsersStore) error

type mapUsersStore struct {
	mu sync.RWMutex
	db map[uuid.UUID]*users.User
}

func NewMapUserStore(opts ...optionsFunc) (*mapUsersStore, error) {
	store := &mapUsersStore{
		db: map[uuid.UUID]*users.User{},
	}

	for _, opt := range opts {
		if err := opt(store); err != nil {
			return nil, err
		}
	}

	return store, nil
}

func WithDB(db map[uuid.UUID]*users.User) optionsFunc {
	return func(c *mapUsersStore) error {
		c.db = db
		return nil
	}
}

func (s *mapUsersStore) InsertUser(acc *users.User) bool {
	foundUser := s.GetUserByEmail(acc.Email)

	if foundUser != nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[acc.ID] = acc

	return true
}

func (s *mapUsersStore) GetUserByID(id uuid.UUID) *users.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	foundUser, ok := s.db[id]
	if !ok {
		return nil
	}

	return foundUser
}

func (s *mapUsersStore) GetUserByEmail(email string) *users.User {
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
