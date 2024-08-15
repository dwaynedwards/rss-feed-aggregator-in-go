package password

import (
	"github.com/alexedwards/argon2id"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
)

func Hash(plaintextPassword string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(plaintextPassword, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func Matches(plaintextPassword, hashedPassword string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plaintextPassword, hashedPassword)
	if err != nil {
		return false, err
	}
	if !match {
		return false, errors.Unauthorizedf(errors.ErrInvalidCredentials)
	}
	return true, nil
}
