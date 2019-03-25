package _test

import (
	"github.com/Andrew-M-C/go-tools/jsonconv"
	"github.com/Andrew-M-C/go-tools/log"
)

var strStandard = `{
	"a-string": "这是一个string",
	"an-int": 12345678,
	"a-float": 12345.12345678,
	"a-true": true,
	"a-false": false,
	"a-null": null,
	"an-object": {
		"sub-string": "string in an object",
		"sub-object": {
			"another-sub-string": "\"string\" in an object in an object",
			"another-sub-array": [1, "string in sub
array", true, null],
			"complex":"\u4e2d\t\u6587",
			"illegel":"illegal\yillegal"
		}
	},
	"an-array": [
		{"sub-string": "string in an object in an array"},
		56789,
		false,
		null
	]
}`

func testKeyInObject(obj *jsonconv.JsonValue, key interface{}, keys... interface{}) {
	child, err := obj.Get(key, keys...)
	if err != nil {
		log.Error("Failed to get child: %s", err.Error())
	} else {
		if child.IsString() {
			log.Info("Get child, type: %s, value '%s'", child.TypeString(), child.String())
		} else {
			log.Info("Get child, type: %s", child.TypeString())
		}
	}
}

func TestJsonValue() {
	obj, err := jsonconv.NewFromString(strStandard)
	if err != nil {
		log.Error("Failed to parse json: %s", err.Error())
	} else {
		testKeyInObject(obj, "an-object", "sub-object", "complex")
		testKeyInObject(obj, "an-object", "sub-object", "illegal")
		testKeyInObject(obj, "an-array", 0, "sub-string")
		testKeyInObject(obj, "an-object", "sub-object", "another-sub-array", 1)
	}

	json_str := ""
	json_str, _ = obj.Marshal()
	log.Info("re-package json: %s", json_str)
	json_str, _ = obj.Marshal(jsonconv.Option{EnsureAscii: true})
	log.Info("re-package json: %s", json_str)
	json_str, _ = obj.Marshal(jsonconv.Option{FloatDigits: 2})
	log.Info("re-package json: %s", json_str)
	return
}
