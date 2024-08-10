package service

import (
	"context"
	"errors"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
)

type FeedService struct {
	store rf.FeedStore
}

func NewFeedService(store rf.FeedStore) *FeedService {
	return &FeedService{
		store: store,
	}
}

func (fs *FeedService) AddFeed(ctx context.Context, feed *rf.Feed) (int64, error) {
	args := FeedArgs{
		store: fs.store,
		feed:  feed,
	}

	if err := args.validateAddFeed(); err != nil {
		return 0, err
	}

	err := rf.RunStateMachine(ctx, args, addFeedState)
	if err != nil {
		return 0, err
	}

	return feed.ID, nil
}

func (fs *FeedService) RemoveFeed(ctx context.Context, feed *rf.Feed) error {
	return errors.New("remove err")
}

func (fs *FeedService) GetFeeds(ctx context.Context, feed *rf.Feed) ([]rf.Feed, error) {
	foundFeeds, err := fs.store.ListUserFeeds(ctx, feed.UserID)
	if err != nil {
		return nil, err
	}
	return foundFeeds, nil
}

func (fs *FeedService) GetFeed(ctx context.Context, feed *rf.Feed) (*rf.Feed, error) {
	foundFeed, err := fs.store.FindUserFeedByID(ctx, feed.UserID, feed.ID)
	if err != nil {
		return nil, err
	}

	return foundFeed, nil
}
