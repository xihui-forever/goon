package session

import "errors"

var (
	ErrSessionGenerateFail = errors.New("session generation failed")
	ErrSessionNotExist     = errors.New("session not exists")
)
