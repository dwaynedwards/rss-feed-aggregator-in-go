package authservice

import (
	"context"
	"log"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/builder"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/statemachine"
)

type AuthStore interface {
	CreateAuthAndUser(ctx context.Context, auth *rf.Auth) error
	FindByEmail(ctx context.Context, email string) (*rf.Auth, error)
}

type AuthService struct {
	store AuthStore
}

func NewAuthService(store AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (as *AuthService) SignUp(ctx context.Context, req *rf.SignUpRequest) (string, error) {
	args := AuthArgs{
		store: as.store,
		auth: builder.NewAuthBuilder().
			WithUser(builder.NewUserBuilder().
				WithName(req.Name)).
			WithBasicAuth(builder.NewBasicAuthBuilder().
				WithEmail(req.Email).
				WithPassword(req.Password)).
			Build(),
	}

	if err := args.validateSignUp(); err != nil {
		log.Printf("Val: %#v", err)
		return "", err
	}

	result, err := statemachine.Run(ctx, args, canSignUpCheckState)
	if err != nil {
		return "", err
	}

	return result.auth.Token, nil
}

func (as *AuthService) SignIn(ctx context.Context, req *rf.SignInRequest) (string, error) {
	args := AuthArgs{
		store: as.store,
		auth: builder.NewAuthBuilder().
			WithBasicAuth(builder.NewBasicAuthBuilder().
				WithEmail(req.Email).
				WithPassword(req.Password)).
			Build(),
	}

	if err := args.validateSignIn(); err != nil {
		return "", err
	}

	result, err := statemachine.Run(ctx, args, canSignInCheckState)
	if err != nil {
		return "", err
	}

	return result.auth.Token, nil
}
