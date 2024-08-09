package builder

import (
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
)

type userBuilder struct {
	user *rf.User
}

func NewUserBuilder() *userBuilder {
	return &userBuilder{
		user: &rf.User{},
	}
}

func (b *userBuilder) WithID(id int64) *userBuilder {
	b.user.ID = id
	return b
}

func (b *userBuilder) WithName(name string) *userBuilder {
	b.user.Name = name
	return b
}

func (b *userBuilder) WithCreatedAt(createdAt time.Time) *userBuilder {
	b.user.CreatedAt = createdAt
	return b
}

func (b *userBuilder) WithModifiedAt(modifiedAt time.Time) *userBuilder {
	b.user.ModifiedAt = modifiedAt
	return b
}

func (b *userBuilder) Build() *rf.User {
	return b.user
}
