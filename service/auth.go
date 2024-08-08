package service

import (
	"context"
	"time"

	"github.com/alexedwards/argon2id"
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

func (a *AuthService) SignUp(ctx context.Context, auth *rf.Auth) (string, error) {
	if err := validateSignUp(auth); err != nil {
		return "", err
	}

	hashedPassword, err := argon2id.CreateHash(auth.BasicAuth.Password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	auth.BasicAuth.Password = hashedPassword

	err = a.store.Create(ctx, auth)
	if err != nil {
		return "", err
	}

	token, err := rf.GenerateAndSignJWT(auth.UserID, time.Now().Add(time.Hour*24))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) SignIn(ctx context.Context, auth *rf.Auth) (string, error) {
	if err := validateSignIn(auth); err != nil {
		return "", err
	}

	authFound, err := a.store.FindByEmail(ctx, auth.BasicAuth.Email)
	if err != nil {
		return "", err
	}

	if authFound == nil {
		return "", rf.AppErrorf(rf.ECUnautherized, rf.EMInvlidCredentials)
	}

	match, err := argon2id.ComparePasswordAndHash(auth.BasicAuth.Password, authFound.BasicAuth.Password)
	if err != nil {
		return "", err
	}
	if !match {
		return "", rf.AppErrorf(rf.ECUnautherized, rf.EMInvlidCredentials)
	}

	auth.UserID = authFound.UserID

	token, err := rf.GenerateAndSignJWT(authFound.UserID, time.Now().Add(time.Hour*24))
	if err != nil {
		return "", err
	}

	return token, nil
}

func validateSignUp(auth *rf.Auth) error {
	errs := map[string]string{}

	if auth.BasicAuth.Email == "" {
		errs["email"] = rf.EMEmailRequired
	}

	if auth.BasicAuth.Password == "" {
		errs["password"] = rf.EMPasswordRequired
	}

	if auth.User == nil || auth.User.Name == "" {
		errs["name"] = rf.EMNameRequired
	}

	if len(errs) > 0 {
		return rf.BadRequestAppError(errs)
	}

	return nil
}

func validateSignIn(auth *rf.Auth) error {
	errs := map[string]string{}

	if auth.BasicAuth.Email == "" {
		errs["email"] = rf.EMEmailRequired
	}

	if auth.BasicAuth.Password == "" {
		errs["password"] = rf.EMPasswordRequired
	}

	if len(errs) > 0 {
		return rf.BadRequestAppError(errs)
	}

	return nil
}
