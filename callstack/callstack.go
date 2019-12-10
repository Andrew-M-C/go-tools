package callstack

import(
	"strings"
	"runtime"
	"path/filepath"
)

func CallerInfo(invokeLevel int) (fileName string, line int, funcName string) {
	funcName = "FILE"
	line = -1
	fileName = "FUNC"

	if invokeLevel <= 0 {
		invokeLevel = 2
	} else {
		invokeLevel += 1
	}

	pc, fi, line, ok := runtime.Caller(invokeLevel)
	if ok {
		fileName = filepath.Base(fi)
		fu := runtime.FuncForPC(pc).Name()
		fu = filepath.Ext(fu)
		funcName = strings.TrimPrefix(fu, ".")
	}
	return
}
