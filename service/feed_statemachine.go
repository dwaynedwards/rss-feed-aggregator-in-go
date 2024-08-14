package service

import (
	"context"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
)

type FeedArgs struct {
	store rf.FeedStore
	feed  *rf.Feed
}

func (fs FeedArgs) validateAddFeed() error {
	if fs.store == nil {
		return rf.NewAppError(rf.ECIntenal, "store cannot be nil")
	}

	if fs.feed == nil || fs.feed.URL == "" {
		return rf.NewAppError(rf.ECInvalid, rf.EMURLRequired)
	}

	return nil
}

func addFeedState(ctx context.Context, args FeedArgs) (FeedArgs, rf.StateFn[FeedArgs], error) {
	hasFeed, err := args.store.FindByURL(ctx, args.feed.URL)
	if err != nil {
		return args, nil, err
	}

	if hasFeed != nil {
		args.feed.ID = hasFeed.ID
		return args, createUserFeedState, nil
	}

	return args, createFeedState, nil
}

func createFeedState(ctx context.Context, args FeedArgs) (FeedArgs, rf.StateFn[FeedArgs], error) {
	if err := args.store.CreateFeed(ctx, args.feed); err != nil {
		return args, nil, err
	}

	return args, createUserFeedState, nil
}

func createUserFeedState(ctx context.Context, args FeedArgs) (FeedArgs, rf.StateFn[FeedArgs], error) {
	if err := args.store.CreateUserFeed(ctx, args.feed); err != nil {
		return args, nil, err
	}

	return args, nil, nil
}
