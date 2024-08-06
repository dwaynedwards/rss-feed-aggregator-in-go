package postgres

import (
	"context"

	rssfeed "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/jackc/pgx/v5"
)

type UserStore struct{}

func createUser(ctx context.Context, tx *Tx, user *rssfeed.User) error {
	user.CreatedAt = tx.now
	user.ModifiedAt = user.CreatedAt

	query := `
  INSERT INTO users (name, created_at, modified_at)
  VALUES (@name, @createdAt, @modifiedAt)
  RETURNING id`
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
