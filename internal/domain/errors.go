package domain

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
	ErrValidation = errors.New("bad request")
	ErrInternal   = errors.New("internal error")
)
