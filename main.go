package main

import (
	test "github.com/Andrew-M-C/go-tools/_test"
	"github.com/Andrew-M-C/go-tools/log"
)

func main() {
	test.TestSqlToJson()
	test.TestJsonValue()
	test.TestStr()
	log.Info("demo done")
	return
}
