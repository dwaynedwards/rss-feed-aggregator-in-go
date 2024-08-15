package feedservice

import (
	"context"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/builder"
	rfcontext "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/context"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/statemachine"
)

type FeedStore interface {
	CreateFeed(ctx context.Context, feed *rf.Feed) error
	CreateUserFeed(ctx context.Context, feed *rf.Feed) error
	ListUserFeeds(ctx context.Context, userID int64) ([]rf.Feed, error)
	FindUserFeedByID(ctx context.Context, userID, feedID int64) (*rf.Feed, error)
	FindByURL(ctx context.Context, url string) (*rf.Feed, error)
	DeleteFeed(ctx context.Context, userID, feedID int64) error
}

type FeedService struct {
	store FeedStore
}

func NewFeedService(store FeedStore) *FeedService {
	return &FeedService{
		store: store,
	}
}

func (fs *FeedService) AddFeed(ctx context.Context, req *rf.AddFeedRequest) (int64, error) {
	userID := rfcontext.UserIDFromContext(ctx)

	feed := builder.NewFeedBuilder().
		WithName(req.Name).
		WithURL(req.URL).
		WithUserID(userID).
		Build()

	args := FeedArgs{
		store: fs.store,
		feed:  feed,
	}

	if err := args.validateAddFeed(); err != nil {
		return 0, err
	}

	result, err := statemachine.Run(ctx, args, addFeedState)
	if err != nil {
		return 0, err
	}

	return result.feed.ID, nil
}

func (fs *FeedService) RemoveFeed(ctx context.Context, feedID int64) error {
	userID := rfcontext.UserIDFromContext(ctx)

	err := fs.store.DeleteFeed(ctx, userID, feedID)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FeedService) GetFeeds(ctx context.Context) ([]rf.Feed, error) {
	userID := rfcontext.UserIDFromContext(ctx)

	foundFeeds, err := fs.store.ListUserFeeds(ctx, userID)
	if err != nil {
		return nil, err
	}
	return foundFeeds, nil
}

func (fs *FeedService) GetFeed(ctx context.Context, feedID int64) (*rf.Feed, error) {
	userID := rfcontext.UserIDFromContext(ctx)

	foundFeed, err := fs.store.FindUserFeedByID(ctx, userID, feedID)
	if err != nil {
		return nil, err
	}

	return foundFeed, nil
}
