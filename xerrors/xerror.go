package xerrors

import (
	"errors"
)

var (
	ErrPasswordTooShort = errors.New("password length too short, minimum of 8 characters required")
)
