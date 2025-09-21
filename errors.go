package main

import (
	"errors"
)

type ErrorCode int

const (
	ErrOk ErrorCode = iota
	ErrConfigFileNilOrEmpty
	ErrConfigFileNilOrEmptyButDefaultExists
	ErrConfigFileNotExists
	ErrCifreNotValid
	ErrBooleanNotValid
)

var ErrorMessages = map[ErrorCode]string{
	ErrOk:                   "No error",
	ErrConfigFileNilOrEmpty: "Config file argument is nil or empty",
	ErrConfigFileNilOrEmptyButDefaultExists: "Config file argument is nil or empty, but default config file exists",
	ErrConfigFileNotExists:  "Config file does not exist",
	ErrCifreNotValid:        "Cifre value is not a valid integer",
	ErrBooleanNotValid:      "Verbose value is not a valid boolean",
}

func (e ErrorCode) String() string {
	if msg, exists := ErrorMessages[e]; exists {
		return msg
	}
	return "Unknown error"
}

func (e ErrorCode) Error() error {
	if msg, exists := ErrorMessages[e]; exists {
		return errors.New(msg)
	}
	return errors.New("Unknown error")
}