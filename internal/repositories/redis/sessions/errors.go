package redissessions

import "errors"

var (
	ErrSessionPrivilegeNotFound = errors.New("session privilege not found")
)
