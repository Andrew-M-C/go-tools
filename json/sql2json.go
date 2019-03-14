package json

import (
	"time"
	"strings"
	"strconv"
	"reflect"
	"bytes"
	"fmt"
	"github.com/Andrew-M-C/go-tools/str"
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
		fStr := fmt.Sprintf("%f", f)
		fStr = strings.TrimRight(fStr, "0")
		fStr = strings.TrimRight(fStr, ".")
		return fStr
	} else {
		format := fmt.Sprintf("%%.%df", int(digits))
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

func processFieldInt64(tag string, i int64, valid bool, opt Option, keyList *[]string, valList *[]string) {
	if valid {
		*keyList = append(*keyList, tag)
		*valList = append(*valList, strconv.FormatInt(i, 10))
	} else {
		if false == opt.OmitNull {
			*keyList = append(*keyList, tag)
			*valList = append(*valList, "null")
		}
	}
}

func processFieldString(tag string, s string, valid bool, opt Option, keyList *[]string, valList *[]string) {
	if valid {
		*keyList = append(*keyList, tag)
		*valList = append(*valList, `"` + escapeJsonString(s) + `"`)
	} else {
		if false == opt.OmitNull {
			*keyList = append(*keyList, tag)
			*valList = append(*valList, "null")
		}
	}
}

func processFieldBool(tag string, b bool, valid bool, opt Option, keyList *[]string, valList *[]string) {
	if valid {
		*keyList = append(*keyList, tag)
		if b {
			*valList = append(*valList, "true")
		} else {
			*valList = append(*valList, "false")
		}
	} else {
		if false == opt.OmitNull {
			*keyList = append(*keyList, tag)
			*valList = append(*valList, "null")
		}
	}
}

func processFieldFloat64(tag string, f float64, valid bool, opt Option, keyList *[]string, valList *[]string) {
	if valid {
		*keyList = append(*keyList, tag)
		*valList = append(*valList, convertFloatToString(f, opt.FloatDigits))
	} else {
		if false == opt.OmitNull {
			*keyList = append(*keyList, tag)
			*valList = append(*valList, "null")
		}
	}
}

func processFieldTime(tag string, t time.Time, valid bool, opt Option, keyList *[]string, valList *[]string) {
	if valid {
		*keyList = append(*keyList, tag)
		*valList = append(*valList, `"` + convertTimeToString(t, opt.TimeDigits) + `"`)
	} else {
		if false == opt.OmitNull {
			*keyList = append(*keyList, tag)
			*valList = append(*valList, "null")
		}
	}
}

func processField(field *reflect.StructField, value reflect.Value, opt Option, keyList *[]string, valList *[]string) {
	tag := getFieldTag(field)
	if tag == "" {		// skip ignored fields
		return
	}

	// log.Debug("Tag: %s - %s", field.Name, tag)
	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		processFieldInt64(tag, value.Int(), true, opt, keyList, valList)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		processFieldInt64(tag, int64(value.Uint()), true, opt, keyList, valList)
	case reflect.String:
		processFieldString(tag, value.String(), true, opt, keyList, valList)
	case reflect.Bool:
		processFieldBool(tag, value.Bool(), true, opt, keyList, valList)
	case reflect.Float32, reflect.Float64:
		processFieldFloat64(tag, value.Float(), true, opt, keyList, valList)
	case reflect.Struct:
		type_str := field.Type.String()
		if type_str == "sql.NullString" {
			processFieldString(tag, value.Field(0).String(), value.Field(1).Bool(), opt, keyList, valList)
		} else if type_str == "sql.NullInt64" {
			processFieldInt64(tag, value.Field(0).Int(), value.Field(1).Bool(), opt, keyList, valList)
		} else if type_str == "sql.NullBool" {
			processFieldBool(tag, value.Field(0).Bool(), value.Field(1).Bool(), opt, keyList, valList)
		} else if type_str == "sql.NullFloat64" {
			processFieldFloat64(tag, value.Field(0).Float(), value.Field(1).Bool(), opt, keyList, valList)
		} else if type_str == "mysql.NullTime" {
			timeVal, _ := value.Field(0).Interface().(time.Time)
			processFieldTime(tag, timeVal, value.Field(1).Bool(), opt, keyList, valList)
		} else if type_str == "time.Time" {
			timeVal, _ := value.Interface().(time.Time)
			processFieldTime(tag, timeVal, true, opt, keyList, valList)
		} else {
			// log.Debug("Unrecognized type: %s (%d)", field.Type.String(), int(field.Type.Kind()))
			// ignore
		}
	default:
		// log.Debug("Unrecognized type: %s (%d)", field.Type.String(), int(field.Type.Kind()))
		// ignore
	}
}

func sqlTypeToJson(t reflect.Type, v reflect.Value, opt Option) (string, error) {
	num_field := t.NumField()
	key_list := make([]string, 0, num_field)
	value_list := make([]string, 0, num_field)

	// parse struct
	for i := 0; i < num_field; i ++ {
		field := t.Field(i)
		processField(&field, v.Field(i), opt, &key_list, &value_list)
	}

	// get json string
	b := bytes.Buffer{}
	is_first := true
	b.WriteRune('{')
	for i, value := range value_list {
		if false == is_first {
			b.WriteRune(',')
		} else {
			is_first = false
		}
		key := key_list[i]
		// log.Debug("\"%s\" = %s", key, value)
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
