package jsonconv

import (
	"errors"
	"fmt"
	"strings"
	"bytes"
	"time"
)

var (
	DataTypeError		= errors.New("invalid parameter type")
	ParaError			= errors.New("parameter invalid")
	JsonFormatError		= errors.New("json string format error")
	JsonTypeError		= errors.New("json target type error")
	IndexOutOfBoundsError	= errors.New("index out of bounds")

	NotAnArrayError		= errors.New("target is not an array")
	NotAnObjectError	= errors.New("target is not an object")
	NotAStringError		= errors.New("target is not a string")
	NotANumberError		= errors.New("target is not a number")
	NotABoolError		= errors.New("target is not a bool")

	ObjectNotFoundError	= errors.New("object not found")
)


type Filter int
const (
	Normal		Filter = 0
	IncludeMode	Filter = 1
	ExcludeMode	Filter = 2
)

type Option struct {
	ShowNull		bool
	EnsureAscii		bool
	FloatDigits		uint8
	TimeDigits		uint8
	FilterMode		Filter
	FilterList		[]string
}

var dftOption = Option {
	ShowNull: false,
	EnsureAscii: false,
	FloatDigits: 0,
	TimeDigits: 0,
	FilterMode: Normal,
}

func escapeJsonString(s string, ensureAscii bool) string {
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
			if ensureAscii && chr > '\u0127' {
				b.WriteString(fmt.Sprintf("\\u%04x", chr))
			} else {
				b.WriteRune(chr)
			}
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