package httptransport

import "errors"

var (
	ErrDeleteHimself = errors.New("unacceptable to delete yourself")
	ErrBlockHimself  = errors.New("unacceptable to block yourself")
)
