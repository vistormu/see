package builder

import (
	"github.com/vistormu/go-dsa/errors"
)

const (
	NArgs               errors.ErrorType = "wrong number of arguments"
	FileNotFound        errors.ErrorType = "file not found"
	NotImplemented      errors.ErrorType = "not implemented"
	DirReading          errors.ErrorType = "directory reading error"
	WrongArg            errors.ErrorType = "wrong argument type"
	UnknownArg          errors.ErrorType = "unknown argument"
	ExpectedValue       errors.ErrorType = "expected value for argument"
	ParseError          errors.ErrorType = "argument parsing error"
	UnsupportedFileType errors.ErrorType = "unsupported file type"
)
