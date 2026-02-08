package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Jwt struct {
	key []byte
}

func NewJwt(key string) *Jwt {
	return &Jwt{key: []byte(key)}
}

func (j *Jwt) GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		Subject:   userID.String(),
		Issuer:    "https://sweetops.io",
		Audience:  jwt.ClaimStrings{"https://sweetops.io"},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.key)
}

func (j *Jwt) ParseToken(token string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return j.key, nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
