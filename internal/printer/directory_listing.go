package printer

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"see/internal/builder"

	"github.com/epilande/go-devicons"
	"github.com/vistormu/go-dsa/ansi"
	"github.com/vistormu/go-dsa/errors"
)

var hiddenStyle = ansi.Dim + ansi.Rgb(150, 150, 150)

const (
	treeMiddle      = TreeBranch + TreeHLine + " "
	treeEnd         = TreeLeaf + TreeHLine + " "
	treeVerticalGap = TreeVLine + "  "
	treeEmptyGap    = "   "
	rootDirStyle    = ansi.Bold + ansi.Green
	subDirStyle     = ansi.Bold + ansi.Blue
	fileStyle       = ansi.Cyan2
	itemsStyle      = ansi.Yellow
	sizeStyle       = ansi.Magenta2
	permsStyle      = ansi.Cyan
)

type rowKind int

const (
	rowRoot rowKind = iota
	rowDir
	rowFile
)

type listingRow struct {
	kind   rowKind
	prefix string
	branch string
	depth  int
	dir    *builder.Directory
	file   *builder.File
}

type renderedRow struct {
	source listingRow
	nFiles int
	nDirs  int
	hasDir bool
	size   string
	perms  string
}

func printDirectoryListing(root *builder.Directory, args builder.Args) error {
	sorters := map[string]func(*builder.Directory){
		"name":       sortByName,
		"kind":       sortByKind,
		"size":       sortBySize,
		"git-status": sortByGitStatus,
	}
	sorter, ok := sorters[args.Sort]
	if !ok {
		return errors.New(InvalidArgs).With("sort", args.Sort)
	}

	sorter(root)

	output, copyOutput := renderDirectoryListing(root, args)
	fmt.Printf("%s\n\n", output)

	if args.Copy {
		if err := copyFn(copyOutput); err != nil {
			return err
		}
	}

	return nil
}

func renderDirectoryListing(root *builder.Directory, args builder.Args) (string, string) {
	rows := flattenRows(root, args.Depth)
	if len(rows) == 0 {
		return "", ""
	}

	rendered := make([]renderedRow, len(rows))
	for i, row := range rows {
		rendered[i] = renderRow(row)
	}

	visible := make([]renderedRow, 0, len(rendered))
	for _, row := range rendered {
		if row.source.kind == rowRoot {
			visible = append(visible, row)
			continue
		}

		plain := strings.Join(
			[]string{
				stripAnsi(renderTreeCell(row.source, -1)),
				strconv.Itoa(row.nFiles),
				strconv.Itoa(row.nDirs),
				stripAnsi(row.perms),
				stripAnsi(row.size),
			},
			" ",
		)
		if args.Filter != "" && !strings.Contains(plain, args.Filter) {
			continue
		}
		visible = append(visible, row)
	}

	treeWidth := len("tree")
	for _, row := range visible {
		treeWidth = max(treeWidth, visibleWidth(renderTreeCell(row.source, -1)))
	}
	treeWidth = max(treeWidth, len("(empty)"))

	fileDigits := 1
	dirDigits := 1
	for _, row := range visible {
		if !row.hasDir {
			continue
		}
		if row.nFiles > 0 {
			fileDigits = max(fileDigits, len(strconv.Itoa(row.nFiles)))
		}
		if row.nDirs > 0 {
			dirDigits = max(dirDigits, len(strconv.Itoa(row.nDirs)))
		}
	}

	itemsWidth := len("items")
	sizeWidth := len("size")
	permsWidth := len("perms")
	for _, row := range visible {
		itemsWidth = max(itemsWidth, visibleWidth(itemsToString(row, fileDigits, dirDigits)))
		sizeWidth = max(sizeWidth, visibleWidth(row.size))
		permsWidth = max(permsWidth, visibleWidth(row.perms))
	}

	nCols := 4
	metaWidth := itemsWidth + permsWidth + sizeWidth
	borderWidth := nCols*3 + 1

	maxWidth := termWidth
	if maxWidth <= 0 {
		maxWidth = treeWidth + metaWidth + borderWidth
	}
	maxTreeWidth := maxWidth - metaWidth - borderWidth
	if maxTreeWidth < 1 {
		maxTreeWidth = 1
	}
	treeWidth = min(treeWidth, maxTreeWidth)

	colWidths := []int{treeWidth, itemsWidth, permsWidth, sizeWidth}

	topBorder := buildTableBorder("┌", "┬", "┐", colWidths)
	midBorder := buildTableBorder("├", "┼", "┤", colWidths)
	bottomBorder := buildTableBorder("└", "┴", "┘", colWidths)

	header := formatTableRow(
		ansi.Underline+"tree"+ansi.Reset,
		ansi.Underline+"items"+ansi.Reset,
		ansi.Underline+"perms"+ansi.Reset,
		ansi.Underline+"size"+ansi.Reset,
		treeWidth,
		itemsWidth,
		permsWidth,
		sizeWidth,
	)

	lines := []string{topBorder, header, midBorder}

	for _, row := range visible {
		lines = append(lines, formatTableRow(
			renderTreeCell(row.source, treeWidth),
			itemsToString(row, fileDigits, dirDigits),
			row.perms,
			row.size,
			treeWidth,
			itemsWidth,
			permsWidth,
			sizeWidth,
		))
	}
	lines = append(lines, bottomBorder)

	output := strings.Join(lines, "\n")
	return output, stripAnsi(output)
}

func buildTableBorder(left, middle, right string, widths []int) string {
	var b strings.Builder
	b.WriteString(left)
	for i, width := range widths {
		b.WriteString(repeat("─", width+2))
		if i < len(widths)-1 {
			b.WriteString(middle)
		}
	}
	b.WriteString(right)
	return b.String()
}

func formatTableRow(tree, items, perms, size string, treeWidth, itemsWidth, permsWidth, sizeWidth int) string {
	tree = padRight(truncateWithEllipsis(tree, treeWidth), treeWidth)
	items = padLeft(truncateWithEllipsis(items, itemsWidth), itemsWidth)
	perms = padRight(truncateWithEllipsis(perms, permsWidth), permsWidth)
	size = padLeft(truncateWithEllipsis(size, sizeWidth), sizeWidth)

	return "│ " + tree + " │ " + items + " │ " + perms + " │ " + size + " │"
}

func renderRow(row listingRow) renderedRow {
	switch row.kind {
	case rowRoot:
		return renderedRow{
			source: row,
			nFiles: len(row.dir.Files),
			nDirs:  len(row.dir.Dirs),
			hasDir: true,
			size:   sizeToString(row.dir.Size),
			perms:  permsToString(row.dir.Mode),
		}
	case rowDir:
		return renderedRow{
			source: row,
			nFiles: len(row.dir.Files),
			nDirs:  len(row.dir.Dirs),
			hasDir: true,
			size:   sizeToString(row.dir.Size),
			perms:  permsToString(row.dir.Mode),
		}
	default:
		return renderedRow{
			source: row,
			size:   sizeToString(row.file.Size),
			perms:  permsToString(row.file.Mode),
		}
	}
}

func itemsToString(row renderedRow, fileDigits, dirDigits int) string {
	if !row.hasDir {
		return ""
	}

	fileSlotWidth := fileDigits + 1
	dirSlotWidth := dirDigits + 1

	filePart := repeat(" ", fileSlotWidth)
	if row.nFiles > 0 {
		filePart = itemsStyle + fmt.Sprintf("%*d", fileDigits, row.nFiles) + FileIcon + ansi.Reset
	}

	dirPart := repeat(" ", dirSlotWidth)
	if row.nDirs > 0 {
		dirPart = itemsStyle + fmt.Sprintf("%*d", dirDigits, row.nDirs) + DirIcon + ansi.Reset
	}

	return strings.TrimRight(filePart+" "+dirPart, " ")
}

func renderTreeCell(row listingRow, maxWidth int) string {
	switch row.kind {
	case rowRoot:
		gitStatus, _ := gitStatusToString(row.dir.GitStatus)
		name, _ := dirNameToString(row.dir.Name, true)
		return fitTreeCell(gitStatus+name, maxWidth, func(limit int) string {
			truncated := truncateStyledName(row.dir.Name, true, limit)
			return gitStatus + truncated
		})
	case rowDir:
		gitStatus, _ := gitStatusToString(row.dir.GitStatus)
		basePrefix := row.prefix + row.branch + gitStatus
		nameCell := fitNameWithPrefix(basePrefix, row.dir.Name, true, false, maxWidth)
		candidate := basePrefix + nameCell
		return fitTreeCell(candidate, maxWidth, func(limit int) string {
			return truncateWithEllipsis(stripAnsi(candidate), limit)
		})
	default:
		icon, _ := iconToString(row.file.Name)
		basePrefix := row.prefix + row.branch + icon + " "
		nameCell := fitNameWithPrefix(basePrefix, row.file.Name, false, strings.HasPrefix(row.file.Name, "."), maxWidth)
		candidate := basePrefix + nameCell
		return fitTreeCell(candidate, maxWidth, func(limit int) string {
			return truncateWithEllipsis(stripAnsi(candidate), limit)
		})
	}
}

func fitTreeCell(content string, maxWidth int, fallback func(int) string) string {
	if maxWidth < 0 || visibleWidth(content) <= maxWidth {
		return content
	}

	next := fallback(maxWidth)
	if visibleWidth(next) <= maxWidth {
		return next
	}

	return truncateWithEllipsis(stripAnsi(next), maxWidth)
}

func fitNameWithPrefix(prefix, name string, isDir, hidden bool, maxWidth int) string {
	if maxWidth < 0 {
		if isDir {
			result, _ := dirNameToString(name, false)
			return result
		}
		result, _ := fileNameToString(name)
		return result
	}

	available := maxWidth - visibleWidth(prefix)
	if available <= 1 {
		return truncateWithEllipsis("", 1)
	}

	if isDir {
		return truncateStyledName(name, false, available)
	}
	return truncateFileName(name, hidden, available)
}

func truncateStyledName(name string, isRoot bool, width int) string {
	plainName := name + "/"
	truncated := truncateWithEllipsis(plainName, width)

	style := subDirStyle
	if isRoot {
		style = rootDirStyle
	}
	if strings.HasPrefix(name, ".") {
		style = hiddenStyle
	}

	return style + truncated + ansi.Reset
}

func truncateFileName(name string, hidden bool, width int) string {
	truncated := truncateWithEllipsis(name, width)
	style := fileStyle
	if hidden {
		style = hiddenStyle
	}
	return style + truncated + ansi.Reset
}

func flattenRows(root *builder.Directory, maxDepth int) []listingRow {
	rows := make([]listingRow, 0, 1+len(root.Dirs)+len(root.Files))
	rows = append(rows, listingRow{
		kind:  rowRoot,
		depth: 0,
		dir:   root,
	})

	appendChildren(root, "", &rows, 1, maxDepth)
	return rows
}

func appendChildren(parent *builder.Directory, prefix string, rows *[]listingRow, currentDepth int, maxDepth int) {
	if maxDepth > 0 && currentDepth > maxDepth {
		return
	}

	totalChildren := len(parent.Dirs) + len(parent.Files)
	if totalChildren == 0 {
		return
	}

	for i, dir := range parent.Dirs {
		isLast := i == totalChildren-1
		branch := treeMiddle
		if isLast {
			branch = treeEnd
		}
		*rows = append(*rows, listingRow{
			kind:   rowDir,
			prefix: prefix,
			branch: branch,
			depth:  currentDepth,
			dir:    dir,
		})

		childPrefix := prefix + treeVerticalGap
		if isLast {
			childPrefix = prefix + treeEmptyGap
		}
		appendChildren(dir, childPrefix, rows, currentDepth+1, maxDepth)
	}

	for i, file := range parent.Files {
		j := len(parent.Dirs) + i
		isLast := j == totalChildren-1
		branch := treeMiddle
		if isLast {
			branch = treeEnd
		}
		current := file
		*rows = append(*rows, listingRow{
			kind:   rowFile,
			prefix: prefix,
			branch: branch,
			depth:  currentDepth,
			file:   &current,
		})
	}
}

func gitStatusToString(status builder.GitStatus) (string, int) {
	var gitColor string
	switch status {
	case builder.GitUpToDate:
		gitColor = ansi.Green
	case builder.GitUntracked:
		gitColor = ansi.Yellow
	case builder.GitModified:
		gitColor = ansi.Red
	case builder.GitUnknown:
		return "", 0
	}

	str := fmt.Sprintf("%s%s%s ", gitColor, GitIcon, ansi.Reset)
	return str, 2
}

func dirNameToString(name string, isRoot bool) (string, int) {
	dirStyle := subDirStyle
	switch {
	case strings.HasPrefix(name, "."):
		dirStyle = hiddenStyle
	case isRoot:
		dirStyle = rootDirStyle
	}

	nameStr := fmt.Sprintf("%s%s/%s", dirStyle, name, ansi.Reset)
	return nameStr, len(name) + 1
}

func sizeToString(size int64) string {
	sizeStr := humanizeSize(size)
	if strings.HasSuffix(sizeStr, " B") {
		sizeStr += " "
	}
	return fmt.Sprintf("%s%s%s", sizeStyle, sizeStr, ansi.Reset)
}

func permsToString(mode fmt.Stringer) string {
	return fmt.Sprintf("%s%s%s", permsStyle, mode.String(), ansi.Reset)
}

func iconToString(name string) (string, int) {
	icon := devicons.IconForPath(name)
	iconStr := fmt.Sprintf("%s%s%s", ansi.Hex(icon.Color), icon.Icon, ansi.Reset)
	return iconStr, 1
}

func fileNameToString(name string) (string, int) {
	style := fileStyle
	if strings.HasPrefix(name, ".") {
		style = hiddenStyle
	}

	nameStr := fmt.Sprintf("%s%s%s", style, name, ansi.Reset)
	return nameStr, len(name)
}

func sortByName(root *builder.Directory) {
	slices.SortFunc(root.Dirs, func(a, b *builder.Directory) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
	slices.SortFunc(root.Files, func(a, b builder.File) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	for _, dir := range root.Dirs {
		sortByName(dir)
	}
}

func sortByKind(root *builder.Directory) {
	slices.SortFunc(root.Dirs, func(a, b *builder.Directory) int {
		return strings.Compare(a.Name, b.Name)
	})
	slices.SortFunc(root.Files, func(a, b builder.File) int {
		return strings.Compare(a.Name, b.Name)
	})

	for _, dir := range root.Dirs {
		sortByKind(dir)
	}
}

func sortBySize(root *builder.Directory) {
	slices.SortFunc(root.Dirs, func(a, b *builder.Directory) int {
		if a.Size < b.Size {
			return -1
		}
		if a.Size > b.Size {
			return 1
		}
		return 0
	})
	slices.SortFunc(root.Files, func(a, b builder.File) int {
		if a.Size < b.Size {
			return -1
		}
		if a.Size > b.Size {
			return 1
		}
		return 0
	})

	for _, dir := range root.Dirs {
		sortBySize(dir)
	}
}

func sortByGitStatus(root *builder.Directory) {
	gitStatusOrder := map[builder.GitStatus]int{
		builder.GitUnknown:   0,
		builder.GitUpToDate:  1,
		builder.GitUntracked: 2,
		builder.GitModified:  3,
	}

	slices.SortFunc(root.Dirs, func(a, b *builder.Directory) int {
		return gitStatusOrder[a.GitStatus] - gitStatusOrder[b.GitStatus]
	})

	for _, dir := range root.Dirs {
		sortByGitStatus(dir)
	}
}
