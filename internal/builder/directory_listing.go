package builder

import (
	"bytes"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getGitStatus(path string) GitStatus {
	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return GitUnknown
	}

	cmd := exec.Command("git", "-C", path, "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return GitUnknown
	}

	status := out.String()
	if status == "" {
		return GitUpToDate
	}
	if strings.Contains(status, "??") {
		return GitUntracked
	}
	return GitModified
}

func buildDirectoryTree(root string, maxDepth int) (*Directory, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	rootInfo, err := os.Stat(rootAbs)
	if err != nil {
		return nil, err
	}

	// map absolute dir path → *Directory nodes for depth ≤ maxDepth
	dirMap := map[string]*Directory{}
	rootDir := &Directory{
		Name:      filepath.Base(rootAbs),
		Path:      rootAbs,
		GitStatus: getGitStatus(rootAbs),
		Mode:      rootInfo.Mode(),
	}
	dirMap[rootAbs] = rootDir

	filepath.WalkDir(rootAbs, func(path string, d fs.DirEntry, wErr error) error {
		if wErr != nil {
			return nil
		}
		if path == rootAbs {
			return nil
		}

		// compute depth relative to root
		rel, _ := filepath.Rel(rootAbs, path)
		depth := len(strings.Split(rel, string(os.PathSeparator)))
		limitedDepth := maxDepth > 0

		if limitedDepth && depth > maxDepth {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}

			dir := &Directory{
				Name:      d.Name(),
				Path:      path,
				GitStatus: getGitStatus(path),
				Mode:      info.Mode(),
			}
			dirMap[path] = dir
			parent := dirMap[filepath.Dir(path)]
			parent.Dirs = append(parent.Dirs, dir)

			return nil
		}

		// it's a file: get its size
		info, err := d.Info()
		if err != nil {
			return err
		}
		size := info.Size()
		f := File{
			Name: d.Name(),
			Path: path,
			Size: size,
			Mode: info.Mode(),
		}

		parentPath := filepath.Dir(path)
		if pd, ok := dirMap[parentPath]; ok {
			pd.Files = append(pd.Files, f)
			// propagate this file’s size up to all ancestors in dirMap
			for p := parentPath; ; {
				if dd, exists := dirMap[p]; exists {
					dd.Size += size
				}
				if p == rootAbs {
					break
				}
				p = filepath.Dir(p)
			}
		}

		return nil
	})

	return rootDir, nil
}

func buildDirectoryListing(args Args) (*Directory, error) {
	metadataDepth := args.Depth
	if metadataDepth > 0 {
		metadataDepth++
	}

	// build info
	rootInfo, err := buildDirectoryTree(args.Element, metadataDepth)
	if err != nil {
		return nil, err
	}

	return rootInfo, nil
}
