package reposessions

import "errors"

var (
	ErrSessionPrivilegesEmpty   = errors.New("session privileges empty")
	ErrSessionPrivilegeNotFound = errors.New("session privilege not found")
	ErrSessionCartNotFound      = errors.New("session cart not found")
)
