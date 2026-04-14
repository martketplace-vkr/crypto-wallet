package wallet

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("deposit address not found")
	ErrUnsupported     = errors.New("unsupported network or asset")
)
