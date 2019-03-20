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
			"another-sub-string": "string in an object in an object",
			"another-sub-array": [1, "array", true, null]
		}
	},
	"an-array": [
		{"sub-string": "string in an object in an array"},
		56789,
		false,
		null
	]
}`

// TODO: "string in an object in an object" 后面的逗号如果没有了，会有问题，但是代码里面并没有检查出来

func TestJsonObj() {
	obj, err := jsonconv.NewFromString(strStandard)
	if err != nil {
		log.Error("Failed to parse json: %s", err.Error())
	} else {
		child, err := obj.GetByKey("an-object", "sub-object", "another-sub-string")
		if err != nil {
			log.Error("Failed to get child: %s", err.Error())
		} else {
			log.Info("Get child, type: %s, value '%s'", child.TypeString(), child.String())
		}
	}
	return
}
