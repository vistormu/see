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
	xansi "github.com/charmbracelet/x/ansi"
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

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func addLineNumbers(src string, start int) string {
	if src == "" {
		return src
	}

	lines := strings.Split(src, "\n")
	width := len(strconv.Itoa(start + len(lines) - 1))

	var buf bytes.Buffer
	buf.Grow(len(src) + len(lines)*(width+4))

	format := hiddenStyle + "%*" + "d │ " + ansi.Reset
	for i, line := range lines {
		buf.WriteString(fmt.Sprintf(format, width, start+i))
		buf.WriteString(line)

		if i < len(lines)-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func clampLineOption(value int) int {
	if value < 0 {
		return -1
	}
	return max(0, value)
}

func selectContentWindow(content string, head, tail int) (string, int) {
	lines := strings.Split(content, "\n")
	if len(lines) == 1 && lines[0] == "" {
		return "", 1
	}

	head = clampLineOption(head)
	tail = clampLineOption(tail)

	if head >= 0 {
		if head < len(lines) {
			lines = lines[:head]
		}
		return strings.Join(lines, "\n"), 1
	}

	if tail >= 0 {
		start := max(0, len(lines)-tail)
		return strings.Join(lines[start:], "\n"), start + 1
	}

	return content, 1
}

func buildFileHeader(fileContent *builder.FileContent, visibleLines int, lineWord string) string {
	name := fmt.Sprintf("%s%s%s%s", ansi.Bold, ansi.Green, fileContent.File.Name, ansi.Reset)
	size := fmt.Sprintf("%s%s%s", ansi.Magenta2, humanizeSize(fileContent.File.Size), ansi.Reset)
	lines := fmt.Sprintf("%s%d %s%s", ansi.Yellow2, visibleLines, lineWord, ansi.Reset)

	return name + "  " + lines + "  " + size
}

func framedContent(header, content string) string {
	lines := strings.Split(content, "\n")

	maxVisible := visibleWidth(header)
	for _, line := range lines {
		maxVisible = max(maxVisible, visibleWidth(line))
	}

	maxInner := termWidth - 4
	if maxInner < 24 {
		maxInner = 24
	}
	innerWidth := min(maxVisible, maxInner)
	if innerWidth < 24 {
		innerWidth = 24
	}

	wrappedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		wrapped := xansi.Wrap(line, innerWidth, "")
		parts := strings.Split(wrapped, "\n")
		wrappedLines = append(wrappedLines, parts...)
	}
	if len(wrappedLines) == 0 {
		wrappedLines = []string{""}
	}

	top := "┌" + repeat("─", innerWidth+2) + "┐"
	mid := "├" + repeat("─", innerWidth+2) + "┤"
	bottom := "└" + repeat("─", innerWidth+2) + "┘"

	out := []string{
		top,
		"│ " + padRight(truncateWithEllipsis(header, innerWidth), innerWidth) + " │",
		mid,
	}
	for _, line := range wrappedLines {
		out = append(out, "│ "+padRight(line, innerWidth)+" │")
	}
	out = append(out, bottom)

	return strings.Join(out, "\n")
}

func printFileContent(fileContent *builder.FileContent, args builder.Args) error {
	selected, startLine := selectContentWindow(fileContent.Content, args.Head, args.Tail)
	if args.Filter != "" {
		selected = filterLines(selected, args.Filter)
		selected = strings.TrimRight(selected, "\n")
	}

	visibleLines := 0
	if selected != "" {
		visibleLines = strings.Count(selected, "\n") + 1
	}
	lineWord := "lines"
	if visibleLines == 1 {
		lineWord = "line"
	}

	header := buildFileHeader(fileContent, visibleLines, lineWord)

	copyContent := selected

	content := selected
	if content != "" {
		hlContent, err := highlight(fileContent.File.Path, content, "catppuccin-mocha")
		if err == nil {
			content = hlContent
		}
		content = addLineNumbers(content, startLine)
	}

	if content == "" {
		content = fmt.Sprintf("%s(empty)%s", hiddenStyle, ansi.Reset)
	}

	fmt.Printf("%s\n\n", framedContent(header, content))

	if args.Copy {
		if err := copyFn(copyContent); err != nil {
			return err
		}
	}

	return nil
}
