package rf

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAndSignJWT(userID int64, ttl time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"ttl":    ttl.Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", NewAppError(ECIntenal, fmt.Sprintf("JWT generation failed: %v", err))
	}

	return tokenString, nil
}

func ParseAndVerifyJWT(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, NewAppError(ECIntenal, fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, NewAppError(ECUnautherized, fmt.Sprintf("JWT parse failed: %v", err))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, NewAppError(ECUnautherized, fmt.Sprintf("JWT claims failed: %v", err))
	}

	ttl := int64(claims["ttl"].(float64))
	userID := int64(claims["userID"].(float64))

	if ttl < time.Now().Unix() {
		return 0, NewAppError(ECUnautherized, fmt.Sprintf("JWT expired: %v", err))
	}

	return userID, nil
}
