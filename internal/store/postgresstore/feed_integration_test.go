package postgresstore_test

import (
	"context"
	"testing"

	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/builder"
	rfcontext "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/context"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/service/authservice"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/service/feedservice"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/store/postgresstore"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/testcontainers"
	"github.com/matryer/is"
)

func TestPostgresDBFeedServiceIntegration(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	ctx := context.Background()

	container, err := testcontainers.NewPostgres(ctx)
	is.NoErr(err)

	migration, err := postgresstore.NewPostgresMigration(container.DB, "migrations")
	is.NoErr(err)

	migration.Up()
	is.NoErr(err)

	t.Cleanup(func() {
		err := migration.Reset()
		is.NoErr(err)
		err = migration.Close()
		is.NoErr(err)
		err = container.Cleanup(ctx)
		is.NoErr(err) // failed to terminate pgContainer
	})

	authStore := postgresstore.NewAuthStore(container.DB)
	authService := authservice.NewAuthService(authStore)

	feedStore := postgresstore.NewFeedStore(container.DB)
	feedService := feedservice.NewFeedService(feedStore)

	signUpReq := builder.NewSignUpRequestBuilder().
		WithName("Gopher").
		WithEmail("gopher1@go.com").
		WithPassword("gogopher1").
		Build()

	token, err := authService.SignUp(ctx, signUpReq)
	is.NoErr(err)           // should sign up
	is.True(len(token) > 0) // should receive token

	userID := int64(1)

	ctxWithUserID := rfcontext.SetUserIDToContext(ctx, userID)
	feeds, err := feedService.GetFeeds(ctxWithUserID)

	is.NoErr(err)           // should have no errors when no feeds are found
	is.Equal(len(feeds), 0) // should find no feeds

	feedName := "The Gopher Podcast"
	feedURL := "http://feed.com/rss"
	feedAddSuccess := builder.NewAddFeedBuilder().
		WithName(feedName).
		WithURL(feedURL).
		Build()

	feedID, err := feedService.AddFeed(ctxWithUserID, feedAddSuccess)

	is.NoErr(err)
	is.Equal(feedID, int64(1))

	feed, err := feedService.GetFeed(ctxWithUserID, feedID)

	is.NoErr(err)                 // should find a feed
	is.Equal(feed.UserID, userID) // should have the user id used to find it
	is.Equal(feed.ID, feedID)     // should have the feed id used to find it
	is.Equal(feed.Name, feedName) // should have feed name that was crated

	feeds, err = feedService.GetFeeds(ctxWithUserID)

	is.NoErr(err)                     // should find feeds
	is.Equal(len(feeds), 1)           // should find 1 feed
	is.Equal(feeds[0].ID, feedID)     // should have feed id
	is.Equal(feeds[0].UserID, userID) // should have user id
	is.Equal(feeds[0].Name, feedName) // should have feed name
	is.Equal(feeds[0].URL, feedURL)   // should have feed url

	invalidFeedID := int64(100)
	feed, err = feedService.GetFeed(ctxWithUserID, invalidFeedID)

	is.NoErr(err)        // should get no error when no feed is found
	is.True(feed == nil) // should not find feed

	invalidUserID := int64(100)
	ctxWithInvalidUserID := rfcontext.SetUserIDToContext(ctx, invalidUserID)
	feed, err = feedService.GetFeed(ctxWithInvalidUserID, feedID)

	is.NoErr(err)        // should get no error when no feed is found
	is.True(feed == nil) // should not find feed
}
