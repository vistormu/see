package printer

import (
	"golang.org/x/term"
	"os"

	"see/internal/builder"
)

var termWidth int

func init() {
	var err error
	termWidth, _, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth = 80 // default width if we can't get the terminal size
	}
}

func Print(command any, args builder.Args) error {
	switch cmd := command.(type) {
	case *builder.Directory:
		return printDirectoryListing(cmd, args)

	case *builder.FileContent:
		return printFileContent(cmd, args)
	}

	return nil
}
