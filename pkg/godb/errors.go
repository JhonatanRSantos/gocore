package godb

import "errors"

var (
	ErrUnknown                   = errors.New("unknown error")
	ErrInvalidDBType             = errors.New("invalid database type")
	ErrConnectionTimeoutExceeded = errors.New("connection timeout exceeded")
)
