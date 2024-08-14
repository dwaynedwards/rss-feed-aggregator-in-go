package service

import (
	"context"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/builder"
)

type FeedService struct {
	store rf.FeedStore
}

func NewFeedService(store rf.FeedStore) *FeedService {
	return &FeedService{
		store: store,
	}
}

func (fs *FeedService) AddFeed(ctx context.Context, req *rf.AddFeedRequest) (int64, error) {
	userID := rf.UserIDFromContext(ctx)

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

	result, err := rf.RunStateMachine(ctx, args, addFeedState)
	if err != nil {
		return 0, err
	}

	return result.feed.ID, nil
}

func (fs *FeedService) RemoveFeed(ctx context.Context, feedID int64) error {
	userID := rf.UserIDFromContext(ctx)

	err := fs.store.DeleteFeed(ctx, userID, feedID)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FeedService) GetFeeds(ctx context.Context) ([]rf.Feed, error) {
	userID := rf.UserIDFromContext(ctx)

	foundFeeds, err := fs.store.ListUserFeeds(ctx, userID)
	if err != nil {
		return nil, err
	}
	return foundFeeds, nil
}

func (fs *FeedService) GetFeed(ctx context.Context, feedID int64) (*rf.Feed, error) {
	userID := rf.UserIDFromContext(ctx)

	foundFeed, err := fs.store.FindUserFeedByID(ctx, userID, feedID)
	if err != nil {
		return nil, err
	}

	return foundFeed, nil
}
