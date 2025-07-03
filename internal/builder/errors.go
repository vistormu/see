package builder

import (
	"github.com/vistormu/go-dsa/errors"
)

const (
	FileInfo       errors.ErrorType = "error retrieving file information"
	FileNotFound   errors.ErrorType = "file not found"
	ZoxideNotFound errors.ErrorType = "zoxide command not found"
)
