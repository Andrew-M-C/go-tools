package logger

import (
    "github.com/Andrew-M-C/go-tools/callstack"
    "github.com/petermattis/goid"
    "strconv"
    "bytes"
    "time"
    "fmt"
)

type Console struct {}

func NewConsole() *Console {
    return &Console{}
}

func packPrefix(level string) string {
    fi, li, fu := callstack.CallerInfo(2)
    gid := goidToStr(goid.Get())
    t := time.Now().Local().Format("2006-01-02 15:04:05.000000")

    buff := bytes.Buffer{}
    buff.WriteString(t + " | ")
    buff.WriteString(gid + " | ")
    buff.WriteString(level + " | ")
    buff.WriteString(fi + ", ")
    buff.WriteString("Line " + strconv.Itoa(li) + ", ")
    buff.WriteString(fu + "() | ")
    return buff.String()
}

func goidToStr(id int64) string {
    s := strconv.Itoa(int(id))
    switch len(s) {
    default:
        return s
    case 1:
        return "0000" + s
    case 2:
        return "000" + s
    case 3:
        return "00" + s
    case 4:
        return "0" + s
    }
}

func (Console) Debug(v ...interface{}) {
    param := []interface{}{packPrefix("DEBUG")}
    param = append(param, v...)
    fmt.Println(param...)
}

func (Console) Debugf(format string, a ...interface{}) {
    fmt.Printf(packPrefix("DEBUG") + format + "\n", a...)
}

func (Console) Info(v ...interface{}) {
    param := []interface{}{packPrefix("INFO ")}
    param = append(param, v...)
    fmt.Println(param...)
}

func (Console) Infof(format string, a ...interface{}) {
    fmt.Printf(packPrefix("INFO ") + format + "\n", a...)
}

func (Console) Warn(v ...interface{}) {
    param := []interface{}{packPrefix("WARN ")}
    param = append(param, v...)
    fmt.Println(param...)
}

func (Console) Warnf(format string, a ...interface{}) {
    fmt.Printf(packPrefix("WARN ") + format + "\n", a...)
}

func (Console) Error(v ...interface{}) {
    param := []interface{}{packPrefix("ERROR")}
    param = append(param, v...)
    fmt.Println(param...)
}

func (Console) Errorf(format string, a ...interface{}) {
    fmt.Printf(packPrefix("ERROR") + format + "\n", a...)
}
