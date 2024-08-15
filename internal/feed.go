package rf

import (
	"time"
)

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
