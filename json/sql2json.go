package json

import (
	"time"
	"strings"
	"strconv"
	"reflect"
	"bytes"
	"fmt"
	"../log"
	"../str"
	// "database/sql"
	// "github.com/go-sql-driver/mysql"
)

type Option struct {
	OmitNull		bool
	FloatDigits		uint8
	TimeDigits		uint8
}

var dftOption = Option {
	OmitNull: true,
	FloatDigits: 0,
	TimeDigits: 0,
}

func escapeJsonString(s string) string {
	b := bytes.Buffer{}
	for _, chr := range s {
		switch chr {
		case '"':
			b.WriteString("\\\"")
		case '/':
			b.WriteString("\\/")
		case '\b':
			b.WriteString("\\b")
		case '\f':
			b.WriteString("\\f")
		case '\t':
			b.WriteString("\\t")
		case '\n':
			b.WriteString("\\n")
		case '\r':
			b.WriteString("\\r")
		default:
			b.WriteRune(chr)
		}
	}
	return b.String()
}

func convertFloatToString(f float64, digits uint8) string {
	if 0 == digits {
		return fmt.Sprintf("%f", f)
	} else {
		format := fmt.Sprintf("%%df", digits)
		return fmt.Sprintf(format, f)
	}
}

func convertTimeToString(t time.Time, digits uint8) string {
	if 0 == digits {
		return t.Format("2006-01-02 15:04:05")
	} else {
		postfix := strings.Repeat("0", int(digits))
		return t.Format("2006-01-02 15:04:05." + postfix)
	}
}

func getFieldTag(field *reflect.StructField) string {
	var tag_list []string

	// read from "json"
	tag_list = strings.SplitN(field.Tag.Get("json"), ",", 2)
	if str.Valid(tag_list[0]) {
		tag := tag_list[0]
		if tag == "-" {
			return ""
		}
		return escapeJsonString(tag)
	}

	// read from "db"
	tag_list = strings.SplitN(field.Tag.Get("db"), ",", 2)
	if str.Valid(tag_list[0]) {
		tag := tag_list[0]
		return escapeJsonString(tag)
	}

	// read field name
	return escapeJsonString(field.Name)
}

func processFieldInt64(tag string, i int64, valid bool, opt Option, keyValue map[string]string) {
	if valid {
		keyValue[tag] = strconv.FormatInt(i, 10)
	} else {
		if false == opt.OmitNull {
			keyValue[tag] = "null"
		}
	}
}

func processFieldString(tag string, s string, valid bool, opt Option, keyValue map[string]string) {
	if valid {
		keyValue[tag] = `"` + escapeJsonString(s) + `"`
	} else {
		if false == opt.OmitNull {
			keyValue[tag] = "null"
		}
	}
}

func processFieldBool(tag string, b bool, valid bool, opt Option, keyValue map[string]string) {
	if valid {
		if b {
			keyValue[tag] = "true"
		} else {
			keyValue[tag] = "false"
		}
	} else {
		if false == opt.OmitNull {
			keyValue[tag] = "null"
		}
	}
}

func processFieldFloat64(tag string, f float64, valid bool, opt Option, keyValue map[string]string) {
	if valid {
		keyValue[tag] = convertFloatToString(f, opt.FloatDigits)
	} else {
		if false == opt.OmitNull {
			keyValue[tag] = "null"
		}
	}
}

func processFieldTime(tag string, t time.Time, valid bool, opt Option, keyValue map[string]string) {
	// TODO
	if valid {
		keyValue[tag] = convertTimeToString(t, opt.TimeDigits)
	} else {
		if false == opt.OmitNull {
			keyValue[tag] = "null"
		}
	}
}

func processField(field *reflect.StructField, value reflect.Value, opt Option, keyValue map[string]string) {
	tag := getFieldTag(field)
	if tag == "" {		// skip ignored fields
		return
	}

	// log.Debug("Tag: %s - %s", field.Name, tag)
	switch field.Type.Kind() {
		// integers
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		processFieldInt64(tag, value.Int(), true, opt, keyValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		processFieldInt64(tag, int64(value.Uint()), true, opt, keyValue)
	case reflect.String:
		processFieldString(tag, value.String(), true, opt, keyValue)
	case reflect.Bool:
		processFieldBool(tag, value.Bool(), true, opt, keyValue)
	case reflect.Float32, reflect.Float64:
		processFieldFloat64(tag, value.Float(), true, opt, keyValue)
	case reflect.Struct:
		type_str := field.Type.String()
		if type_str == "sql.NullString" {
			processFieldString(tag, value.Field(0).String(), value.Field(1).Bool(), opt, keyValue)
		} else if type_str == "sql.NullInt64" {
			processFieldInt64(tag, value.Field(0).Int(), value.Field(1).Bool(), opt, keyValue)
		} else if type_str == "sql.NullBool" {
			processFieldBool(tag, value.Field(0).Bool(), value.Field(1).Bool(), opt, keyValue)
		} else if type_str == "sql.NullFloat64" {
			processFieldFloat64(tag, value.Field(0).Float(), value.Field(1).Bool(), opt, keyValue)
		} else if type_str == "mysql.NullTime" {
			timeVal, _ := value.Field(0).Interface().(time.Time)
			processFieldTime(tag, timeVal, value.Field(1).Bool(), opt, keyValue)
		} else if type_str == "time.Time" {
			timeVal, _ := value.Interface().(time.Time)
			processFieldTime(tag, timeVal, true, opt, keyValue)
		} else {
			log.Debug("Unrecognized type: %s (%d)", field.Type.String(), int(field.Type.Kind()))
			// ignore
		}
	default:
		log.Debug("Unrecognized type: %s (%d)", field.Type.String(), int(field.Type.Kind()))
		// ignore
	}
}

func sqlTypeToJson(t reflect.Type, v reflect.Value, opt Option) (string, error) {
	num_field := t.NumField()
	key_value := make(map[string]string)

	// parse struct
	for i := 0; i < num_field; i ++ {
		field := t.Field(i)
		processField(&field, v.Field(i), opt, key_value)
	}

	// get json string
	b := bytes.Buffer{}
	is_first := true
	b.WriteRune('{')
	for key, value := range key_value {
		if false == is_first {
			b.WriteRune(',')
		} else {
			is_first = false
		}
		b.WriteRune('"')
		b.WriteString(key)
		b.WriteRune('"')
		b.WriteRune(':')
		b.WriteString(value)
	}
	b.WriteRune('}')

	return b.String(), nil
}

/**
 * Valid parameter type: struct ptr
 */
func SqlToJson(u interface{}, options... Option) (string, error) {
	// check parameter type
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)
	var opt *Option
	// log.Debug("type: %s", t.String())

	if len(options) > 0 {
		opt = &(options[0])
	} else {
		opt = &dftOption
	}

	switch(t.Kind()) {
	case reflect.Ptr:
		return sqlTypeToJson(t.Elem(), v.Elem(), *opt)
	case reflect.Struct:
		return sqlTypeToJson(t, v, *opt)
	default:
		return "", DataTypeError
	}
}

// type Kind uint
// const (
// 	Invalid Kind = iota	// 0
// 	Bool				// 1
// 	Int					// 2
// 	Int8				// 3
// 	Int16				// 4
// 	Int32				// 5
// 	Int64				// 6
// 	Uint				// 7
// 	Uint8				// 8
// 	Uint16				// 9
// 	Uint32				// 10
// 	Uint64				// 11
// 	Uintptr				// 12
// 	Float32				// 13
// 	Float64				// 14
// 	Complex64			// 15
// 	Complex128			// 16
// 	Array				// 17
// 	Chan				// 18
// 	Func				// 19
// 	Interface			// 20
// 	Map					// 21
// 	Ptr					// 22
// 	Slice				// 23
// 	String				// 24
// 	Struct				// 25
// 	UnsafePointer		// 26
// )
