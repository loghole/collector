package codes

import (
	"net/http"
)

const (
	systemCodes   = 1000
	usecasesCodes = 2000
)

const (
	DatabaseError = systemCodes + iota
	SystemError
)

const (
	UnmarshalError = usecasesCodes + iota
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
