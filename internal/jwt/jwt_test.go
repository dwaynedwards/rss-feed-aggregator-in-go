package jwt_test

import (
	"testing"
	"time"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/jwt"
	"github.com/matryer/is"
)

func TestJWT(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	token, err := jwt.GenerateAndSignUserID(1, time.Now())
	is.NoErr(err)
	is.True(len(token) > 0)

	userID, err := jwt.ParseAndVerifyUserID(token)
	is.NoErr(err)
	is.Equal(userID, int64(userID))

	token, err = jwt.GenerateAndSignUserID(1, time.Now().Add(time.Minute*-1))
	is.NoErr(err)
	is.True(len(token) > 0)

	_, err = jwt.ParseAndVerifyUserID(token)
	is.True(err != nil)                                        // should error from token expiration
	is.Equal(errors.ToReferenceCode(err), errors.Unauthorized) // shoud have error code
	is.Equal(errors.ToErr(err), errors.ErrTokenExpired)        // shoud have token expired error
}
