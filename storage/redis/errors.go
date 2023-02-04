package redis

import "errors"

var (
	ErrRedisSetExpireFailure = errors.New("timeOut set fail")
	ErrSidNotExist           = errors.New("sid not exists")
)
