package tools

import (
	"strings"
)

func CreateLike(value string) string {
	return strings.Join([]string{"%", value, "%"}, "")
}
