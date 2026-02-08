package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
