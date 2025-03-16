package authtoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	SessionID  string
	IssuedAt   time.Time
	ExpiredAt  time.Time
	AccessOnly bool
	Valid      bool
}

func ParseToken(signingKey []byte, signedString []byte) (*Token, error) {
	var (
		claims  Claims
		keyFunc = func(t *jwt.Token) (any, error) {
			return signingKey, nil
		}
	)

	token, err := jwt.ParseWithClaims(string(signedString), &claims, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("%w | %v", ErrTokenParse, err)
	}

	return &Token{
		SessionID:  claims.Session,
		IssuedAt:   claims.IssuedAt.Time,
		ExpiredAt:  claims.ExpiresAt.Time,
		AccessOnly: claims.AccessOnly,
		Valid:      token.Valid,
	}, nil
}
