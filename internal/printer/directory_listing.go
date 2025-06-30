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
	rootDirStyle = ansi.Bold + ansi.Green
	subDirStyle  = ansi.Bold + ansi.Blue
	fileStyle    = ansi.Cyan2
	itemsStyle   = ansi.Yellow
	sizeStyle    = ansi.Magenta2
	emptyStyle   = ansi.Red2
)

// =====
// print
// =====
type dirPrint struct {
	divider       string
	dividerLength int
	git           string
	gitLength     int
	name          string
	nameLength    int
	items         string
	itemsLength   int
	size          string
	sizeLength    int
	spaces        int
}

type filePrint struct {
	divider       string
	dividerLength int
	icon          string
	iconLength    int
	name          string
	nameLength    int
	size          string
	sizeLength    int
	spaces        int
}

func printDirectoryListing(root *builder.Directory) error {
	// TODO: this does not work if the depth is greater than 1
	dirs := make([]dirPrint, len(root.Dirs)+1)
	files := make([]filePrint, len(root.Files))

	// root line
	rootDirPrint := printRootDir(root)
	dirs[0] = rootDirPrint

	// print directories
	for i, dir := range root.Dirs {
		isLast := i == len(root.Dirs)-1 && len(root.Files) == 0
		subDirPrint := printSubDir(dir, isLast)
		dirs[i+1] = subDirPrint
	}

	// print files
	for i, file := range root.Files {
		isLast := i == len(root.Files)-1
		filePrint := printFile(file, isLast)
		files[i] = filePrint
	}

	// result
	longestSizeLength := 0
	for _, dir := range dirs {
		if dir.sizeLength > longestSizeLength {
			longestSizeLength = dir.sizeLength
		}
	}
	for _, file := range files {
		if file.sizeLength > longestSizeLength {
			longestSizeLength = file.sizeLength
		}
	}

	result := ""
	for _, dir := range dirs {
		nSpaces := longestSizeLength - dir.sizeLength + 1 // +1 for the space before the size
		result += fmt.Sprintf("%s%s%s%s%s%s%s\n",
			dir.divider,
			dir.git,
			dir.name,
			repeat(" ", dir.spaces-nSpaces+1),
			dir.items,
			repeat(" ", nSpaces),
			dir.size,
		)
	}

	for _, file := range files {
		result += fmt.Sprintf("%s%s %s%s%s\n",
			file.divider,
			file.icon,
			file.name,
			repeat(" ", file.spaces),
			file.size,
		)
	}

	fmt.Print(result)

	return nil
}

// =======
// helpers
// =======
func printRootDir(dir *builder.Directory) dirPrint {
	rootGitStatus, rootGitStatusLength := gitStatusToString(dir.GitStatus)
	if rootGitStatusLength != 0 {
		rootGitStatus += " "
		rootGitStatusLength += 1 // account for the space after the git status
	}
	rootName, rootNameLenght := dirNameToString(dir.Name, true)
	rootItems, rootItemsLength := itemsToString(len(dir.Dirs), len(dir.Files))
	rootSize, rootSizeLength := sizeToString(dir.Size)
	nSpaces := termWidth - (rootGitStatusLength + rootNameLenght + rootItemsLength + rootSizeLength)

	return dirPrint{
		divider:       "",
		dividerLength: 0,
		git:           rootGitStatus,
		gitLength:     rootGitStatusLength,
		name:          rootName,
		nameLength:    rootNameLenght,
		items:         rootItems,
		itemsLength:   rootItemsLength,
		size:          rootSize,
		sizeLength:    rootSizeLength,
		spaces:        nSpaces - 2,
	}
}

func printSubDir(dir *builder.Directory, isLast bool) dirPrint {
	// divider
	divider := TreeBranch + TreeHLine + " "
	if isLast {
		divider = TreeLeaf + TreeHLine + " "
	}
	dividerLength := len(divider)

	// print
	dirGitStatus, dirGitStatusLength := gitStatusToString(dir.GitStatus)
	if dirGitStatusLength != 0 {
		dirGitStatus += " "
		dirGitStatusLength += 1 // account for the space after the git status
	}
	dirName, dirNameLength := dirNameToString(dir.Name, false)
	dirItems, dirItemsLength := itemsToString(len(dir.Dirs), len(dir.Files))
	dirSize, dirSizeLength := sizeToString(dir.Size)
	nSpaces := termWidth - (dividerLength + dirGitStatusLength + dirNameLength + dirItemsLength + dirSizeLength)

	return dirPrint{
		divider:       divider,
		dividerLength: dividerLength,
		git:           dirGitStatus,
		gitLength:     dirGitStatusLength,
		name:          dirName,
		nameLength:    dirNameLength,
		items:         dirItems,
		itemsLength:   dirItemsLength,
		size:          dirSize,
		sizeLength:    dirSizeLength,
		spaces:        nSpaces + 2,
	}
}

func printFile(file builder.File, isLast bool) filePrint {
	divider := TreeBranch + TreeHLine + " "
	if isLast {
		divider = TreeLeaf + TreeHLine + " "
	}
	dividerLength := len(divider)

	// print
	fileIcon, fileIconLength := iconToString(file.Name)
	fileName, fileNameLength := fileNameToString(file.Name)
	fileSize, fileSizeLength := sizeToString(file.Size)
	nSpaces := termWidth - (dividerLength + fileIconLength + fileNameLength + fileSizeLength)

	return filePrint{
		divider:       divider,
		dividerLength: dividerLength,
		icon:          fileIcon,
		iconLength:    fileIconLength,
		name:          fileName,
		nameLength:    fileNameLength,
		size:          fileSize,
		sizeLength:    fileSizeLength,
		spaces:        nSpaces + 2,
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

	str := fmt.Sprintf("%s%s%s", gitColor, Git, ansi.Reset)
	length := 1

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

func itemsToString(nDirs, nFiles int) (string, int) {
	itemsStr := ""
	switch {
	case nDirs == 0 && nFiles == 0:
		itemsStr = ""

	case nDirs == 0 && nFiles > 0:
		itemsStr = fmt.Sprintf("%d files", nFiles)

	case nDirs > 0 && nFiles == 0:
		itemsStr = fmt.Sprintf("%d dirs", nDirs)

	case nDirs > 0 && nFiles > 0:
		itemsStr = fmt.Sprintf("%d dirs, %d files", nDirs, nFiles)
	}

	length := len(itemsStr)
	itemsStr = fmt.Sprintf("%s%s%s", itemsStyle, itemsStr, ansi.Reset)

	return itemsStr, length
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
