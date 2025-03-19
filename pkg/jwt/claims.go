package authidjwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	Session    string `json:"session_id"`
	AccessOnly bool   `json:"access_only"`
}
