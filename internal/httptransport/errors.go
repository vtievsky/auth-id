package httptransport

import "errors"

var (
	ErrDeleteHimself      = errors.New("unacceptable to delete yourself")
	ErrDeleteHimselfRoles = errors.New("unacceptable to delete your roles")
	ErrAddHimselfRoles    = errors.New("unacceptable to add your roles")
	ErrBlockHimself       = errors.New("unacceptable to block yourself")
)
