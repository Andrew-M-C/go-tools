package jsonconv

import (
	"errors"
)

var (
	DataTypeError		= errors.New("invalid parameter type")
	ParaError			= errors.New("parameter invalid")
	JsonFormatError		= errors.New("json string format error")
	IndexOutOfBoundsError	= errors.New("index out of bounds")

	NotAnArrayError		= errors.New("target is not an array")
	NotAnObjectError	= errors.New("target is not an object")
	NotAStringError		= errors.New("target is not a string")
	NotANumberError		= errors.New("target is not a number")
	NotABoolError		= errors.New("target is not a bool")

	ObjectNotFoundError	= errors.New("object not found")
)
