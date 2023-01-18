package admin

import "errors"

var (
	ErrAdminExist           = errors.New("admin exist")
	ErrAdminNotExist        = errors.New("admin not exist")
	ErrPasswordWrong        = errors.New("password wrong")
	ErrPasswordChangeFailed = errors.New("password change failed")
)
