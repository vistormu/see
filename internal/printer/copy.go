package printer

import (
	"io"
	"os/exec"
	"runtime"

	"github.com/vistormu/go-dsa/errors"
)

var copyFn = copyToClipboard

func clipboardCommand() (string, []string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "pbcopy", []string{}, nil
	case "windows":
		return "clip", []string{}, nil
	default:
		if _, err := exec.LookPath("wl-copy"); err == nil {
			return "wl-copy", []string{}, nil
		}
		if _, err := exec.LookPath("xclip"); err == nil {
			return "xclip", []string{"-selection", "clipboard"}, nil
		}
		if _, err := exec.LookPath("xsel"); err == nil {
			return "xsel", []string{"--clipboard", "--input"}, nil
		}
	}

	return "", nil, errors.New(ClipboardCopy).With("reason", "no clipboard command found")
}

func copyToClipboard(content string) error {
	command, args, err := clipboardCommand()
	if err != nil {
		return err
	}

	cmd := exec.Command(command, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.New(ClipboardCopy).Wrap(err)
	}

	if err := cmd.Start(); err != nil {
		return errors.New(ClipboardCopy).Wrap(err)
	}

	if _, err := io.WriteString(stdin, content); err != nil {
		return errors.New(ClipboardCopy).Wrap(err)
	}

	if err := stdin.Close(); err != nil {
		return errors.New(ClipboardCopy).Wrap(err)
	}

	if err := cmd.Wait(); err != nil {
		return errors.New(ClipboardCopy).Wrap(err)
	}

	return nil
}
