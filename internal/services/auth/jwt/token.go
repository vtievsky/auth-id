package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	HEADER_SESSION     = "session_id"
	HEADER_ISSUED_AT   = "iat"
	HEADER_EXPIRED_AT  = "exp"
	HEADER_ACCESS_ONLY = "access_only"
)

type TokenOpts struct {
	SessionID string
	ExpiredAt time.Time
}

func NewAccessToken(signingKey string, opts *TokenOpts) ([]byte, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		HEADER_SESSION:     opts.SessionID,
		HEADER_ISSUED_AT:   time.Now().Unix(),
		HEADER_EXPIRED_AT:  opts.ExpiredAt.Unix(),
		HEADER_ACCESS_ONLY: true,
	})

	s, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return nil, err
	}

	return []byte(s), nil
}

func NewRefreshToken(signingKey string, opts *TokenOpts) ([]byte, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		HEADER_SESSION:     opts.SessionID,
		HEADER_ISSUED_AT:   time.Now().Unix(),
		HEADER_EXPIRED_AT:  opts.ExpiredAt.Unix(),
		HEADER_ACCESS_ONLY: false,
	})

	s, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return nil, err
	}

	return []byte(s), nil
}
