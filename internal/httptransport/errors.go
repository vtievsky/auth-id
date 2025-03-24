package httptransport

import "errors"

var (
	ErrDeleteHimself      = errors.New("unacceptable to delete yourself")
	ErrDeleteHimselfRoles = errors.New("unacceptable to delete your roles")
	ErrAddHimselfRoles    = errors.New("unacceptable to add your roles")
	ErrHimself            = errors.New("unacceptable to processing yourself")
	ErrBlockHimself       = errors.New("unacceptable to block yourself")
)
