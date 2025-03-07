package dberrors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserScan          = errors.New("user scan error")
	ErrUserAlreadyExists = errors.New("user already exists")
)
