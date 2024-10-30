package utils

import "errors"

var (
	ErrDirNotExist  = errors.New("no such directory")
	ErrFileNotExist = errors.New("no such file")
)
