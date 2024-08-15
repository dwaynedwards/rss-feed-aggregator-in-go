package password_test

import (
	"testing"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/password"
	"github.com/matryer/is"
)

func TestXxx(t *testing.T) {
	is := is.New(t)

	pass := "password1"
	hashedPassword, err := password.Hash(pass)

	is.NoErr(err)

	is.True(pass != hashedPassword)
	match, err := password.Matches(pass, hashedPassword)

	is.NoErr(err)
	is.True(match)
}
