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
