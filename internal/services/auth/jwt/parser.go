package jwt

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

func ParseToken(signingKey string, signedString []byte) (*Token, error) {
	token, err := jwt.Parse(string(signedString), func(atoken *jwt.Token) (interface{}, error) {
		if _, ok := atoken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenParse
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrTokenClaimsInvalid
	}

	session, ok := claims[HEADER_SESSION].(string)
	if !ok {
		return nil, fmt.Errorf("%w | %s", ErrTokenClaimsInvalid, HEADER_SESSION)
	}

	issuedAt := claims[HEADER_ISSUED_AT].(float64)
	if !ok {
		return nil, fmt.Errorf("%w | %s", ErrTokenClaimsInvalid, HEADER_ISSUED_AT)
	}

	expiredAt := claims[HEADER_EXPIRED_AT].(float64)
	if !ok {
		return nil, fmt.Errorf("%w | %s", ErrTokenClaimsInvalid, HEADER_EXPIRED_AT)
	}

	accessOnly := claims[HEADER_ACCESS_ONLY].(bool)
	if !ok {
		return nil, fmt.Errorf("%w | %s", ErrTokenClaimsInvalid, HEADER_ACCESS_ONLY)
	}

	return &Token{
		SessionID:  session,
		IssuedAt:   time.Unix(int64(issuedAt), 0),
		ExpiredAt:  time.Unix(int64(expiredAt), 0),
		AccessOnly: accessOnly,
		Valid:      token.Valid,
	}, nil
}
