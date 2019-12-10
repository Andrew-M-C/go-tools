package logger

import (
    "testing"
    "time"
)

func TestStdout(t *testing.T) {
    l := NewConsole()
    l.Debugf("Hello, world on %v", time.Now().Local())
}
