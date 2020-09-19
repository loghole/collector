package codes

import (
	"net/http"
)

const (
	system   = 1000
	internal = 2000
)

const (
	DatabaseError = system + iota
	SystemError
)

const (
	UnmarshalError = internal + iota
)

func ToHTTP(code int) int {
	switch code {
	case DatabaseError, SystemError:
		return http.StatusInternalServerError
	case UnmarshalError:
		return http.StatusBadRequest
	default:
		return http.StatusTeapot
	}
}
