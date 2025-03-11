package usersvc

import "errors"

var (
	ErrInvalidName      = errors.New("invalid name error")
	ErrInvalidLogin     = errors.New("invalid login error")
	ErrInvalidPassword  = errors.New("invalid password error")
	ErrGeneratePassword = errors.New("generate password error")
)
