package builder

import (
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
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

type signUpRequestBuilder struct {
	signUpRequest *rf.SignUpRequest
}

func NewSignUpRequestBuilder() *signUpRequestBuilder {
	return &signUpRequestBuilder{
		signUpRequest: &rf.SignUpRequest{},
	}
}

func (b *signUpRequestBuilder) WithEmail(email string) *signUpRequestBuilder {
	b.signUpRequest.Email = email
	return b
}

func (b *signUpRequestBuilder) WithPassword(password string) *signUpRequestBuilder {
	b.signUpRequest.Password = password
	return b
}

func (b *signUpRequestBuilder) WithName(name string) *signUpRequestBuilder {
	b.signUpRequest.Name = name
	return b
}

func (b *signUpRequestBuilder) Build() *rf.SignUpRequest {
	return b.signUpRequest
}

type signInRequestBuilder struct {
	signInRequest *rf.SignInRequest
}

func NewSignInRequestBuilder() *signInRequestBuilder {
	return &signInRequestBuilder{
		signInRequest: &rf.SignInRequest{},
	}
}

func (b *signInRequestBuilder) WithEmail(email string) *signInRequestBuilder {
	b.signInRequest.Email = email
	return b
}

func (b *signInRequestBuilder) WithPassword(password string) *signInRequestBuilder {
	b.signInRequest.Password = password
	return b
}

func (b *signInRequestBuilder) Build() *rf.SignInRequest {
	return b.signInRequest
}
