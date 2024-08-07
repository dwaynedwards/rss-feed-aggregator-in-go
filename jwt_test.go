package rf_test

import (
	"testing"
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/matryer/is"
)

func TestJWT(t *testing.T) {
	is := is.New(t)

	token, err := rf.GenerateAndSignJWT(1, time.Now())
	is.NoErr(err)
	is.True(len(token) > 0)

	userID, err := rf.ParseAndVerifyJWT(token)
	is.NoErr(err)
	is.Equal(userID, int64(userID))

	token, err = rf.GenerateAndSignJWT(1, time.Now().Add(time.Minute*-1))
	is.NoErr(err)
	is.True(len(token) > 0)

	_, err = rf.ParseAndVerifyJWT(token)
	is.True(err != nil)                               // should error from token expiration
	is.Equal(rf.AppErrorCode(err), rf.ECUnautherized) // shoud have error code
}
