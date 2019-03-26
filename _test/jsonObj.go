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
			"complex":"\u4e2d\t\u6587"
		}
	},
	"an-array": [
		{
			"sub-string": "string in an object in an array",
			"sub-sub-array": [
				{
					"sub-sub-string": "string in an object in an array in an object in an string"
				}
			]
		},
		56789,
		false,
		null
	]
}`

func testKeyInObject(obj *jsonconv.JsonValue, key interface{}, keys... interface{}) {
	log.Info("Test %v %v", key, keys)
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
	log.Info("======== Start testing JsonValue")
	obj, err := jsonconv.NewFromString(strStandard)
	if err != nil {
		log.Error("Failed to parse json: %s", err.Error())
	} else {
		testKeyInObject(obj, "an-object", "sub-object", "complex")
		// testKeyInObject(obj, "an-object", "sub-object", "illegal")
		testKeyInObject(obj, "an-array", 0, "sub-string")		// ERROR OCCURRED !!!
		testKeyInObject(obj, "an-object", "sub-object", "another-sub-array", 1)
	}

	json_str := ""
	json_str, _ = obj.Marshal()
	log.Info("re-package json: %s", json_str)
	json_str, _ = obj.Marshal(jsonconv.Option{EnsureAscii: true})
	log.Info("re-package json: %s", json_str)
	json_str, _ = obj.Marshal(jsonconv.Option{FloatDigits: 2})
	log.Info("re-package json: %s", json_str)

	// test foreach
	log.Info("Now test object foreach")
	obj.ObjectForeach(func (key string, value *jsonconv.JsonValue) error {
		log.Info("Key - %s, type %s", key, value.TypeString())
		return nil
	})

	log.Info("Now test array foreach")
	arr, err := obj.Get("an-array")
	log.Info("array type: %s", arr.TypeString())
	if err != nil {
		log.Error("error: %s", err.Error())
	} else {
		arr.ArrayForeach(func (index int, value *jsonconv.JsonValue) error {
			log.Info("[%d] type %s", index, value.TypeString())
			return nil
		})
	}

	// test modification
	err = obj.Set(jsonconv.NewString("THIS IS A FULL NEW STRING"), "an-array", 0, "sub-string")
	if err != nil {
		log.Error("Failed to set: %s", err.Error())
	}
	err = obj.Set(jsonconv.NewString("THIS IS ANOTHER FULL NEW STRING"), "an-array", 0, "sub-sub-array", 0, "sub-sub-string")
	if err != nil {
		log.Error("Failed to set: %s", err.Error())
	}
	json_str, _ = obj.Marshal()
	log.Info("new json after modification: %s", json_str)

	// return
	return
}
