package _test

import (
	"github.com/Andrew-M-C/go-tools/str"
	"github.com/Andrew-M-C/go-tools/log"
)

func TestStr() {
	s := str.JoinBy(",", "a", "BB", "ccc", "4444")
	log.Debug("Join result: %s", s)
}
