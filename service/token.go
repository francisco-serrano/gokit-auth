package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const key = "abc123"

type customClaims struct {
	jwt.StandardClaims
	SessionID string
}

func CreateToken(sessionID string) (string, error) {
	claims := &customClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		SessionID: sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error while signing JWT: %w", err)
	}

	return signedToken, nil
}

func ParseToken(token string) (string, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(key), nil
	})
	if err != nil {
		return "", fmt.Errorf("error while parsing JWT: %s", err)
	}

	if !parsedToken.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := parsedToken.Claims.(*customClaims)
	if !ok {
		return "", fmt.Errorf("could not obtain token claims: %T", parsedToken.Claims)
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return "", fmt.Errorf("token expired")
	}

	return claims.SessionID, nil
}
