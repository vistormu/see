package builder

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/vistormu/go-dsa/errors"
)

var envVarPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func findWithZoxide(pattern string) (string, error) {
	cmd := exec.Command("zoxide", "query", "--list", pattern)
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New(ZoxideNotFound).With("pattern", pattern).Wrap(err)
	}

	dir := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	if dir == "" {
		return "", errors.New(ZoxideNotFound).With("pattern", pattern)
	}

	return dir, nil
}

func findWithEnvVar(input string) (*EnvVariable, bool) {
	key := strings.TrimPrefix(input, "$")
	if !envVarPattern.MatchString(key) {
		return nil, false
	}

	value, ok := os.LookupEnv(key)
	if !ok {
		return nil, false
	}

	return &EnvVariable{
		Name:  key,
		Value: value,
	}, true
}

func findWithEnvValue(value string) (*EnvVariable, bool) {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if parts[1] != value {
			continue
		}

		return &EnvVariable{
			Name:  parts[0],
			Value: parts[1],
		}, true
	}

	return nil, false
}

func BuildCommand(args Args) (any, error) {
	if args.Depth < 0 {
		return nil, errors.New(InvalidArgs).With("depth", args.Depth)
	}

	info, err := os.Stat(args.Element)
	if os.IsNotExist(err) {
		if envVar, ok := findWithEnvVar(args.Element); ok {
			return envVar, nil
		}
		if envVar, ok := findWithEnvValue(args.Element); ok {
			return envVar, nil
		}

		_, err := exec.LookPath("zoxide")
		if err != nil {
			return nil, errors.New(FileNotFound).With("path", args.Element).Wrap(err)
		}

		foundDir, err := findWithZoxide(args.Element)
		if err != nil {
			return nil, errors.New(FileNotFound).With("path", args.Element).Wrap(err)
		}

		args.Element = foundDir
		info, err = os.Stat(args.Element)
		if err != nil {
			return nil, errors.New(FileInfo).With("path", args.Element).Wrap(err)
		}
	}
	if err != nil {
		return nil, errors.New(FileInfo).With("path", args.Element).Wrap(err)
	}

	mode := info.Mode()
	if !mode.IsRegular() && !mode.IsDir() {
		if envVar, ok := findWithEnvValue(args.Element); ok {
			return envVar, nil
		}
		return nil, errors.New(FileInfo).With("path", args.Element)
	}

	if info.IsDir() {
		return buildDirectoryListing(args)
	} else {
		return buildFileContent(args)
	}
}
