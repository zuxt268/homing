package entity

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrInternalServerError = errors.New("internal server error")
	ErrNotImplemented      = errors.New("not implemented")
)

type HomingErr struct {
	Err     error
	Message string
}

func (h *HomingErr) Unwrap() error { return h.Err }

func (h *HomingErr) Error() string {
	if h.Err == nil {
		return fmt.Sprintf("msg: %s", h.Message)
	}
	return fmt.Sprintf("err: %s, msg: %s", h.Err.Error(), h.Message)
}
