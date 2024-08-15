package feedservice

import (
	"context"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/statemachine"
)

type FeedArgs struct {
	store FeedStore
	feed  *rf.Feed
}

func (fs FeedArgs) validateAddFeed() error {
	if fs.store == nil {
		return errors.InternalErrorf("store cannot be nil")
	}

	if fs.feed == nil || fs.feed.URL == "" {
		return errors.InvalidDataf(errors.ErrURLRequired)
	}

	return nil
}

func addFeedState(ctx context.Context, args FeedArgs) (FeedArgs, statemachine.StateFn[FeedArgs], error) {
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

func createFeedState(ctx context.Context, args FeedArgs) (FeedArgs, statemachine.StateFn[FeedArgs], error) {
	if err := args.store.CreateFeed(ctx, args.feed); err != nil {
		return args, nil, err
	}

	return args, createUserFeedState, nil
}

func createUserFeedState(ctx context.Context, args FeedArgs) (FeedArgs, statemachine.StateFn[FeedArgs], error) {
	if err := args.store.CreateUserFeed(ctx, args.feed); err != nil {
		return args, nil, err
	}

	return args, nil, nil
}
