package log

import (
	"fmt"
	"time"
	"runtime"
	"strings"
	"path/filepath"
)

var (
	shouldDebug	= false
	shouldInfo	= true
	shouldWarn	= true
	shouldError	= true
)

const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_NONE
)

func SetLogLevel(level int) {
	switch level {
	case LEVEL_DEBUG:
		shouldDebug = true
		shouldInfo = true
		shouldWarn = true
		shouldError = true
	case LEVEL_INFO:
		shouldDebug = false
		shouldInfo = true
		shouldWarn = true
		shouldError = true
	case LEVEL_WARN:
		shouldDebug = false
		shouldInfo = false
		shouldWarn = true
		shouldError = true
	case LEVEL_ERROR:
		shouldDebug = false
		shouldInfo = false
		shouldWarn = false
		shouldError = true
	default:
		shouldDebug = false
		shouldInfo = false
		shouldWarn = false
		shouldError = false
	}
	return
}

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
	if false == shouldDebug {
		return
	}
	datetime := getTimeStr()
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	fmt.Printf("%s - DEBUG - %s, Line %d, %s() - %s\n", datetime, file, line, function, text)
	return
}


func Info(format string, v ...interface{}) {
	if false == shouldInfo {
		return
	}
	datetime := getTimeStr()
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	fmt.Printf("%s - INFO  - %s, Line %d, %s() - %s\n", datetime, file, line, function, text)
	return
}


func Warn(format string, v ...interface{}) {
	if false == shouldWarn {
		return
	}
	datetime := getTimeStr()
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	fmt.Printf("%s - WARN  - %s, Line %d, %s() - %s\n", datetime, file, line, function, text)
	return
}


func Error(format string, v ...interface{}) {
	if false == shouldError {
		return
	}
	datetime := getTimeStr()
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	fmt.Printf("%s - ERROR - %s, Line %d, %s() - %s\n", datetime, file, line, function, text)
	return
}
