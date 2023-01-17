package session

import "errors"

var (
	ErrSessionNotExist     = errors.New("session not exist")
	ErrSessionExpired      = errors.New("session expired")
	ErrSessionGenerateFail = errors.New("session generate failed")
)
