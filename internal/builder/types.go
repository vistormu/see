package builder

import "os"

type Args struct {
	Element string
	Sort    string
	Filter  string
	Depth   int
	Nerd    bool
	Head    int
	Tail    int
	Copy    bool
}

type File struct {
	Name string
	Path string
	Size int64
	Mode os.FileMode
}

type GitStatus int

const (
	GitUpToDate GitStatus = iota
	GitModified
	GitUntracked
	GitUnknown
)

type Directory struct {
	Name      string
	Path      string
	Dirs      []*Directory
	Files     []File
	GitStatus GitStatus
	Size      int64
	Mode      os.FileMode
}

type EnvVariable struct {
	Name  string
	Value string
}
