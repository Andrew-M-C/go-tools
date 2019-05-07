package main

import (
	test "github.com/Andrew-M-C/go-tools/_test"
	"github.com/Andrew-M-C/go-tools/log"
)

func main() {
	log.SetLogLevel(log.LEVEL_DEBUG)
	// test.TestSqlToJson()
	// test.TestJsonValue()
	// test.TestStr()
	// for i := 0; i < 200; i ++ {
	// 	log.Info("TestOrigJsonEffenciency")
	// 	test.TestOrigJsonEffenciency()
	// }
	test.TestAwsomeEscapingJson()
	log.Info("demo done")
	return
}
