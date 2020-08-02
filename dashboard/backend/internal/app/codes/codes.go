package codes

import (
	"net/http"
	"strconv"
)

const (
	systemCodes     = 1000
	usecasesCodes   = 2000
	validationCodes = 200_000
)

const (
	DatabaseError = systemCodes + iota
	InvalidJSONError
)

const (
	InvalidSuggestType = usecasesCodes + iota
)

func ToHTTP(code int) int {
	switch code {
	case DatabaseError:
		return http.StatusInternalServerError
	case InvalidJSONError, InvalidSuggestType:
		return http.StatusBadRequest
	default:
		return http.StatusTeapot
	}
}

// Validation error codes
type Code int

func (s Code) String() string {
	return strconv.Itoa(int(s))
}

const (
	ValidMinLimit Code = validationCodes + iota
	ValidMaxLimit
	ValidMinOffset
	ValidQueryParamsTypeRequired
	ValidQueryParamsTypeIn
	ValidQueryParamsKeyRequired
	ValidQueryParamsValueRequired
	ValidQueryParamsValueItemEmpty
	ValidQueryParamsValueListEmpty
	ValidQueryParamsOperatorRequired
	ValidQueryParamsOperatorIn
	ValidSuggestTypeIn
	ValidSuggestValueLength
)
