package dberrors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserScan          = errors.New("user scan error")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleScan          = errors.New("role scan error")
	ErrRoleAlreadyExists = errors.New("role already exists")
)
