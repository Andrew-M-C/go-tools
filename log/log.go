package log

import (
	"fmt"
	"time"
)

func Debug(format string, v ...interface{}) {
	now := time.Now()
	datetime := now.Local().Format("2006-01-02 15:04:05.000000")
	fmt.Printf(datetime + " - DEBUG - " + format + "\n", v...)
	return
}
