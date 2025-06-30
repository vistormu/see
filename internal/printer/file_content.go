package printer

import (
	"fmt"

	"see/internal/builder"

	"github.com/vistormu/go-dsa/ansi"
)

func printFileContent(fileContent *builder.FileContent) error {
	name := fmt.Sprintf("%s%s%s%s",
		ansi.Bold,
		ansi.Green,
		fileContent.File.Name,
		ansi.Reset,
	)
	nameLength := len(fileContent.File.Name)

	nLines := fmt.Sprintf("%s%d lines %s",
		ansi.Yellow2,
		fileContent.NLines,
		ansi.Reset,
	)
	nLinesLength := len(fmt.Sprintf("%d lines ", fileContent.NLines))

	size := fmt.Sprintf("%s%s%s",
		ansi.Magenta2,
		humanizeSize(fileContent.File.Size),
		ansi.Reset,
	)
	nSizeLength := len(humanizeSize(fileContent.File.Size))

	if fileContent.File.Size == 0 {
		size = fmt.Sprintf("%s%s%s",
			ansi.Red2,
			"empty",
			ansi.Reset,
		)
		nSizeLength = len("empty")
		nLines = ""
		nLinesLength = 0
	}

	spaces := repeat(" ", termWidth-nameLength-nLinesLength-nSizeLength-1)

	fmt.Printf("%s%s%s%s\n\n%s\n",
		name,
		spaces,
		nLines,
		size,
		fileContent.Content,
	)

	return nil
}
