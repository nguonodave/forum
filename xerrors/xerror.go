package xerrors

import (
	"errors"
)

var (
	ErrPasswordTooShort = errors.New("password length too short, minimum of 8 characters required")
	ErrInvalidUser      = errors.New("invalid user ID")
	ErrInvalidVoteType  = errors.New("invalid vote type")
	ErrInvalidRequest   = errors.New("invalid vote request")
	ErrWrongEmailFormat = errors.New("wrong email format")
)
