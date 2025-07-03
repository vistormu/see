package printer

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"see/internal/builder"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/vistormu/go-dsa/ansi"
)

func highlight(path, content, styleName string) (string, error) {
	lexer := lexers.Match(filepath.Base(path))
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(styleName)
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.TTY16m

	iterator, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func addLineNumbers(src string) string {
	if src == "" {
		return src
	}

	lines := strings.Split(src, "\n")
	width := len(strconv.Itoa(len(lines)))

	var buf bytes.Buffer
	buf.Grow(len(src) + len(lines)*width + len(lines)*3) // simple capacity guess

	format := "%*" + "d │ " // e.g. "%3d │ "
	for i, line := range lines {
		buf.WriteString(fmt.Sprintf(format, width, i+1))
		buf.WriteString(line)

		if i < len(lines)-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

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

	halfFreeSpace := termWidth/2 - nameLength - nLinesLength - nSizeLength - 1

	freeSpace := halfFreeSpace

	spaces := repeat(" ", freeSpace)

	// final touches
	content := strings.ReplaceAll(fileContent.Content, "\t", "    ")
	hlContent, err := highlight(fileContent.File.Path, content, "catppuccin-mocha")
	if err == nil {
		content = hlContent
	}
	content = addLineNumbers(content)

	fmt.Printf("%s%s%s%s\n\n%s\n\n",
		name,
		spaces,
		nLines,
		size,
		content,
	)

	return nil
}
