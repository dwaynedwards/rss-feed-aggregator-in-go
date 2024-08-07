package rf

import "time"

type authBuilder struct {
	auth *Auth
}

func NewAuthBuilder() *authBuilder {
	return &authBuilder{
		auth: &Auth{},
	}
}

func (b *authBuilder) WithID(id int64) *authBuilder {
	b.auth.ID = id
	return b
}

func (b *authBuilder) WithUserID(userID int64) *authBuilder {
	b.auth.UserID = userID
	return b
}

func (b *authBuilder) WithUser(builder *userBuilder) *authBuilder {
	b.auth.User = builder.Build()
	return b
}

func (b *authBuilder) WithBasicAuth(builder *basicAuthBuilder) *authBuilder {
	b.auth.BasicAuth = builder.Build()
	return b
}

func (b *authBuilder) WithCreatedAt(createdAt time.Time) *authBuilder {
	b.auth.CreatedAt = createdAt
	return b
}

func (b *authBuilder) WithModifiedAt(modifiedAt time.Time) *authBuilder {
	b.auth.ModifiedAt = modifiedAt
	return b
}

func (b *authBuilder) WithLastSignedInAt(lastSignedInAt time.Time) *authBuilder {
	b.auth.LastSignedInAt = lastSignedInAt
	return b
}

func (b *authBuilder) Build() *Auth {
	return b.auth
}

type basicAuthBuilder struct {
	basicAuth *BasicAuth
}

func NewBasicAuthBuilder() *basicAuthBuilder {
	return &basicAuthBuilder{
		basicAuth: &BasicAuth{},
	}
}

func (b *basicAuthBuilder) WithEmail(email string) *basicAuthBuilder {
	b.basicAuth.Email = email
	return b
}

func (b *basicAuthBuilder) WithPassword(password string) *basicAuthBuilder {
	b.basicAuth.Password = password
	return b
}

func (b *basicAuthBuilder) Build() *BasicAuth {
	return b.basicAuth
}
