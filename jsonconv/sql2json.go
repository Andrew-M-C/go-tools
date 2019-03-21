package jsonconv

import (
	"time"
	"strings"
	"strconv"
	"reflect"
	"bytes"
	"github.com/Andrew-M-C/go-tools/str"
	// "github.com/Andrew-M-C/go-tools/log"
	// "database/sql"
	// "github.com/go-sql-driver/mysql"
)

func getFieldTag(field *reflect.StructField, filterMode Filter, filterMap map[string]int, ensureAscii bool) string {
	tag := ""

	// read from "json"
	{
		tag_list := strings.SplitN(field.Tag.Get("json"), ",", 2)
		if str.Valid(tag_list[0]) {
			tag = tag_list[0]
			if tag == "-" {
				return ""
			}
		}
	}

	// read from "db"
	if str.Empty(tag) {
		tag_list := strings.SplitN(field.Tag.Get("db"), ",", 2)
		if str.Valid(tag_list[0]) {
			tag = tag_list[0]
		}
	}

	// read from field name
	if str.Empty(tag) {
		tag = field.Name
	}

	// filter
	switch filterMode {
	case IncludeMode:
		_, exist := filterMap[tag]
		if false == exist {
			return ""
		}
	case ExcludeMode:
		_, exist := filterMap[tag]
		if exist {
			return ""
		}
	default:
		// do nothing
	}

	// done
	return escapeJsonString(tag, ensureAscii)
}

func processFieldInt64(tag string, i int64, valid bool, opt Option, keyList *[]string, valList *[]string) {
	if valid {
		*keyList = append(*keyList, tag)
		*valList = append(*valList, strconv.FormatInt(i, 10))
	} else {
		if opt.ShowNull {
			*keyList = append(*keyList, tag)
			*valList = append(*valList, "null")
		}
	}
}

func processFieldString(tag string, s string, valid bool, opt Option, keyList *[]string, valList *[]string) {
	if valid {
		*keyList = append(*keyList, tag)
		*valList = append(*valList, `"` + escapeJsonString(s, opt.EnsureAscii) + `"`)
	} else {
		if opt.ShowNull {
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
		if opt.ShowNull {
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
		if opt.ShowNull {
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
		if opt.ShowNull {
			*keyList = append(*keyList, tag)
			*valList = append(*valList, "null")
		}
	}
}

func processField(field *reflect.StructField, value reflect.Value, opt Option, filterMap map[string]int, keyList *[]string, valList *[]string) {
	tag := getFieldTag(field, opt.FilterMode, filterMap, opt.EnsureAscii)
	if str.Empty(tag) {		// skip ignored fields
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

	// parse filter mode
	filter_map := make(map[string]int)
	if opt.FilterMode == IncludeMode || opt.FilterMode == ExcludeMode {
		for _, key := range opt.FilterList {
			filter_map[key] = 1
		}
	}

	// parse struct
	for i := 0; i < num_field; i ++ {
		field := t.Field(i)
		processField(&field, v.Field(i), opt, filter_map, &key_list, &value_list)
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
