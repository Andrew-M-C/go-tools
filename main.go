package main

import (
	test "github.com/Andrew-M-C/go-tools/_test"
	"github.com/Andrew-M-C/go-tools/log"
)

func main() {
	// log.SetLogLevel(log.LEVEL_NONE)
	// test.TestSqlToJson()
	// test.TestJsonValue()
	// test.TestJsonObjSort()
	// test.TestJsonArraySort()
	// test.TestStr()

	// {
	// 	var t1, t2, t3, t4, t5 int64
	// 	COUNT := int64(10000)
	// 	for i := int64(0); i < COUNT; i++ {
	// 		e1, e2, e3, e4, e5 := test.TestOrigJsonEffenciency()
	// 		t1 += e1
	// 		t2 += e2
	// 		t3 += e3
	// 		t4 += e4
	// 		t5 += e5
	// 	}
	// 	t1 /= COUNT
	// 	t2 /= COUNT
	// 	t3 /= COUNT
	// 	t4 /= COUNT
	// 	t5 /= COUNT
		log.SetLogLevel(log.LEVEL_DEBUG)
	// 	log.Info("unmarshal to map: %v", t1)
	// 	log.Info("unmarshal to map and parse: %v", t2)
	// 	log.Info("unmarshal to struct: %v", t3)
	// 	log.Info("jsonparser: %v", t4)
	// 	log.Info("jsonconv: %v", t5)
	// }

	// test.TestAwsomeEscapingJson()
	// test.TestJsonMerge()

	// test.TestReadSqlKVs()

	test.TestXmlconv()

	log.Info("demo done")
	return
}
