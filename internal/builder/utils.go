package builder

import (
	"os"
	"strconv"
	"strings"

	"github.com/vistormu/go-dsa/errors"
)

func validFlag(args []string, next int) bool {
	return next < len(args) && !strings.HasPrefix(args[next], "--") && args[next] != ""
}

func parseArgs(args []string) (map[string]any, error) {
	parsedArgs := map[string]any{
		"dir":       "",
		"file":      "",
		"--depth":   1,
		"--sort-by": "name",
		"--head":    -1,
		"--tail":    -1,
	}

	if len(args) == 0 {
		return parsedArgs, nil
	}

	skip := false
	for i, arg := range args {
		if skip {
			skip = false
			continue
		}

		// flag case
		if strings.HasPrefix(arg, "--") {
			typ, ok := parsedArgs[arg]
			if !ok {
				return nil, errors.New(UnknownArg).With("arg", arg)
			}

			switch typ.(type) {
			case string:
				if !validFlag(args, i+1) {
					return nil, errors.New(ExpectedValue).With("arg", arg)
				}
				parsedArgs[arg] = args[i+1]
				skip = true
				continue

			case int:
				if !validFlag(args, i+1) {
					return nil, errors.New(ExpectedValue).With("arg", arg)
				}
				value := args[i+1]
				intValue, err := strconv.Atoi(value)
				if err != nil {
					return nil, errors.New(ParseError).With("arg", arg).With("value", value).Wrap(err)
				}
				parsedArgs[arg] = intValue
				skip = true
				continue

			case bool:
				parsedArgs[arg] = true
				continue
			}
		}

		// positional argument case
		fileInfo, err := os.Stat(arg)
		if err != nil {
			return nil, errors.New(FileNotFound).With("file", arg).Wrap(err)
		}

		if fileInfo.IsDir() {
			if parsedArgs["dir"] != "" {
				return nil, errors.New(NArgs).With("expected", "1 directory").With("got", len(args))
			}
			parsedArgs["dir"] = arg
		} else {
			if parsedArgs["file"] != "" {
				return nil, errors.New(NArgs).With("expected", "1 file").With("got", len(args))
			}

			if !fileInfo.Mode().IsRegular() {
				return nil, errors.New(UnsupportedFileType).With("file", arg)
			}

			parsedArgs["file"] = arg
		}

	}

	return parsedArgs, nil
}
