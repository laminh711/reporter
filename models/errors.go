package models

import "errors"

var (
	// ErrInternalServerError will throw if any internal server error happen
	ErrInternalServerError = errors.New("Internal Server Error")
)
