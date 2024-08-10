package service

import (
	"context"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
)

type AuthService struct {
	store rf.AuthStore
}

func NewAuthService(store rf.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (as *AuthService) SignUp(ctx context.Context, auth *rf.Auth) (string, error) {
	args := AuthArgs{
		store: as.store,
		auth:  auth,
	}

	if err := args.validateSignUp(); err != nil {
		return "", err
	}

	err := rf.RunStateMachine(ctx, args, canSignUpCheckState)
	if err != nil {
		return "", err
	}

	return auth.Token, nil
}

func (as *AuthService) SignIn(ctx context.Context, auth *rf.Auth) (string, error) {
	args := AuthArgs{
		store: as.store,
		auth:  auth,
	}

	if err := args.validateSignIn(); err != nil {
		return "", err
	}

	err := rf.RunStateMachine(ctx, args, canSignInCheckState)
	if err != nil {
		return "", err
	}

	return auth.Token, nil
}
