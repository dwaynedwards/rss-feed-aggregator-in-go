package postgres

import (
	"context"
	"errors"
	"fmt"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/jackc/pgx/v5"
)

type FeedStore struct {
	db *DB
}

func NewFeedStore(db *DB) *FeedStore {
	return &FeedStore{
		db: db,
	}
}

func (fs *FeedStore) CreateFeed(ctx context.Context, feed *rf.Feed) error {
	tx, err := fs.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	feed.CreatedAt = tx.now
	feed.ModifiedAt = feed.CreatedAt
	feed.LastSyncedAt = feed.CreatedAt

	query := `
	INSERT INTO feeds (url, created_at, modified_at, last_synced_at)
	VALUES (@url, @createdAt, @modifiedAt, @lastSyncedAt)
	RETURNING id
	`
	args := pgx.NamedArgs{
		"url":          feed.URL,
		"createdAt":    feed.CreatedAt,
		"modifiedAt":   feed.ModifiedAt,
		"lastSyncedAt": feed.LastSyncedAt,
	}

	err = tx.QueryRow(ctx, query, args).Scan(&feed.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (fs *FeedStore) CreateUserFeed(ctx context.Context, feed *rf.Feed) error {
	tx, err := fs.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
	INSERT INTO user_feeds (user_id, feed_id, name, created_at, modified_at)
	VALUES (@userID, @feedID, @name, @createdAt, @modifiedAt)
	`
	args := pgx.NamedArgs{
		"userID":     feed.UserID,
		"feedID":     feed.ID,
		"name":       feed.Name,
		"createdAt":  feed.CreatedAt,
		"modifiedAt": feed.ModifiedAt,
	}

	result, err := tx.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if result.RowsAffected() != 1 {
		return rf.NewAppError(rf.ECIntenal, fmt.Sprintf("no rows were inserted into user_feeds: %s", result.String()))
	}

	return tx.Commit(ctx)
}

func (fs *FeedStore) ListUserFeeds(ctx context.Context, userID int64) ([]rf.Feed, error) {
	tx, err := fs.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
	SELECT user_feeds.user_id as user_id,
				 user_feeds.feed_id as feed_id,
				 user_feeds.name as name,
				 feeds.url as url
		FROM user_feeds
		LEFT JOIN feeds
			ON user_feeds.feed_id = feeds.id
		WHERE user_id = @userID
	`
	args := pgx.NamedArgs{
		"userID": userID,
	}

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	feeds, err := pgx.CollectRows(rows, pgx.RowToStructByName[rf.Feed])
	if err != nil {
		return nil, err
	}

	return feeds, nil
}

func (fs *FeedStore) FindUserFeedByID(ctx context.Context, userID, feedID int64) (*rf.Feed, error) {
	tx, err := fs.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	feed := &rf.Feed{
		ID:     feedID,
		UserID: userID,
	}

	query := `
	SELECT name
	FROM user_feeds
	WHERE user_id = @userID AND feed_id = @feedID
	`
	args := pgx.NamedArgs{
		"userID": feed.UserID,
		"feedID": feed.ID,
	}

	err = tx.QueryRow(ctx, query, args).Scan(&feed.Name)
	if err != nil {
		if ok := errors.Is(err, pgx.ErrNoRows); ok {
			return nil, nil
		}
		return nil, err
	}

	return feed, nil
}

func (fs *FeedStore) FindByURL(ctx context.Context, url string) (*rf.Feed, error) {
	return nil, nil
}

func (fs *FeedStore) DeleteFeed(ctx context.Context, userID, feedID int64) error {
	return nil
}
