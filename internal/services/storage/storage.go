package storage

import "errors"

var (
	ErrorUserNotFound = errors.New("user not found")
	ErrorUserExists   = errors.New("user already exists")
	ErrorAppNotFound  = errors.New("app not found")
)
