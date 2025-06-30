package builder

import (
	"github.com/vistormu/go-dsa/errors"
)

func BuildCommand(args []string) (any, error) {
	// parse args
	parsedArgs, err := parseArgs(args)
	if err != nil {
		return nil, err
	}

	// default case
	if parsedArgs["dir"] == "" && parsedArgs["file"] == "" {
		parsedArgs["dir"] = "."
	}

	if parsedArgs["dir"] != "" {
		return buildDirectoryListing(parsedArgs)
	}

	if parsedArgs["file"] != "" {
		return buildFileContent(parsedArgs)
	}

	return nil, errors.New(NotImplemented)
}
