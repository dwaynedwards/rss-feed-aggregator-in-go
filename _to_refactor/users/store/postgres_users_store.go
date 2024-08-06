package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/common"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/users"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type postgresUsersStore struct {
	db *pgx.Conn
}

func NewPostgresUsersStore() (*postgresUsersStore, func(), error) {
	db, err := pgx.Connect(context.Background(), common.GetEnvVar("DATABASE_URL"))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	cleanup := func() { db.Close(context.Background()) }

	if err := db.Ping(context.Background()); err != nil {
		return nil, cleanup, err
	}

	return &postgresUsersStore{
		db: db,
	}, cleanup, nil
}

func (s *postgresUsersStore) InsertUser(user *users.User) error {
	// Doing a check to prevent an id increment on a failed insert due to the unique email constraint
	foundUser, _ := s.getUserByEmail(user.Email)
	if foundUser != nil {
		// InvalidUserExists mornally this workflow would be handled with a status 201 and a message saying an email was sent to
		// verify the user. When this error is hit, an email would be sent saying if you're trying to create
		// an you can trying executing the forgot password workflow instead of leaking internal info to the user
		// that an user already exists with the email provided, but this is outside of the scope of this project
		return common.InvalidUserExists()
	}

	query := `
	INSERT INTO users (email, password, name)
	VALUES ($1, $2, $3)
	`

	_, err := s.db.Exec(context.Background(), query, user.Email, user.Password, user.Name)
	if err != nil {
		var pg *pgconn.PgError
		if ok := errors.As(err, &pg); ok {
			switch pg.Code {
			case pgerrcode.UniqueViolation:
				// InvalidUserExists mornally this workflow would be handled with a status 201 and a message saying an email was sent to
				// verify the user. When this error is hit, an email would be sent saying if you're trying to create
				// an you can trying executing the forgot password workflow instead of leaking internal info to the user
				// that an user already exists with the email provided, but this is outside of the scope of this project
				return common.InvalidUserExists()
			default:
				return err
			}
		}
		return err
	}

	return nil
}

func (s *postgresUsersStore) GetUserByID(id int64) (*users.User, error) {
	query := `SELECT email, password, name FROM users WHERE id = $1`

	foundUser := &users.User{
		ID: id,
	}

	row := s.db.QueryRow(context.Background(), query, id)
	err := row.Scan(&foundUser.Email, &foundUser.Password, &foundUser.Name)
	if err != nil {
		if ok := errors.Is(err, pgx.ErrNoRows); ok {
			return nil, nil
		}
		return nil, err
	}

	return foundUser, nil
}

func (s *postgresUsersStore) GetUserByEmail(email string) (*users.User, error) {
	return s.getUserByEmail(email)
}

func (s *postgresUsersStore) getUserByEmail(email string) (*users.User, error) {
	query := `SELECT id, password, name FROM users WHERE email = $1`

	foundUser := &users.User{
		Email: email,
	}

	row := s.db.QueryRow(context.Background(), query, email)
	err := row.Scan(&foundUser.ID, &foundUser.Password, &foundUser.Name)
	if err != nil {
		if ok := errors.Is(err, pgx.ErrNoRows); ok {
			return nil, nil
		}
		return nil, err
	}

	return foundUser, nil
}
