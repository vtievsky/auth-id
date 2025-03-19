package authidjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenOpts struct {
	SessionID string
	ExpiredAt time.Time
}

func NewAccessToken(signingKey []byte, opts *TokenOpts) ([]byte, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{ //nolint:exhaustruct
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: opts.ExpiredAt},
		},
		Session:    opts.SessionID,
		AccessOnly: true,
	})

	s, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return []byte(s), nil
}

func NewRefreshToken(signingKey []byte, opts *TokenOpts) ([]byte, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{ //nolint:exhaustruct
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: opts.ExpiredAt},
		},
		Session:    opts.SessionID,
		AccessOnly: false,
	})

	s, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return []byte(s), nil
}
