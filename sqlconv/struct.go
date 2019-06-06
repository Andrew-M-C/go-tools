package sqlconv
import (
	"github.com/Andrew-M-C/go-tools/str"
	"github.com/Andrew-M-C/go-tools/log"
	"reflect"
	"bytes"
	"strconv"
	"time"
	"strings"
)

func GetValidKVsFromStruct(s interface{}, quote string) ([]string, []string, error) {
	// check parameter type
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	switch(t.Kind()) {
	case reflect.Ptr:
		return getKVs(t.Elem(), v.Elem(), quote)
	case reflect.Struct:
		return getKVs(t, v, quote)
	default:
		return nil, nil, DataTypeError
	}
}

func getKVs(t reflect.Type, v reflect.Value, quote string) ([]string, []string, error) {
	num_field := t.NumField()
	key_list := make([]string, 0, num_field)
	value_list := make([]string, 0, num_field)

	for i := 0; i < num_field; i ++ {
		field := t.Field(i)
		tag := getFieldTag(&field)
		value := v.Field(i)
		if str.Empty(tag) {
			continue
		}

		log.Debug("Tag: %s - %s", field.Name, tag)
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			key_list = append(key_list, tag)
			value_list = append(value_list, strconv.FormatInt(value.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			key_list = append(key_list, tag)
			value_list = append(value_list, strconv.FormatUint(value.Uint(), 10))
		case reflect.String:
			key_list = append(key_list, tag)
			value_list = append(value_list, escapeStringWithQuote(value.String(), quote))
		case reflect.Bool:
			key_list = append(key_list, tag)
			value_list = append(value_list, strconv.FormatBool(value.Bool()))
		case reflect.Float32, reflect.Float64:
			key_list = append(key_list, tag)
			value_list = append(value_list, strconv.FormatFloat(value.Float(), 'f', 6, 64))
		case reflect.Struct:
			type_str := field.Type.String()
			if type_str == "sql.NullString" {
				if value.Field(1).Bool() {
					key_list = append(key_list, tag)
					value_list = append(value_list, escapeStringWithQuote(value.Field(0).String(), quote))
				}
			} else if type_str == "sql.NullInt64" {
				if value.Field(1).Bool() {
					key_list = append(key_list, tag)
					value_list = append(value_list, strconv.FormatInt(value.Field(0).Int(), 10))
				}
			} else if type_str == "sql.NullBool" {
				if value.Field(1).Bool() {
					key_list = append(key_list, tag)
					value_list = append(value_list, strconv.FormatBool(value.Field(0).Bool()))
				}
			} else if type_str == "sql.NullFloat64" {
				if value.Field(1).Bool() {
					key_list = append(key_list, tag)
					value_list = append(value_list, strconv.FormatFloat(value.Field(0).Float(), 'f', 6, 64))
				}
			} else if type_str == "mysql.NullTime" {
				if value.Field(1).Bool() {
					t, _ := value.Field(0).Interface().(time.Time)
					key_list = append(key_list, tag)
					value_list = append(value_list, getTimeString(t))
				}
			} else if type_str == "time.Time" {
				t, _ := value.Interface().(time.Time)
				key_list = append(key_list, tag)
				value_list = append(value_list, getTimeString(t))
			} else {
				log.Debug("Unrecognized type: %s (%d)", field.Type.String(), int(field.Type.Kind()))
				// ignore
			}
		default:
			log.Debug("Unrecognized type: %s (%d)", field.Type.String(), int(field.Type.Kind()))
			// ignore
		}
	}

	return key_list, value_list, nil
}


func getFieldTag(field *reflect.StructField) string {
	tag := ""

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

	// done
	return tag
}


func getTimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}


func escapeStringWithQuote(s string, quote string) string {
	is_single_quote := false
	is_double_quote := false
	is_back_quote := false
	if quote == "'" {
		is_single_quote = true
	} else if quote == "\"" {
		is_double_quote = true
	} else if quote == "`" {
		is_back_quote = true
	} else {
		return s
	}

	buff := bytes.Buffer{}

	// string start
	buff.WriteString(quote)

	// string body
	for _, c := range s {
		switch c {
		case '\'':
			if is_single_quote {
				buff.WriteString("''")
			} else {
				buff.WriteRune('\'')
			}
		case '"':
			if is_double_quote {
				buff.WriteString("\\\"")
			} else {
				buff.WriteRune('"')
			}
		case '&':
			buff.WriteString("&&")
		case '`':
			if is_back_quote {
				buff.WriteString("\\`")
			} else {
				buff.WriteRune('`')
			}
		case '\\':
			buff.WriteString("\\\\")
		default:
			buff.WriteRune(c)
		}
	}

	// string end
	buff.WriteString(quote)

	return buff.String()
}
