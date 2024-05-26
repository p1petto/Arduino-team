package storage

import "errors"

var (
	ErrRoomExists   = errors.New("room already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)
