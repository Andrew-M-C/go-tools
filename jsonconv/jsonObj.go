package jsonconv

import (
	"github.com/buger/jsonparser"
	// "github.com/Andrew-M-C/go-tools/log"
	"strconv"
	"strings"
	"bytes"
	"reflect"
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
// internal functions

var escapeMap = map[rune]rune {
	'"': '"',
	'/': '/',
	'b': '\b',
	'f': '\f',
	't': '\t',
	'n': '\n',
}

func stringFromEscapedBytes(input []byte) (string, error) {
	b := bytes.Buffer{}
	s := string(input)
	escaping := false
	skip := 0

	for i, chr := range s {
		if skip > 0 {
			skip --
		} else if escaping {
			escaping = false
			switch chr {
			case '"', '/', 'b', 'f', 't', 'n', 'r':
				write_chr, exist := escapeMap[chr]
				if exist {
					b.WriteRune(write_chr)
				}
			case 'u':
				// parse unicode
				sub_str := s[i+1:i+5]
				unicode, err := strconv.ParseInt(sub_str, 16, 32)
				if err != nil {
					// err
					// log.Error("err: %s", err.Error())
					return "", JsonFormatError
				} else {
					skip = 4
					b.WriteRune(rune(unicode))
				}
			default:
				// illegal character
				// log.Error("Illegal escaped char: %c", chr)
				return "", JsonFormatError
			}
		} else {
			switch chr {
			case '\\':
				// get ready to escape
				escaping = true
			default:
				b.WriteRune(chr)
			}
		}
	}
	return b.String(), nil
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
			// log.Debug("Skip: %c", chr)
			// skip
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
	add_child := func (obj *JsonObj, key []byte, child *JsonObj) {
		key_str, key_err := stringFromEscapedBytes(key)
		if key_err == nil {
			obj.objChildren[key_str] = child
		}
	}

	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, _ int) error {
		// log.Debug("key: %s", string(key))
		// log.Debug("value: %s", string(value))
		switch dataType {
		case jsonparser.String:
			str_value, err := stringFromEscapedBytes(value)
			if err == nil {
				child := NewStringObject(str_value)
				add_child(obj, key, child)
			}
		case jsonparser.Number:
			child := new(JsonObj)
			child.valueType = Number
			child.origBytes = value
			add_child(obj, key, child)
		case jsonparser.Object:
			child := NewObjectObject()
			err := child.parseObject(value)
			if err == nil {
				add_child(obj, key, child)
			}
		case jsonparser.Array:
			child := NewArrayObject()
			err := child.parseArray(value)
			if err == nil {
				add_child(obj, key, child)
			}
		case jsonparser.Boolean:
			b, err := strconv.ParseBool(string(value))
			if err == nil {
				child := NewBoolObject(b)
				add_child(obj, key, child)
			}
		case jsonparser.Null:
			child := NewNullObject()
			add_child(obj, key, child)
		default:
			// log.Debug("Invalid type: %d", int(dataType))
		}
		return nil
	})
	return nil
}

func (obj *JsonObj) parseArray(data []byte) error {
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, _ int, _ error) {
		switch dataType {
		case jsonparser.String:
			str_value, err := stringFromEscapedBytes(value)
			if err == nil {
				child := NewStringObject(str_value)
				obj.arrChildren = append(obj.arrChildren, child)
			}
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
	if index < 0 {
		return nil, ParaError
	} else if obj.valueType == Array {
		if index < len(obj.arrChildren) {
			return obj.arrChildren[index], nil
		} else {
			return nil, IndexOutOfBoundsError
		}
	} else {
		return nil, NotAnArrayError
	}
}

func (obj *JsonObj) Get(first interface{}, keys ...interface{}) (*JsonObj, error) {
	switch first.(type) {
	case string:
		child, err := obj.GetByKey(first.(string))
		if err != nil {
			return nil, err
		} else if len(keys) == 1 {
			return child.Get(keys[0])
		} else if len(keys) > 1 {
			return child.Get(keys[0], keys[1:]...)
		} else {
			return child, nil
		}
	case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint:
		value := reflect.ValueOf(first)
		index := int(value.Int())
		child, err := obj.GetAtIndex(index)
		if err != nil {
			return nil, err
		} else if len(keys) == 1 {
			return child.Get(keys[0])
		} else if len(keys) > 1 {
			return child.Get(keys[0], keys[1:]...)
		} else {
			return child, nil
		}
	default:
		return nil, ParaError
	}
}

func (obj *JsonObj) GetString(first interface{}, keys ...interface{}) (string, error) {
	child, err := obj.Get(first, keys...)
	if err != nil {
		return "", err
	}
	if false == child.IsString() {
		return "", NotAStringError
	}
	return child.String(), nil
}

func (obj *JsonObj) GetInt(first interface{}, keys ...interface{}) (int64, error) {
	child, err := obj.Get(first, keys...)
	if err != nil {
		return 0, err
	}
	if false == child.IsNumber() {
		return 0, NotANumberError
	}
	return child.Int(), nil
}

func (obj *JsonObj) GetFloat(first interface{}, keys ...interface{}) (float64, error) {
	child, err := obj.Get(first, keys...)
	if err != nil {
		return 0.0, err
	}
	if false == child.IsNumber() {
		return 0.0, NotANumberError
	}
	return child.Float(), nil
}

func (obj *JsonObj) GetBool(first interface{}, keys ...interface{}) (bool, error) {
	child, err := obj.Get(first, keys...)
	if err != nil {
		return false, err
	}
	if false == child.IsBool() {
		return false, NotABoolError
	}
	return child.Bool(), nil
}
