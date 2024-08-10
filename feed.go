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
}

type FeedService interface {
	AddFeed(ctx context.Context, feed *Feed) (int64, error)
	RemoveFeed(ctx context.Context, feed *Feed) error
	GetFeeds(ctx context.Context, feed *Feed) ([]Feed, error)
	GetFeed(ctx context.Context, feed *Feed) (*Feed, error)
}

type Feed struct {
	ID           int64     `json:"id" db:"feed_id"`
	Name         string    `json:"name" db:"name"`
	URL          string    `json:"url" db:"url"`
	Enabled      bool      `json:"enabled" db:"-"`
	Deleted      bool      `json:"deleted" db:"-"`
	CreatedAt    time.Time `json:"createdAt" db:"-"`
	ModifiedAt   time.Time `json:"modifiedAt" db:"-"`
	LastSyncedAt time.Time `json:"lastSyncedAt" db:"-"`

	UserID int64 `json:"userID" db:"user_id"`
}
