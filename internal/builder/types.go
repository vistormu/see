package builder

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
