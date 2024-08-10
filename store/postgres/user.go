package postgres

import (
	"context"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/jackc/pgx/v5"
)

type UserStore struct {
	db *DB
}

func NewUserStore(db *DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func createUser(ctx context.Context, tx *Tx, user *rf.User) error {
	user.CreatedAt = tx.now
	user.ModifiedAt = user.CreatedAt

	query := `
  INSERT INTO users (name, created_at, modified_at)
  VALUES (@name, @createdAt, @modifiedAt)
  RETURNING id
	`
	args := pgx.NamedArgs{
		"name":       user.Name,
		"createdAt":  user.CreatedAt,
		"modifiedAt": user.ModifiedAt,
	}

	err := tx.QueryRow(ctx, query, args).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}
