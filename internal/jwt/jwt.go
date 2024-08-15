package jwt

import (
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/dwaynedwards/rss-feed-aggregator-in-go/internal/errors"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAndSignUserID(userID int64, ttl time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"ttl":    ttl.Unix(),
	})

	tokenString, err := token.SignedString([]byte(rf.Config.JWTSecret))
	if err != nil {
		return "", errors.InternalErrorf("%s: %v", errors.ErrTokenGenerationFailed, err)
	}

	return tokenString, nil
}

func ParseAndVerifyUserID(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.InternalErrorf("%s: %v", errors.ErrTokenUnexpactedSigningMethod, token.Header["alg"])
		}

		return []byte(rf.Config.JWTSecret), nil
	})
	if err != nil {
		return 0, errors.Unauthorizedf("%s: %v", errors.ErrTokenParseFailed, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.Unauthorizedf(errors.ErrTokenClaimsFailed)
	}

	ttl := int64(claims["ttl"].(float64))
	userID := int64(claims["userID"].(float64))

	if ttl < time.Now().Unix() {
		return 0, errors.Unauthorizedf(errors.ErrTokenExpired)
	}

	return userID, nil
}
