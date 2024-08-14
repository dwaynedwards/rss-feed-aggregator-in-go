package rf

import (
	"context"
	"time"
)

type FeedStore interface {
	CreateFeed(ctx context.Context, feed *Feed) error
	CreateUserFeed(ctx context.Context, feed *Feed) error
	ListUserFeeds(ctx context.Context, userID int64) ([]Feed, error)
	FindUserFeedByID(ctx context.Context, userID, feedID int64) (*Feed, error)
	FindByURL(ctx context.Context, url string) (*Feed, error)
	DeleteFeed(ctx context.Context, userID, feedID int64) error
}

type FeedService interface {
	AddFeed(ctx context.Context, req *AddFeedRequest) (int64, error)
	RemoveFeed(ctx context.Context, feedID int64) error
	GetFeeds(ctx context.Context) ([]Feed, error)
	GetFeed(ctx context.Context, feedID int64) (*Feed, error)
}

type Feed struct {
	ID           int64     `db:"feed_id"`
	Name         string    `db:"name"`
	URL          string    `db:"url"`
	Enabled      bool      `db:"-"`
	Deleted      bool      `db:"-"`
	CreatedAt    time.Time `db:"-"`
	ModifiedAt   time.Time `db:"-"`
	LastSyncedAt time.Time `db:"-"`

	UserID int64 `db:"user_id"`
}

type AddFeedRequest struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	UserID int64
}
