package types

import "strings"

func IsUniqueErr(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
