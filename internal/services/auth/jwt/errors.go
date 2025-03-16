package jwt

import "errors"

var (
	ErrTokenParse         = errors.New("there was an error in parsing")
	ErrTokenClaimsInvalid = errors.New("token claims invalid")
)
