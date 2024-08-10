package postgres_test

import (
	"context"
	"testing"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	rfbuilder "github.com/dwaynedwards/rss-feed-aggregator-in-go/builder"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/service"
	rfpg "github.com/dwaynedwards/rss-feed-aggregator-in-go/store/postgres"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/testcontainers"
	"github.com/matryer/is"
)

func TestPostgresDBFeedServiceIntegration(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	ctx := context.Background()

	container, err := testcontainers.NewPostgresTestContainer(ctx)
	is.NoErr(err)

	t.Cleanup(func() {
		err := container.Cleanup(ctx)
		is.NoErr(err) // failed to terminate pgContainer
	})

	authStore := rfpg.NewAuthStore(container.DB)
	authService := service.NewAuthService(authStore)

	feedStore := rfpg.NewFeedStore(container.DB)
	feedService := service.NewFeedService(feedStore)

	auth := rfbuilder.NewAuthBuilder().
		WithUser(rfbuilder.NewUserBuilder().WithName("Gopher")).
		WithBasicAuth(rfbuilder.NewBasicAuthBuilder().
			WithEmail("gopher1@go.com").
			WithPassword("gogopher1")).
		Build()

	token, err := authService.SignUp(ctx, auth)
	is.NoErr(err)           // should sign up
	is.True(len(token) > 0) // should receive token

	userID := auth.UserID

	feedGetFeedsWithNoFeedsSuccess := &rf.Feed{
		UserID: userID,
	}
	feeds, err := feedService.GetFeeds(ctx, feedGetFeedsWithNoFeedsSuccess)

	is.NoErr(err)           // should have no errors when no feeds are found
	is.Equal(len(feeds), 0) // should find no feeds

	feedName := "The Gopher Podcast"
	feedURL := "http://feed.com/rss"
	feedAddSuccess := rfbuilder.NewFeedBuilder().
		WithUserID(userID).
		WithName(feedName).
		WithURL(feedURL).
		Build()

	feedID, err := feedService.AddFeed(ctx, feedAddSuccess)

	is.NoErr(err)
	is.Equal(feedID, int64(1))

	feedGetFeedSuccess := &rf.Feed{
		ID:     feedID,
		UserID: userID,
	}
	feed, err := feedService.GetFeed(ctx, feedGetFeedSuccess)

	is.NoErr(err)                 // should find a feed
	is.Equal(feed.UserID, userID) // should have the user id used to find it
	is.Equal(feed.ID, feedID)     // should have the feed id used to find it
	is.Equal(feed.Name, feedName) // should have feed name that was crated

	feedGetFeedsWith1FeedSuccess := &rf.Feed{
		UserID: userID,
	}
	feeds, err = feedService.GetFeeds(ctx, feedGetFeedsWith1FeedSuccess)

	is.NoErr(err)                     // should find feeds
	is.Equal(len(feeds), 1)           // should find 1 feed
	is.Equal(feeds[0].ID, feedID)     // should have feed id
	is.Equal(feeds[0].UserID, userID) // should have user id
	is.Equal(feeds[0].Name, feedName) // should have feed name
	is.Equal(feeds[0].URL, feedURL)   // should have feed url

	invalidFeedID := int64(100)
	feedGetFeedInvalidFeedIDFailure := &rf.Feed{
		ID:     invalidFeedID,
		UserID: userID,
	}
	feed, err = feedService.GetFeed(ctx, feedGetFeedInvalidFeedIDFailure)

	is.NoErr(err)        // should get no error when no feed is found
	is.True(feed == nil) // should not find feed

	invalidUserID := int64(100)
	feedGetFeedInvalidUserIDFailure := &rf.Feed{
		ID:     feedID,
		UserID: invalidUserID,
	}
	feed, err = feedService.GetFeed(ctx, feedGetFeedInvalidUserIDFailure)

	is.NoErr(err)        // should get no error when no feed is found
	is.True(feed == nil) // should not find feed
}
