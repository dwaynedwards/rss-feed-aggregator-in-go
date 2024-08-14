package service

import (
	"context"
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
)

type AuthArgs struct {
	store          rf.AuthStore
	auth           *rf.Auth
	authToValidate *rf.Auth
}

func (as AuthArgs) validateSignUp() error {
	if as.store == nil {
		return rf.NewAppError(rf.ECIntenal, "store cannot be nil")
	}

	errs := map[string]string{}

	if as.auth == nil || as.auth.BasicAuth == nil || as.auth.BasicAuth.Email == "" {
		errs["email"] = rf.EMEmailRequired
	}

	if as.auth == nil || as.auth.BasicAuth == nil || as.auth.BasicAuth.Password == "" {
		errs["password"] = rf.EMPasswordRequired
	}

	if as.auth == nil || as.auth.User == nil || as.auth.User.Name == "" {
		errs["name"] = rf.EMNameRequired
	}

	if len(errs) > 0 {
		return rf.NewAppError(rf.ECInvalid, errs)
	}

	return nil
}

func (as AuthArgs) validateSignIn() error {
	if as.store == nil {
		return rf.NewAppError(rf.ECIntenal, "store cannot be nil")
	}

	if as.auth == nil {
		return rf.NewAppError(rf.ECIntenal, "auth cannot be nil")
	}

	errs := map[string]string{}

	if as.auth.BasicAuth == nil || as.auth.BasicAuth.Email == "" {
		errs["email"] = rf.EMEmailRequired
	}

	if as.auth.BasicAuth == nil || as.auth.BasicAuth.Password == "" {
		errs["password"] = rf.EMPasswordRequired
	}

	if len(errs) > 0 {
		return rf.NewAppError(rf.ECInvalid, errs)
	}

	return nil
}

func canSignUpCheckState(ctx context.Context, args AuthArgs) (AuthArgs, rf.StateFn[AuthArgs], error) {
	hasAuth, err := args.store.FindByEmail(ctx, args.auth.BasicAuth.Email)
	if err != nil {
		return args, nil, err
	}

	if hasAuth != nil {
		return args, nil, rf.NewAppError(rf.ECInvalid, rf.EMUserExists)
	}

	return args, createAuthAndUserState, nil
}

func canSignInCheckState(ctx context.Context, args AuthArgs) (AuthArgs, rf.StateFn[AuthArgs], error) {
	hasAuth, err := args.store.FindByEmail(ctx, args.auth.BasicAuth.Email)
	if err != nil {
		return args, nil, err
	}

	if hasAuth == nil {
		return args, nil, rf.NewAppError(rf.ECUnautherized, rf.EMInvlidCredentials)
	}

	args.authToValidate = hasAuth
	return args, validateAuthState, nil
}

func createAuthAndUserState(ctx context.Context, args AuthArgs) (AuthArgs, rf.StateFn[AuthArgs], error) {
	hashedPassword, err := rf.Hash(args.auth.BasicAuth.Password)
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

func validateAuthState(ctx context.Context, args AuthArgs) (AuthArgs, rf.StateFn[AuthArgs], error) {
	match, err := rf.Matches(args.auth.BasicAuth.Password, args.authToValidate.BasicAuth.Password)
	if err != nil {
		return args, nil, err
	}
	if !match {
		return args, nil, rf.NewAppError(rf.ECUnautherized, rf.EMInvlidCredentials)
	}

	args.auth.UserID = args.authToValidate.UserID
	return args, generateUserTokenState, nil
}

func generateUserTokenState(ctx context.Context, args AuthArgs) (AuthArgs, rf.StateFn[AuthArgs], error) {
	token, err := rf.GenerateAndSignUserIDJWT(args.auth.UserID, time.Now().Add(time.Hour*24))
	if err != nil {
		return args, nil, err
	}

	args.auth.Token = token
	return args, nil, nil
}
