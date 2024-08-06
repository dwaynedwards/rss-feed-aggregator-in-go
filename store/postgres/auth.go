package postgres

import (
	"context"
	"errors"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/jackc/pgx/v5"
)

type AuthStore struct {
	db *DB
}

func NewAuthStore(db *DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

func (a *AuthStore) Create(ctx context.Context, auth *rf.Auth) error {
	tx, err := a.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	authFound, err := a.FindByEmail(ctx, auth.BasicAuth.Email)
	if err != nil {
		return err
	}

	if authFound != nil {
		return rf.AppErrorf(rf.ECInvalid, rf.EMUserExists)
	}

	user := &rf.User{
		Name: auth.User.Name,
	}
	err = createUser(ctx, tx, user)
	if err != nil {
		return err
	}

	auth.UserID = user.ID
	auth.CreatedAt = tx.now
	auth.ModifiedAt = auth.CreatedAt
	auth.LastSignedInAt = auth.CreatedAt

	query := `
	INSERT INTO auths (user_id, email, password, created_at, modified_at, last_signed_in_at)
	VALUES (@userID, @email, @password, @createdAt, @modifiedAt, @lastSignedInAt) RETURNING id`
	args := pgx.NamedArgs{
		"userID":         auth.UserID,
		"email":          auth.BasicAuth.Email,
		"password":       auth.BasicAuth.Password,
		"createdAt":      auth.CreatedAt,
		"modifiedAt":     auth.ModifiedAt,
		"lastSignedInAt": auth.LastSignedInAt,
	}

	err = tx.QueryRow(ctx, query, args).Scan(&auth.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (a *AuthStore) FindByEmail(ctx context.Context, email string) (*rf.Auth, error) {
	tx, err := a.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	return findByEmail(ctx, tx, email)
}

func findByEmail(ctx context.Context, tx *Tx, email string) (*rf.Auth, error) {
	auth := &rf.Auth{
		BasicAuth: &rf.BasicAuth{
			Email: email,
		},
	}

	query := `
	SELECT user_id, password FROM auths WHERE email = @email
	`
	args := pgx.NamedArgs{
		"email": email,
	}

	err := tx.QueryRow(ctx, query, args).Scan(&auth.UserID, &auth.BasicAuth.Password)
	if err != nil {
		if ok := errors.Is(err, pgx.ErrNoRows); ok {
			return nil, nil
		}
		return nil, err
	}

	return auth, nil
}
