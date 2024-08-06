package store

import (
	"sync"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/users"
)

type MapUsersDB map[int64]*users.User

type optionsFunc func(c *mapUsersStore) error

type mapUsersStore struct {
	mu sync.RWMutex
	db MapUsersDB
}

func NewMapUserStore(opts ...optionsFunc) (*mapUsersStore, error) {
	store := &mapUsersStore{
		db: MapUsersDB{},
	}

	for _, opt := range opts {
		if err := opt(store); err != nil {
			return nil, err
		}
	}

	return store, nil
}

func WithDB(db MapUsersDB) optionsFunc {
	return func(c *mapUsersStore) error {
		c.db = db
		return nil
	}
}

func (s *mapUsersStore) InsertUser(user *users.User) error {
	foundUser := s.getUserByEmail(user.Email)

	if foundUser != nil {
		// InvalidUserExists mornally this workflow would be handled with a status 201 and a message saying an email was sent to
		// verify the user. When this error is hit, an email would be sent saying if you're trying to create
		// an you can trying executing the forgot password workflow instead of leaking internal info to the user
		// that an user already exists with the email provided, but this is outside of the scope of this project
		return common.InvalidUserExists()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = int64(len(s.db)) + 1

	s.db[user.ID] = user

	return nil
}

func (s *mapUsersStore) GetUserByID(id int64) (*users.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	foundUser, ok := s.db[id]
	if !ok {
		return nil, nil
	}

	return foundUser, nil
}

func (s *mapUsersStore) GetUserByEmail(email string) (*users.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getUserByEmail(email), nil
}

func (s *mapUsersStore) getUserByEmail(email string) *users.User {
	for _, user := range s.db {
		if user.Email != email {
			continue
		}
		return user
	}

	return nil
}
