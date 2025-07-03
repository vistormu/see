package builder

import (
	"os"
	"os/exec"
	"strings"

	"github.com/vistormu/go-dsa/errors"
)

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

func BuildCommand(args Args) (any, error) {
	info, err := os.Stat(args.Element)
	if os.IsNotExist(err) {
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

	if info.IsDir() {
		return buildDirectoryListing(args)
	} else {
		return buildFileContent(args)
	}
}
