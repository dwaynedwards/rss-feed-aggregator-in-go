package builder

import (
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
)

type authBuilder struct {
	auth *rf.Auth
}

func NewAuthBuilder() *authBuilder {
	return &authBuilder{
		auth: &rf.Auth{},
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

func (b *authBuilder) AsEnabled(enabled bool) *authBuilder {
	b.auth.Enabled = enabled
	return b
}

func (b *authBuilder) AsDeleted(deleted bool) *authBuilder {
	b.auth.Deleted = deleted
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

func (b *authBuilder) Build() *rf.Auth {
	return b.auth
}

type basicAuthBuilder struct {
	basicAuth *rf.BasicAuth
}

func NewBasicAuthBuilder() *basicAuthBuilder {
	return &basicAuthBuilder{
		basicAuth: &rf.BasicAuth{},
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

func (b *basicAuthBuilder) Build() *rf.BasicAuth {
	return b.basicAuth
}

type signUpAuthRequestBuilder struct {
	signUpAuthRequest *rf.SignUpAuthRequest
}

func NewSignUpAuthRequestBuilder() *signUpAuthRequestBuilder {
	return &signUpAuthRequestBuilder{
		signUpAuthRequest: &rf.SignUpAuthRequest{},
	}
}

func (b *signUpAuthRequestBuilder) WithEmail(email string) *signUpAuthRequestBuilder {
	b.signUpAuthRequest.Email = email
	return b
}

func (b *signUpAuthRequestBuilder) WithPassword(password string) *signUpAuthRequestBuilder {
	b.signUpAuthRequest.Password = password
	return b
}

func (b *signUpAuthRequestBuilder) WithName(name string) *signUpAuthRequestBuilder {
	b.signUpAuthRequest.Name = name
	return b
}

func (b *signUpAuthRequestBuilder) Build() *rf.SignUpAuthRequest {
	return b.signUpAuthRequest
}

type signInAuthRequestBuilder struct {
	signInAuthRequest *rf.SignInAuthRequest
}

func NewSignInAuthRequestBuilder() *signInAuthRequestBuilder {
	return &signInAuthRequestBuilder{
		signInAuthRequest: &rf.SignInAuthRequest{},
	}
}

func (b *signInAuthRequestBuilder) WithEmail(email string) *signInAuthRequestBuilder {
	b.signInAuthRequest.Email = email
	return b
}

func (b *signInAuthRequestBuilder) WithPassword(password string) *signInAuthRequestBuilder {
	b.signInAuthRequest.Password = password
	return b
}

func (b *signInAuthRequestBuilder) Build() *rf.SignInAuthRequest {
	return b.signInAuthRequest
}
