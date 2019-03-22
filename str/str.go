package str

import (
	"strings"
)

func Empty(s string) bool {
	return 0 == strings.Compare(s, "")
}

func Valid(s string) bool {
	return 0 != strings.Compare(s, "")
}

func JoinBy(sep string, a ...string) string {
	return strings.Join(a, sep)
}
