package _test
import (
	"github.com/Andrew-M-C/go-tools/log"
	"github.com/Andrew-M-C/go-tools/jsonconv"
	"encoding/json"
	"github.com/buger/jsonparser"
	"strings"
	"strconv"
	"time"
)

type teststruct struct {
	S	string	`json:"str"`
	I	int64	`json:"int"`
	F	float64	`json:"float"`
	A	[]interface{}	`json:"array"`
	O	map[string]interface{}	`json:"obj"`
}

func elapsed(s time.Time) int64 {
	e := time.Now().Sub(s)
	usec := e.Nanoseconds()
	return usec
}

func TestOrigJsonEffenciency() (mapIntf, mapIntfAll, structure, parser, conv int64) {
	// This is an example for how to parse a structure-unknown JSON string
	// refs:
	//  - [preserve int64 values when parsing json in Go](https://stackoverflow.com/questions/16946306/preserve-int64-values-when-parsing-json-in-go)
	//  - [buger/jsonparser](https://github.com/buger/jsonparser)
	var err error

	var start time.Time
	json_bytes := []byte("{\"str\": \"hello, json\", \"int\": 123, \"float\": 10.2, \"array\": [1, \"2\"], \"obj\": {\"num\": 10}, \"bool\": true, \"int64\": 4418489049307132905}")
	var func_parse_array func([]interface{}, int)
	var func_parse_obj func(map[string]interface{}, int)

	func_parse_array = func(array []interface{} ,level int) {
		prefix := strings.Repeat("    ", level)
		for index, intf := range array {
			switch intf.(type) {
			case string:
				log.Debug(prefix + "[%d] = %s", index, intf.(string))
			case float64:
				integer := int64(intf.(float64))
				if (float64(integer) == intf.(float64)) {
					log.Debug(prefix + "[%d] = %d", index, integer)
				} else {
					log.Debug(prefix + "[%d] = %f", index, intf.(float64))
				}
			case bool:
				log.Debug(prefix + "[%d] = %t", index, intf.(bool))
			case map[string]interface{}:
				log.Debug(prefix + "[%d] is an object", index)
				func_parse_obj(intf.(map[string]interface{}), level + 1)
			case []interface{}:
				log.Debug(prefix + "[%d] is an array", index)
				func_parse_array(intf.([]interface{}), level + 1)
			}
			index = 0
			intf = nil
		}
		return
	}	// func_parse_array ends

	func_parse_obj = func(obj map[string]interface{}, level int) {
		prefix := strings.Repeat("    ", level)
		for key, intf := range obj {
			switch intf.(type) {
			case string:
				log.Debug(prefix + "\"%s\" = %s", key, intf.(string))
			case float64:
				integer := int64(intf.(float64))
				if (float64(integer) == intf.(float64)) {
					log.Debug(prefix + "\"%s\" = %d", key, integer)
				} else {
					log.Debug(prefix + "\"%s\" = %f", key, intf.(float64))
				}
			case bool:
				log.Debug(prefix + "\"%s\" = %t", key, intf.(bool))
			case map[string]interface{}:
				log.Debug(prefix + "\"%s\" is an object", key)
				func_parse_obj(intf.(map[string]interface{}), level + 1)
			case []interface{}:
				log.Debug(prefix + "\"%s\" is an array", key)
				func_parse_array(intf.([]interface{}), level + 1)
			default:
				// do nothing
			}
			key = ""
			intf = nil
		}
		return
	}	// func_parse_obj ends

	var data map[string]interface{}
	log.Info("Now try encoding/json without map parsing")
	start = time.Now()
	err = json.Unmarshal(json_bytes, &data)
	if err != nil {
		log.Error("unmarshal json error: " + err.Error())
	}
	mapIntf = elapsed(start)
	log.Info("End of encoding/json without map parsing: %v", mapIntf)

	log.Info("Now try encoding/json with map parsing")
	start = time.Now()
	err = json.Unmarshal(json_bytes, &data)
	if err != nil {
		log.Error("unmarshal json error: " + err.Error())
	} else {
		func_parse_obj(data, 0)
	}
	mapIntfAll = elapsed(start)
	log.Info("End of encoding/json with map parsing: %v", mapIntfAll)

	log.Info("Now try encoding/json struct")
	aStruct := teststruct{}
	start = time.Now()
	err = json.Unmarshal(json_bytes, &aStruct)
	if err != nil {
		log.Error("unmarshal json error: " + err.Error())
	}
	structure = elapsed(start)
	log.Info("End of encoding/json struct: %v", structure)

	var func_obj_each func([]byte, []byte, jsonparser.ValueType, int) error
	var func_array_each func([]byte, jsonparser.ValueType, int, error)
	prefix_level := 0

	func_obj_each = func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		prefix := strings.Repeat("    ", prefix_level)
		switch dataType {
		case jsonparser.String:
			log.Debug(prefix + "\"%s\" = %s", string(key), string(value))
		case jsonparser.Number:
			str := string(value)
			if index := strings.Index(string(value), "."); index >= 0 {
				flt_val, _ := strconv.ParseFloat(str, 64)
				log.Debug(prefix + "\"%s\" = %f", string(key), flt_val)
			} else {
				int_val, _ := strconv.ParseInt(str, 10, 64)
				log.Debug(prefix + "\"%s\" = %d", string(key), int_val)
			}
		case jsonparser.Object:
			log.Debug(prefix + "\"%s\" is an object", string(key))
			prefix_level ++
			jsonparser.ObjectEach(value, func_obj_each)
			prefix_level --
		case jsonparser.Array:
			log.Debug("\"%s\" is an array", string(key))
			prefix_level ++
			jsonparser.ArrayEach(value, func_array_each)
			prefix_level --
		case jsonparser.Boolean:
			log.Debug(prefix + "\"%s\" = %s", string(key), string(value))
		default:
			log.Debug(prefix + "\"%s\" type: %d", string(key), int(dataType))
		}
		return nil
	}

	func_array_each = func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		prefix := strings.Repeat("    ", prefix_level)
		switch dataType {
		case jsonparser.String:
			log.Debug(prefix + "[] = %s", string(value))
		case jsonparser.Number:
			str := string(value)
			if index := strings.Index(string(value), "."); index >= 0 {
				flt_val, _ := strconv.ParseFloat(str, 64)
				log.Debug(prefix + "[] = %f", flt_val)
			} else {
				int_val, _ := strconv.ParseInt(str, 10, 64)
				log.Debug(prefix + "[] = %d", int_val)
			}
		case jsonparser.Object:
			log.Debug(prefix + "[] is an object")
			prefix_level ++
			jsonparser.ObjectEach(value, func_obj_each)
			prefix_level --
		case jsonparser.Array:
			log.Debug(prefix + "[] is an array")
			prefix_level ++
			jsonparser.ArrayEach(value, func_array_each)
			prefix_level --
		case jsonparser.Boolean:
			log.Debug(prefix + "[] = %s", string(value))
		default:
			log.Debug(prefix + "[] type: %d", int(dataType))
		}
		return
	}

	log.Info("Now try jsonparser")
	start = time.Now()
	jsonparser.ObjectEach(json_bytes, func_obj_each)
	prefix_level = 0
	parser = elapsed(start)
	log.Info("End of jsonparser: %v", parser)

	log.Info("Now try jsonconv")
	start = time.Now()
	_, _ = jsonconv.NewFromString(string(json_bytes))
	conv = elapsed(start)
	log.Info("End of jsonconv: %v", conv)

	return
}
