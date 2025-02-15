package xerrors

import (
	"errors"
)

var (
	ErrPasswordTooShort   = errors.New("password length too short, minimum of 8 characters required")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoSuchUser         = errors.New("no such user")
	ErrEmptyTitle         = errors.New("post title cannot be empty")
	ErrEmptyContent       = errors.New("post content cannot be empty")
	ErrInvalidUser        = errors.New("invalid user ID")
	ErrInvalidPost        = errors.New("post not found")
	ErrNoCategory         = errors.New("invalid category")
	ErrInvalidVoteType    = errors.New("invalid vote type")
	ErrInvalidRequest     = errors.New("invalid vote request")
)
