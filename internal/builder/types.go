package builder

type Args struct {
	Element string
	Sort    string
	Filter  string
	Depth   int
	Nerd    bool
}

type File struct {
	Name string
	Path string
	Size int64
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
}
