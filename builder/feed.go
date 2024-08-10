package builder

import (
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
)

type feedBuilder struct {
	feed *rf.Feed
}

func NewFeedBuilder() *feedBuilder {
	return &feedBuilder{
		feed: &rf.Feed{},
	}
}

func (b *feedBuilder) WithID(id int64) *feedBuilder {
	b.feed.ID = id
	return b
}

func (b *feedBuilder) WithUserID(userID int64) *feedBuilder {
	b.feed.UserID = userID
	return b
}

func (b *feedBuilder) WithName(name string) *feedBuilder {
	b.feed.Name = name
	return b
}

func (b *feedBuilder) WithURL(url string) *feedBuilder {
	b.feed.URL = url
	return b
}

func (b *feedBuilder) AsEnabled(enabled bool) *feedBuilder {
	b.feed.Enabled = enabled
	return b
}

func (b *feedBuilder) AsDeleted(deleted bool) *feedBuilder {
	b.feed.Deleted = deleted
	return b
}

func (b *feedBuilder) WithCreatedAt(createdAt time.Time) *feedBuilder {
	b.feed.CreatedAt = createdAt
	return b
}

func (b *feedBuilder) WithModifiedAt(modifiedAt time.Time) *feedBuilder {
	b.feed.ModifiedAt = modifiedAt
	return b
}

func (b *feedBuilder) WithLastSyncedAt(lastSyncedAt time.Time) *feedBuilder {
	b.feed.LastSyncedAt = lastSyncedAt
	return b
}

func (b *feedBuilder) Build() *rf.Feed {
	return b.feed
}
