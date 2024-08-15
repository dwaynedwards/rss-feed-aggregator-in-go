package authservice

import (
	"context"
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/jwt"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/password"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/statemachine"
)

type AuthArgs struct {
	store          AuthStore
	auth           *rf.Auth
	authToValidate *rf.Auth
}

func (as AuthArgs) validateSignUp() error {
	if as.store == nil {
		return errors.InternalErrorf("store cannot be nil")
	}

	errs := map[string]string{}

	if as.auth == nil || as.auth.BasicAuth == nil || as.auth.BasicAuth.Email == "" {
		errs["email"] = errors.ErrEmailRequired
	}

	if as.auth == nil || as.auth.BasicAuth == nil || as.auth.BasicAuth.Password == "" {
		errs["password"] = errors.ErrPasswordRequired
	}

	if as.auth == nil || as.auth.User == nil || as.auth.User.Name == "" {
		errs["name"] = errors.ErrNameRequired
	}

	if len(errs) > 0 {
		return errors.InvalidError(errs)
	}

	return nil
}

func (as AuthArgs) validateSignIn() error {
	if as.store == nil {
		return errors.InternalErrorf("store cannot be nil")
	}

	errs := map[string]string{}

	if as.auth == nil || as.auth.BasicAuth == nil || as.auth.BasicAuth.Email == "" {
		errs["email"] = errors.ErrEmailRequired
	}

	if as.auth == nil || as.auth.BasicAuth == nil || as.auth.BasicAuth.Password == "" {
		errs["password"] = errors.ErrPasswordRequired
	}

	if len(errs) > 0 {
		return errors.InvalidError(errs)
	}

	return nil
}

func canSignUpCheckState(ctx context.Context, args AuthArgs) (AuthArgs, statemachine.StateFn[AuthArgs], error) {
	hasAuth, err := args.store.FindByEmail(ctx, args.auth.BasicAuth.Email)
	if err != nil {
		return args, nil, err
	}

	if hasAuth != nil {
		return args, nil, errors.InvalidDataf(errors.ErrCouldNotProcess)
	}

	return args, createAuthAndUserState, nil
}

func canSignInCheckState(ctx context.Context, args AuthArgs) (AuthArgs, statemachine.StateFn[AuthArgs], error) {
	hasAuth, err := args.store.FindByEmail(ctx, args.auth.BasicAuth.Email)
	if err != nil {
		return args, nil, err
	}

	if hasAuth == nil {
		return args, nil, errors.Unauthorizedf(errors.ErrInvalidCredentials)
	}

	args.authToValidate = hasAuth
	return args, validateAuthState, nil
}

func createAuthAndUserState(ctx context.Context, args AuthArgs) (AuthArgs, statemachine.StateFn[AuthArgs], error) {
	hashedPassword, err := password.Hash(args.auth.BasicAuth.Password)
	if err != nil {
		return args, nil, err
	}

	args.auth.BasicAuth.Password = hashedPassword

	err = args.store.CreateAuthAndUser(ctx, args.auth)
	if err != nil {
		return args, nil, err
	}

	return args, generateUserTokenState, nil
}

func validateAuthState(ctx context.Context, args AuthArgs) (AuthArgs, statemachine.StateFn[AuthArgs], error) {
	match, err := password.Matches(args.auth.BasicAuth.Password, args.authToValidate.BasicAuth.Password)
	if err != nil {
		return args, nil, err
	}
	if !match {
		return args, nil, errors.Unauthorizedf(errors.ErrInvalidCredentials)
	}

	args.auth.UserID = args.authToValidate.UserID
	return args, generateUserTokenState, nil
}

func generateUserTokenState(ctx context.Context, args AuthArgs) (AuthArgs, statemachine.StateFn[AuthArgs], error) {
	token, err := jwt.GenerateAndSignUserID(args.auth.UserID, time.Now().Add(time.Hour*24))
	if err != nil {
		return args, nil, err
	}

	args.auth.Token = token
	return args, nil, nil
}
