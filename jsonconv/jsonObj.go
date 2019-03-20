package jsonconv

import (
	"github.com/buger/jsonparser"
	"github.com/Andrew-M-C/go-tools/log"
	"strconv"
	"strings"
)

// data definitions same as jsonparser
type ValueType int
const (
	NotExist = jsonparser.NotExist
	String	= jsonparser.String
	Number	= jsonparser.Number
	Object	= jsonparser.Object
	Array	= jsonparser.Array
	Boolean	= jsonparser.Boolean
	Null	= jsonparser.Null
	Unknown	= jsonparser.Unknown
)

type JsonObj struct {
	origBytes		[]byte
	// type
	valueType		jsonparser.ValueType
	// values
	stringValue		string
	intValue		int64
	floatValue		float64
	boolValue		bool
	// object children
	objChildren		map[string]*JsonObj
	// array children
	arrChildren		[]*JsonObj
}

// ====================
// New() functions

func NewFromString(s string) (*JsonObj, error) {
	var obj *JsonObj
	var err error

	// get first character
	for index, chr := range s {
		switch chr {
		case ' ', '\r', '\n', '\t':
			// continue
		case '{':
			obj = NewObjectObject()
			err = obj.parseObject([]byte(s[index:]))
			if err != nil {
				return nil, err
			} else {
				return obj, nil
			}
		case '[':
			obj = NewArrayObject()
			err = obj.parseArray([]byte(s[index:]))
			if err != nil {
				return nil, err
			} else {
				return obj, nil
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			obj = new(JsonObj)
			obj.valueType = Number
			obj.intValue, err = strconv.ParseInt(s[index:], 10, 64)
			if err != nil {
				return nil, err
			}
			obj.floatValue, err = strconv.ParseFloat(s[index:], 64)
			if err != nil {
				return nil, err
			}
			return obj, nil
		case '"':
			obj = new(JsonObj)
			obj.valueType = String
			// search for next quote
			next := strings.IndexRune(s[index + 1:], '"')
			if next < 0 {
				return nil, JsonFormatError
			} else {
				obj.stringValue = s[index + 1 : index + 1 + next]
				return obj, nil
			}
		case 't':
			if s[index:] == "true" {
				obj = new(JsonObj)
				obj.valueType = Boolean
				obj.boolValue = true
				return obj, nil
			} else {
				return nil, JsonFormatError
			}
		case 'f':
			if s[index:] == "false" {
				obj = new(JsonObj)
				obj.valueType = Boolean
				obj.boolValue = false
				return obj, nil
			} else {
				return nil, JsonFormatError
			}
		default:
			log.Debug("Skip: %c", chr)
		}
	}
	return nil, JsonFormatError
}

func NewStringObject(s string) *JsonObj {
	obj := new(JsonObj)
	obj.valueType = String
	obj.stringValue = s
	return obj
}

func NewIntObject(i int64) *JsonObj {
	obj := new(JsonObj)
	obj.valueType = Number
	obj.intValue = i
	obj.floatValue = float64(i)
	return obj
}

func NewFloatObject(f float64) *JsonObj {
	obj := new(JsonObj)
	obj.valueType = Number
	obj.intValue = int64(f)
	obj.floatValue = f
	return obj
}

func NewBoolObject(b bool) *JsonObj {
	obj := new(JsonObj)
	obj.valueType = Boolean
	obj.boolValue = b
	return obj
}

func NewNullObject() *JsonObj {
	obj := new(JsonObj)
	obj.valueType = Null
	return obj
}

func NewObjectObject() *JsonObj {
	obj := new(JsonObj)
	obj.valueType = Object
	obj.objChildren = make(map[string]*JsonObj)
	return obj
}

func NewArrayObject() *JsonObj {
	obj := new(JsonObj)
	obj.valueType = Array
	obj.arrChildren = make([]*JsonObj, 0, 10)
	return obj
}

// ====================
// parse functions

func (obj *JsonObj) parseObject(data []byte) error {
	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, _ int) error {
		log.Debug("key: %s", string(key))
		log.Debug("value: %s", string(value))
		switch dataType {
		case jsonparser.String:
			child := NewStringObject(string(value))
			obj.objChildren[string(key)] = child
		case jsonparser.Number:
			child := new(JsonObj)
			child.valueType = Number
			child.origBytes = value
			obj.objChildren[string(key)] = child
		case jsonparser.Object:
			child := NewObjectObject()
			err := child.parseObject(value)
			if err == nil {
				obj.objChildren[string(key)] = child
			}
		case jsonparser.Array:
			child := NewArrayObject()
			err := child.parseArray(value)
			if err == nil {
				obj.objChildren[string(key)] = child
			}
		case jsonparser.Boolean:
			b, err := strconv.ParseBool(string(value))
			if err == nil {
				child := NewBoolObject(b)
				obj.objChildren[string(key)] = child
			}
		case jsonparser.Null:
			child := NewNullObject()
			obj.objChildren[string(key)] = child
		default:
			log.Debug("Invalid type: %d", int(dataType))
		}
		return nil
	})
	return nil
}

func (obj *JsonObj) parseArray(data []byte) error {
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, _ int, _ error) {
		switch dataType {
		case jsonparser.String:
			child := NewStringObject(string(value))
			obj.arrChildren = append(obj.arrChildren, child)
		case jsonparser.Number:
			child := new(JsonObj)
			child.valueType = Number
			child.origBytes = value
			obj.arrChildren = append(obj.arrChildren, child)
		case jsonparser.Object:
			child := NewObjectObject()
			err := child.parseObject(value)
			if err == nil {
				obj.arrChildren = append(obj.arrChildren, child)
			}
		case jsonparser.Array:
			child := NewArrayObject()
			err := child.parseArray(value)
			if err == nil {
				obj.arrChildren = append(obj.arrChildren, child)
			}
		case jsonparser.Boolean:
			b, err := strconv.ParseBool(string(value))
			if err == nil {
				child := NewBoolObject(b)
				obj.arrChildren = append(obj.arrChildren, child)
			}
		case jsonparser.Null:
			child := NewNullObject()
			obj.arrChildren = append(obj.arrChildren, child)
		}
		return
	})
	return nil
}

// ====================
// content access

// simple values
func (obj *JsonObj) String() string {
	if obj.valueType == String {
		return obj.stringValue
	}
	return ""
}

func (obj *JsonObj) Int() int64 {
	if obj.valueType == Number {
		return obj.intValue
	}
	return 0
}

func (obj *JsonObj) Float() float64 {
	if obj.valueType == Number {
		return obj.floatValue
	}
	return 0.0
}

func (obj *JsonObj) Bool() bool {
	if obj.valueType == Boolean {
		return obj.boolValue
	}
	return false
}

// types
func (obj *JsonObj) Type() ValueType {
	return ValueType(obj.valueType)
}

func (obj *JsonObj) TypeString() string {
	switch obj.valueType {
	case String:
		return "string"
	case Number:
		return "number"
	case Boolean:
		return "boolean"
	case Null:
		return "null"
	case Object:
		return "object"
	case Array:
		return "array"
	default:
		return "unknown"
	}
}

func (obj *JsonObj) IsNull() bool {
	return obj.valueType == Null
}

func (obj *JsonObj) IsString() bool {
	return obj.valueType == String
}

func (obj *JsonObj) IsNumber() bool {
	return obj.valueType == Number
}

func (obj *JsonObj) IsObject() bool {
	return obj.valueType == Object
}

func (obj *JsonObj) IsArray() bool {
	return obj.valueType == Object
}

func (obj *JsonObj) IsBoollean() bool {
	return obj.valueType == Boolean
}

func (obj *JsonObj) IsBool() bool {
	return obj.valueType == Boolean
}

// children access
func (obj *JsonObj) GetByKey(keys ...string) (*JsonObj, error) {
	if obj.valueType != Object {
		return nil, NotAnObjectError
	}
	if 0 == len(keys) {
		return obj, nil
	}
	if 1 == len(keys) {
		child, exist := obj.objChildren[keys[0]]
		if false == exist {
			return nil, ObjectNotFoundError
		} else {
			return child, nil
		}
	}
	// else
	first_child, exist := obj.objChildren[keys[0]]
	if false == exist {
		return nil, ObjectNotFoundError
	}
	child, err := first_child.GetByKey(keys[1:]...)
	if err != nil {
		return nil, ObjectNotFoundError
	} else {
		return child, nil
	}
}

func (obj *JsonObj) GetAtIndex(index int) (*JsonObj, error) {
	if obj.valueType == Array {
		if index < len(obj.arrChildren) {
			return obj.arrChildren[index], nil
		} else {
			return nil, IndexOutOfBoundsError
		}
	} else {
		return nil, NotAnArrayError
	}
}
