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

type JsonValue struct {
	// type
	valueType		jsonparser.ValueType
	// values
	stringValue		string
	intValue		int64
	floatValue		float64
	boolValue		bool
	// object children
	objChildren		map[string]*JsonValue
	// array children
	arrChildren		[]*JsonValue
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

func NewByUnmarshal(s string) (*JsonValue, error) {
	return NewFromString(s)
}

func NewFromString(s string) (*JsonValue, error) {
	var obj *JsonValue
	var err error

	// get first character
	for index, chr := range s {
		switch chr {
		case ' ', '\r', '\n', '\t':
			// continue
		case '{':
			obj = NewObject()
			err = obj.parseObject([]byte(s[index:]))
			if err != nil {
				return nil, err
			} else {
				return obj, nil
			}
		case '[':
			obj = NewArray()
			err = obj.parseArray([]byte(s[index:]))
			if err != nil {
				return nil, err
			} else {
				return obj, nil
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			obj = new(JsonValue)
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
			obj = new(JsonValue)
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
				obj = new(JsonValue)
				obj.valueType = Boolean
				obj.boolValue = true
				return obj, nil
			} else {
				return nil, JsonFormatError
			}
		case 'f':
			if s[index:] == "false" {
				obj = new(JsonValue)
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

func NewString(s string) *JsonValue {
	obj := new(JsonValue)
	obj.valueType = String
	obj.stringValue = s
	return obj
}

func NewInt(i int64) *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Number
	obj.intValue = i
	obj.floatValue = float64(i)
	return obj
}

func NewFloat(f float64) *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Number
	obj.intValue = int64(f)
	obj.floatValue = f
	return obj
}

func NewBoolObject(b bool) *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Boolean
	obj.boolValue = b
	return obj
}

func NewNull() *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Null
	return obj
}

func NewObject() *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Object
	obj.objChildren = make(map[string]*JsonValue)
	return obj
}

func NewArray() *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Array
	obj.arrChildren = make([]*JsonValue, 0, 10)
	return obj
}

// ====================
// parse functions

func (obj *JsonValue) parseObject(data []byte) error {
	add_child := func (obj *JsonValue, key []byte, child *JsonValue) {
		key_str, key_err := stringFromEscapedBytes(key)
		if key_err == nil {
			obj.objChildren[key_str] = child
		}
	}

	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, _ int) error {
		// log.Debug("----------")
		// log.Debug("key: %s", string(key))
		// log.Debug("value: %s", string(value))
		switch dataType {
		case jsonparser.String:
			// log.Debug("string")
			str_value, err := stringFromEscapedBytes(value)
			if err == nil {
				child := NewString(str_value)
				add_child(obj, key, child)
			}
		case jsonparser.Number:
			// log.Debug("number")
			child := new(JsonValue)
			str_value := string(value)
			child.valueType = Number
			child.intValue, _ = strconv.ParseInt(str_value, 10, 64)
			child.floatValue, _ = strconv.ParseFloat(str_value, 64)
			add_child(obj, key, child)
		case jsonparser.Object:
			// log.Debug("object")
			child := NewObject()
			err := child.parseObject(value)
			if err != nil {
				// log.Error("Failed to parse object: %s", err.Error())
			} else {
				// log.Debug("%s ---- object size: %d", string(key), child.Length())
				add_child(obj, key, child)
			}
		case jsonparser.Array:
			// log.Debug("array")
			child := NewArray()
			err := child.parseArray(value)
			if err != nil {
				// log.Error("Failed to parse array: %s", err.Error())
			} else {
				// log.Debug("%s ---- array size: %d", string(key), child.Length())
				add_child(obj, key, child)
			}
		case jsonparser.Boolean:
			// log.Debug("bool")
			b, err := strconv.ParseBool(string(value))
			if err == nil {
				child := NewBoolObject(b)
				add_child(obj, key, child)
			}
		case jsonparser.Null:
			// log.Debug("null")
			child := NewNull()
			add_child(obj, key, child)
		default:
			// log.Debug("Invalid type: %d", int(dataType))
		}
		return nil
	})
	return nil
}

func (obj *JsonValue) parseArray(data []byte) error {
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, _ int, _ error) {
		switch dataType {
		case jsonparser.String:
			str_value, err := stringFromEscapedBytes(value)
			if err == nil {
				child := NewString(str_value)
				obj.arrChildren = append(obj.arrChildren, child)
			}
		case jsonparser.Number:
			child := new(JsonValue)
			str_value := string(value)
			child.valueType = Number
			child.intValue, _ = strconv.ParseInt(str_value, 10, 64)
			child.floatValue, _ = strconv.ParseFloat(str_value, 64)
			obj.arrChildren = append(obj.arrChildren, child)
		case jsonparser.Object:
			child := NewObject()
			err := child.parseObject(value)
			if err == nil {
				obj.arrChildren = append(obj.arrChildren, child)
			}
		case jsonparser.Array:
			child := NewArray()
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
			child := NewNull()
			obj.arrChildren = append(obj.arrChildren, child)
		}
		return
	})
	return nil
}

// ====================
// content access

// simple values
func (obj *JsonValue) String() string {
	if obj.valueType == String {
		return obj.stringValue
	}
	return ""
}

func (obj *JsonValue) Int() int64 {
	if obj.valueType == Number {
		return obj.intValue
	}
	return 0
}

func (obj *JsonValue) Float() float64 {
	if obj.valueType == Number {
		return obj.floatValue
	}
	return 0.0
}

func (obj *JsonValue) Bool() bool {
	if obj.valueType == Boolean {
		return obj.boolValue
	}
	return false
}

func (obj *JsonValue) Boolean() bool {
	if obj.valueType == Boolean {
		return obj.boolValue
	}
	return false
}

func (obj *JsonValue) Length() int {
	if obj.valueType == Array {
		return len(obj.arrChildren)
	} else if obj.valueType == Object {
		return len(obj.objChildren)
	} else {
		return 0
	}
}

// types
func (obj *JsonValue) Type() ValueType {
	return ValueType(obj.valueType)
}

func (obj *JsonValue) TypeString() string {
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

func (obj *JsonValue) IsNull() bool {
	return obj.valueType == Null
}

func (obj *JsonValue) IsString() bool {
	return obj.valueType == String
}

func (obj *JsonValue) IsNumber() bool {
	return obj.valueType == Number
}

func (obj *JsonValue) IsObject() bool {
	return obj.valueType == Object
}

func (obj *JsonValue) IsArray() bool {
	return obj.valueType == Array
}

func (obj *JsonValue) IsBoollean() bool {
	return obj.valueType == Boolean
}

func (obj *JsonValue) IsBool() bool {
	return obj.valueType == Boolean
}

// children access
func (obj *JsonValue) GetByKey(keys ...string) (*JsonValue, error) {
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

func (obj *JsonValue) GetAtIndex(index int) (*JsonValue, error) {
	if obj.valueType == Array {
		if index >= 0 && index < obj.Length() {
			return obj.arrChildren[index], nil
		} else {
			// log.Error("request index %d, but length is %d", index, obj.Length())
			return nil, IndexOutOfBoundsError
		}
	} else {
		return nil, NotAnArrayError
	}
}

func (obj *JsonValue) Get(first interface{}, keys ...interface{}) (*JsonValue, error) {
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

func (obj *JsonValue) GetString(first interface{}, keys ...interface{}) (string, error) {
	child, err := obj.Get(first, keys...)
	if err != nil {
		return "", err
	}
	if false == child.IsString() {
		return "", NotAStringError
	}
	return child.String(), nil
}

func (obj *JsonValue) GetInt(first interface{}, keys ...interface{}) (int64, error) {
	child, err := obj.Get(first, keys...)
	if err != nil {
		return 0, err
	}
	if false == child.IsNumber() {
		return 0, NotANumberError
	}
	return child.Int(), nil
}

func (obj *JsonValue) GetFloat(first interface{}, keys ...interface{}) (float64, error) {
	child, err := obj.Get(first, keys...)
	if err != nil {
		return 0.0, err
	}
	if false == child.IsNumber() {
		return 0.0, NotANumberError
	}
	return child.Float(), nil
}

func (obj *JsonValue) GetBool(first interface{}, keys ...interface{}) (bool, error) {
	child, err := obj.Get(first, keys...)
	if err != nil {
		return false, err
	}
	if false == child.IsBool() {
		return false, NotABoolError
	}
	return child.Bool(), nil
}

// ====================
// Marshal
func Marshal(obj *JsonValue, opts... Option) (string, error) {
	if nil == obj {
		return "", ParaError
	}
	return obj.Marshal(opts...)
}

func (obj *JsonValue)Marshal(opts... Option) (string, error) {
	var opt *Option
	if len(opts) > 0 {
		opt = &(opts[0])
	} else {
		opt = &dftOption
	}

	switch obj.valueType {
	case String:
		return `"` + escapeJsonString(obj.String(), opt.EnsureAscii) + `"`, nil
	case Number:
		i := obj.Int()
		f := obj.Float()
		if float64(i) == f {
			return strconv.FormatInt(i, 10), nil
		} else {
			return convertFloatToString(f, opt.FloatDigits), nil
		}
	case Null:
		return "null", nil
	case Boolean:
		if obj.Bool() {
			return "true", nil
		} else {
			return "false", nil
		}
	case Object:
		is_first := true
		b := bytes.Buffer{}
		b.WriteRune('{')
		for key, child := range obj.objChildren {
			if child.IsNull() && false == opt.ShowNull {
				// do nothing
			} else {
				if is_first {
					is_first = false
				} else {
					b.WriteRune(',')
				}
				b.WriteRune('"')
				b.WriteString(escapeJsonString(key, opt.EnsureAscii))
				b.WriteRune('"')
				b.WriteRune(':')

				child_str, _ := child.Marshal(*opt)
				b.WriteString(child_str)
			}
		}
		b.WriteRune('}')
		return b.String(), nil
	case Array:
		is_first := true
		b := bytes.Buffer{}
		b.WriteRune('[')
		for _, child := range obj.arrChildren {
			if child.IsNull() && false == opt.ShowNull {
				// do nothing
			} else {
				if is_first {
					is_first = false
				} else {
					b.WriteRune(',')
				}
				child_str, _ := child.Marshal(*opt)
				b.WriteString(child_str)
			}
		}
		b.WriteRune(']')
		return b.String(), nil
	default:
		// do nothing
		return "", JsonTypeError
	}
}

// ====================
// object modification
func (obj *JsonValue) Delete(first interface{}, keys ...interface{}) error {
	var parent *JsonValue
	var last_key *interface{}
	var err error

	// get parent
	switch len(keys) {
	case 0:
		last_key = &first
		parent = obj
	case 1:
		last_key = &keys[0]
		parent, err = obj.Get(first)
		if err != nil {
			return ObjectNotFoundError
		}
	default:
		last_index := len(keys) - 1
		last_key = &keys[last_index]
		parent, err = obj.Get(first, keys[1:]...)
		if err != nil {
			return ObjectNotFoundError
		}
	}

	// get child
	switch (*last_key).(type) {
	case string:
		if false == parent.IsObject() {
			return ObjectNotFoundError
		}
		// delect object
		last_key_str := (*last_key).(string)
		// log.Debug("Delete key: %s", last_key_str)
		_, err = parent.Get(last_key_str)
		if err != nil {
			// log.Debug("key %s not found", last_key_str)
			return err
		}
		delete(parent.objChildren, last_key_str)
		return nil

	case uint8, int8, uint16, int16, uint32, int32, uint64, int64, int, uint:
		value := reflect.ValueOf(first)
		index := int(value.Int())
		arr_len := len(parent.arrChildren)
		if index >= 0 && index < arr_len {
			parent.arrChildren = append(parent.arrChildren[index:], parent.arrChildren[index+1:]...)
			return nil
		} else {
			return IndexOutOfBoundsError
		}

	default:
		return DataTypeError
	}
}

func (this *JsonValue) Append(newOne *JsonValue, keys ...interface{}) error {
	if nil == newOne {
		return nil
	}
	if 0 == len(keys) {
		if this.valueType == Array {
			this.arrChildren = append(this.arrChildren, newOne)
			return nil
		} else {
			return NotAnArrayError
		}
	} else {
		child, err := this.Get(keys[0], keys[1:]...)
		if err != nil {
			return err
		} else {
			return child.Append(newOne)
		}
	}
}

func (this *JsonValue) Insert(newOne *JsonValue, index interface{}, keys ...interface{}) error {
	if nil == newOne {
		return nil
	}
	keys_count := len(keys)
	if 0 == keys_count {
		if this.valueType != Array {
			return NotAnArrayError
		}
		switch index.(type) {
		case uint8, int8, uint16, int16, uint32, int32, uint64, int64, int, uint:
			value := reflect.ValueOf(index)
			index := int(value.Int())
			arr_len := len(this.arrChildren)
			if index >= 0 && index < arr_len {
				rear := this.arrChildren[index:]
				this.arrChildren = append(this.arrChildren[0:index], newOne)
				this.arrChildren = append(this.arrChildren, rear...)
				return nil
			} else {
				return IndexOutOfBoundsError
			}
		default:
			return DataTypeError

		}
	} else {
		var err error
		var child *JsonValue
		if 1 == keys_count {
			child, err = this.Get(index, keys[0])
		} else {
			child, err = this.Get(index, keys[:keys_count-2]...)
		}
		if err != nil {
			return err
		}
		return child.Insert(newOne, keys[keys_count-1])
	}
}

func (this *JsonValue) Set(newOne *JsonValue, first interface{}, keys ...interface{}) error {
	// log.Debug("Set \"%v\" (%v)", first, keys)
	keys_count := len(keys)
	switch keys_count {
	case 0:
		switch first.(type) {
		case string:
			if this.IsObject() {
				key := first.(string)
				this.objChildren[key] = newOne
				return nil
			} else {
				// log.Error("Not an object")
				return NotAnObjectError
			}
		default:
			// log.Error("leaf not a string")
			return DataTypeError
		}
	case 1:
		child, err := this.Get(first)
		if err != nil {
			// log.Error("Failed to get: %s", err.Error())
			return err
		}
		return child.Set(newOne, keys[0])
	default:
		child, err := this.Get(first)
		if err != nil {
			// log.Error("Failed to get: %s", err.Error())
			return err
		}
		return child.Set(newOne, keys[0], keys[1:]...)
	}
}

// ====================
// foreach
func (this *JsonValue) ArrayForeach(callback func(index int, value *JsonValue) error) error {
	if false == this.IsArray() {
		// log.Error("object is not an array")
		return NotAnArrayError
	}
	// log.Debug("array size: %d", len(this.arrChildren))
	for i, val := range this.arrChildren {
		err := callback(i, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *JsonValue) ObjectForeach(callback func(key string, value *JsonValue) error) error {
	if false == this.IsObject() {
		return NotAnObjectError
	}
	for k, v := range this.objChildren {
		err := callback(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
