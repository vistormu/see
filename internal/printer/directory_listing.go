package printer

import (
	"fmt"
	"strings"

	"see/internal/builder"

	"github.com/epilande/go-devicons"
	"github.com/vistormu/go-dsa/ansi"
)

// ====
// init
// ====
var hiddenStyle = ansi.Dim + ansi.Rgb(128, 128, 128)

const (
	treeMiddle   = TreeBranch + TreeHLine + " "
	treeEnd      = TreeLeaf + TreeHLine + " "
	rootDirStyle = ansi.Bold + ansi.Green
	subDirStyle  = ansi.Bold + ansi.Blue
	fileStyle    = ansi.Cyan2
	itemsStyle   = ansi.Yellow
	sizeStyle    = ansi.Magenta2
	emptyStyle   = ansi.Red2
)

// =======
// columns
// =======
// tree
type rootDirName struct {
	git    string
	name   string
	length int
}

type subDirName struct {
	divider string
	git     string
	name    string
	length  int
}

type fileName struct {
	divider string
	icon    string
	name    string
	length  int
}

type treeColumn struct {
	root      rootDirName
	dirs      []subDirName
	files     []fileName
	maxLength int
}

// files
type nFiles struct {
	files  string
	length int
}

type filesColumn struct {
	root      nFiles
	items     []nFiles
	maxLength int
}

// dirs
type nDirs struct {
	dirs   string
	length int
}

type dirsColumn struct {
	root      nDirs
	items     []nDirs
	maxLength int
}

// size
type size struct {
	size   string
	length int
}

type sizeColumn struct {
	root      size
	items     []size
	maxLength int
}

// =====
// print
// =====
func printDirectoryListing(root *builder.Directory) error {
	treeColumn := createTreeColumn(root)
	filesColumn := createFilesColumn(root)
	dirsColumn := createDirsColumn(root)
	sizeColumn := createSizeColumn(root)

	// fullFreeSpace := termWidth - treeColumn.maxLength - filesColumn.maxLength - dirsColumn.maxLength - sizeColumn.maxLength
	// fullFreeSpace -= 3 // account for the two extra spaces between columns

	halfFreeSpace := termWidth/2 - treeColumn.maxLength - filesColumn.maxLength - dirsColumn.maxLength - sizeColumn.maxLength

	freeSpace := halfFreeSpace

	spaces := repeat(" ", freeSpace)

	lines := make([]string, 0)
	// header
	line := fmt.Sprintf("%s%s%s%s%s%s%s%s",
		repeat(" ", treeColumn.maxLength+freeSpace),
		ansi.Underline,
		repeat(" ", (filesColumn.maxLength+dirsColumn.maxLength)-len("items")+1)+"items",
		ansi.Reset,
		repeat(" ", 2),
		ansi.Underline,
		repeat(" ", sizeColumn.maxLength-len("size"))+"size",
		ansi.Reset,
	)
	lines = append(lines, line)

	// root
	line = fmt.Sprintf("%s%s%s",
		treeColumn.root.git,
		treeColumn.root.name,
		repeat(" ", treeColumn.maxLength-treeColumn.root.length),
	)

	// spaces
	// line += spaces
	line += repeat(" ", freeSpace)

	// files column
	line += fmt.Sprintf("%s%s",
		repeat(" ", filesColumn.maxLength-filesColumn.root.length),
		filesColumn.root.files,
	)

	line += " "

	// dirs column
	line += fmt.Sprintf("%s%s",
		repeat("-", dirsColumn.maxLength-dirsColumn.root.length),
		dirsColumn.root.dirs,
	)

	line += "  "

	// size
	line += fmt.Sprintf("%s%s",
		repeat(" ", sizeColumn.maxLength-sizeColumn.root.length),
		sizeColumn.root.size,
	)

	lines = append(lines, line)

	// dirs
	for i := range treeColumn.dirs {
		// tree column
		line := fmt.Sprintf("%s%s%s%s",
			treeColumn.dirs[i].divider,
			treeColumn.dirs[i].git,
			treeColumn.dirs[i].name,
			repeat(" ", treeColumn.maxLength-treeColumn.dirs[i].length),
		)

		// spaces
		line += spaces

		// files column
		line += fmt.Sprintf("%s%s",
			repeat(" ", filesColumn.maxLength-filesColumn.items[i].length),
			filesColumn.items[i].files,
		)

		line += " "

		// dirs column
		line += fmt.Sprintf("%s%s",
			repeat(" ", dirsColumn.maxLength-dirsColumn.items[i].length),
			dirsColumn.items[i].dirs,
		)

		line += "  "

		// size column
		line += fmt.Sprintf("%s%s",
			repeat(" ", sizeColumn.maxLength-sizeColumn.items[i].length),
			sizeColumn.items[i].size,
		)

		lines = append(lines, line)
	}
	// files
	for i := range treeColumn.files {
		// tree column
		line = fmt.Sprintf("%s%s %s%s",
			treeColumn.files[i].divider,
			treeColumn.files[i].icon,
			treeColumn.files[i].name,
			repeat(" ", treeColumn.maxLength-treeColumn.files[i].length),
		)

		// spaces
		line += spaces

		// files column
		line += repeat(" ", filesColumn.maxLength)

		line += " "

		// dirs column
		line += repeat(" ", dirsColumn.maxLength)

		line += "  "

		// size column
		line += fmt.Sprintf("%s%s",
			repeat(" ", sizeColumn.maxLength-sizeColumn.items[len(treeColumn.dirs)+i].length),
			sizeColumn.items[len(treeColumn.dirs)+i].size,
		)

		lines = append(lines, line)
	}

	fmt.Println(strings.Join(lines, "\n") + "\n")

	return nil
}

func createTreeColumn(root *builder.Directory) treeColumn {
	// root
	// git status
	rootGitStatus, rootGitStatusLength := gitStatusToString(root.GitStatus)

	// dir name
	rootName, rootNameLength := dirNameToString(root.Name, true)

	// length
	rootLength := rootGitStatusLength + rootNameLength
	maxLength := rootLength

	// dirs
	dirs := make([]subDirName, len(root.Dirs))
	for i, dir := range root.Dirs {
		// divider
		divider := treeMiddle
		if i == len(root.Dirs)-1 && len(root.Files) == 0 {
			divider = treeEnd
		}
		dividerLength := 3

		// git status
		dirGitStatus, dirGitStatusLength := gitStatusToString(dir.GitStatus)

		// name
		dirName, dirNameLength := dirNameToString(dir.Name, false)

		// length
		length := dividerLength + dirGitStatusLength + dirNameLength
		maxLength = max(maxLength, length)

		dirs[i] = subDirName{
			divider: divider,
			git:     dirGitStatus,
			name:    dirName,
			length:  length,
		}
	}

	// files
	files := make([]fileName, len(root.Files))
	for i, file := range root.Files {
		// divider
		divider := treeMiddle
		if i == len(root.Files)-1 {
			divider = treeEnd
		}
		dividerLength := 3

		// icon
		icon, iconLength := iconToString(file.Name)

		// name
		name, nameLength := fileNameToString(file.Name)

		length := dividerLength + 1 + iconLength + nameLength
		maxLength = max(maxLength, length)

		files[i] = fileName{
			divider: divider,
			icon:    icon,
			name:    name,
			length:  length,
		}
	}

	return treeColumn{
		root: rootDirName{
			git:    rootGitStatus,
			name:   rootName,
			length: rootLength,
		},
		dirs:      dirs,
		files:     files,
		maxLength: maxLength,
	}
}

func createFilesColumn(root *builder.Directory) filesColumn {
	// root
	rootFilesStr, rootFilesLength := nFilesToString(len(root.Files))

	maxLength := rootFilesLength

	// dirs
	items := make([]nFiles, len(root.Files)+len(root.Dirs))
	for i, dir := range root.Dirs {
		filesStr, filesLength := nFilesToString(len(dir.Files))
		items[i] = nFiles{
			files:  filesStr,
			length: filesLength,
		}
		maxLength = max(maxLength, filesLength)
	}

	// files
	for i := range root.Files {
		j := len(root.Dirs) + i
		items[j] = nFiles{
			files:  "",
			length: 0,
		}
	}

	return filesColumn{
		root: nFiles{
			files:  rootFilesStr,
			length: rootFilesLength,
		},
		items:     items,
		maxLength: maxLength,
	}
}

func createDirsColumn(root *builder.Directory) dirsColumn {
	// root
	rootDirsStr, rootDirsLength := nDirsToString(len(root.Dirs))

	maxLength := rootDirsLength

	// dirs
	items := make([]nDirs, len(root.Dirs)+len(root.Files))
	for i, dir := range root.Dirs {
		dirsStr, dirsLength := nDirsToString(len(dir.Dirs))
		items[i] = nDirs{
			dirs:   dirsStr,
			length: dirsLength,
		}
		maxLength = max(maxLength, dirsLength)
	}

	// files
	for i := range root.Files {
		j := len(root.Dirs) + i
		items[j] = nDirs{
			dirs:   "",
			length: 0,
		}
	}

	return dirsColumn{
		root: nDirs{
			dirs:   rootDirsStr,
			length: rootDirsLength,
		},
		items:     items,
		maxLength: maxLength,
	}
}

func createSizeColumn(root *builder.Directory) sizeColumn {
	// root
	rootSizeStr, rootSizeLength := sizeToString(root.Size)

	maxLength := rootSizeLength

	// dirs
	items := make([]size, len(root.Dirs)+len(root.Files))
	for i, dir := range root.Dirs {
		sizeStr, sizeLength := sizeToString(dir.Size)
		items[i] = size{
			size:   sizeStr,
			length: sizeLength,
		}
		maxLength = max(maxLength, sizeLength)
	}

	// files
	for i := range root.Files {
		j := len(root.Dirs) + i
		sizeStr, sizeLength := sizeToString(root.Files[i].Size)
		items[j] = size{
			size:   sizeStr,
			length: sizeLength,
		}
		maxLength = max(maxLength, sizeLength)
	}

	return sizeColumn{
		root:      size{size: rootSizeStr, length: rootSizeLength},
		items:     items,
		maxLength: maxLength,
	}
}

// =======
// helpers
// =======
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
	length := 2

	return str, length
}

func dirNameToString(name string, isRoot bool) (string, int) {
	var dirStyle string
	hidden := strings.HasPrefix(name, ".")

	switch {
	case isRoot && hidden:
		dirStyle = hiddenStyle
	case !isRoot && hidden:
		dirStyle = hiddenStyle

	case isRoot && !hidden:
		dirStyle = rootDirStyle

	case !isRoot && !hidden:
		dirStyle = subDirStyle
	}

	nameStr := fmt.Sprintf("%s%s/%s", dirStyle, name, ansi.Reset)
	length := len(name) + 1 // plus the "/" character

	return nameStr, length
}

func nDirsToString(nDirs int) (string, int) {
	if nDirs == 0 {
		return "", 0
	}

	nDirsStr := fmt.Sprintf("%d", nDirs)
	length := len(nDirsStr) + 1
	nDirsStr = fmt.Sprintf("%s%s%s%s", itemsStyle, nDirsStr, DirIcon, ansi.Reset)

	return nDirsStr, length
}

func nFilesToString(nFiles int) (string, int) {
	if nFiles == 0 {
		return "", 0
	}

	nFilesStr := fmt.Sprintf("%d", nFiles)
	length := len(nFilesStr) + 1
	nFilesStr = fmt.Sprintf("%s%s%s%s", itemsStyle, nFilesStr, FileIcon, ansi.Reset)

	return nFilesStr, length
}

func sizeToString(size int64) (string, int) {
	if size == 0 {
		return fmt.Sprintf("%sempty%s", emptyStyle, ansi.Reset), 5
	}

	sizeStr := humanizeSize(size)
	length := len(sizeStr)

	sizeStr = fmt.Sprintf("%s%s%s", sizeStyle, sizeStr, ansi.Reset)

	return sizeStr, length
}

func iconToString(name string) (string, int) {
	icon := devicons.IconForPath(name)
	iconStr := fmt.Sprintf("%s%s%s", ansi.Hex(icon.Color), icon.Icon, ansi.Reset)
	length := 1

	return iconStr, length
}

func fileNameToString(name string) (string, int) {
	style := fileStyle
	if strings.HasPrefix(name, ".") {
		style = hiddenStyle
	}

	nameStr := fmt.Sprintf("%s%s%s", style, name, ansi.Reset)
	length := len(name)

	return nameStr, length
}
