package jsonconv

import (
	"errors"
)

var (
	DataTypeError		= errors.New("invalid parameter type")
	JsonFormatError		= errors.New("json string format error")
	IndexOutOfBoundsError	= errors.New("index out of bounds")

	NotAnArrayError		= errors.New("target is not an array")
	NotAnObjectError	= errors.New("target is not an object")

	ObjectNotFoundError	= errors.New("object not found")
)
