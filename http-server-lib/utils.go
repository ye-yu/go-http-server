package httpserverlib

import "strings"

func concat(a string, b string) string {
	return strings.Join([]string{a, b}, "")
}
