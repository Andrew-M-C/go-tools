package log

import (
	"fmt"
	"time"
	"runtime"
	"strings"
	"path/filepath"
	"strconv"
)

func getCallerInfo(invoke_level int) (fileName string, line int, funcName string) {
	funcName = "FILE"
	line = -1
	fileName = "FUNC"

	if invoke_level <= 0 {
		invoke_level = 2
	} else {
		invoke_level += 1
	}

	pc, file_name, line, ok := runtime.Caller(invoke_level)
	if ok {
		fileName = filepath.Base(file_name)
		func_name := runtime.FuncForPC(pc).Name()
		func_name = filepath.Ext(func_name)
		funcName = strings.TrimPrefix(func_name, ".")
	}
	// fmt.Println(reflect.TypeOf(pc), reflect.ValueOf(pc))
	return
}

func getTimeStr() string {
	return time.Now().Local().Format("2006-01-02 15:04:05.000000")
}

func Debug(format string, v ...interface{}) {
	datetime := getTimeStr()
	file, line, function := getCallerInfo(0);
	fmt.Printf(datetime + " - DEBUG - " + file + ", " + strconv.Itoa(line) + ", " + function + "() - " + format + "\n", v...)
	return
}


func Info(format string, v ...interface{}) {
	datetime := getTimeStr()
	file, line, function := getCallerInfo(0);
	fmt.Printf(datetime + " - INFO  - " + file + ", " + strconv.Itoa(line) + ", " + function + "() - " + format + "\n", v...)
	return
}


func Warn(format string, v ...interface{}) {
	datetime := getTimeStr()
	file, line, function := getCallerInfo(0);
	fmt.Printf(datetime + " - WARN  - " + file + ", " + strconv.Itoa(line) + ", " + function + "() - " + format + "\n", v...)
	return
}


func Error(format string, v ...interface{}) {
	datetime := getTimeStr()
	file, line, function := getCallerInfo(0);
	fmt.Printf(datetime + " - ERROR - " + file + ", " + strconv.Itoa(line) + ", " + function + "() - " + format + "\n", v...)
	return
}
