package rssfeed

import "time"

type userBuilder struct {
	user *User
}

func NewUserBuilder() *userBuilder {
	return &userBuilder{
		user: &User{},
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

func (b *userBuilder) Build() *User {
	return b.user
}
