package builder

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vistormu/go-dsa/errors"
)

func getGitStatus(path string) GitStatus {
	// check if the directory is a git repository
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

func buildDirectoryTree(rootPath string, depth int) (*Directory, error) {
	if depth < 0 {
		return nil, errors.New("depth cannot be negative")
	}

	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}

	dir := &Directory{
		Name:      info.Name(),
		Path:      absPath,
		Dirs:      []*Directory{},
		Files:     []File{},
		GitStatus: getGitStatus(absPath),
		Size:      0,
	}

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, errors.New("error reading directory").With("path", rootPath).Wrap(err)
	}

	var totalSize int64

	for _, entry := range entries {
		entryPath := filepath.Join(rootPath, entry.Name())
		entryInfo, err := entry.Info()
		if err != nil {
			continue // Ignore errors for entries
		}

		// If depth is greater than 0, explore subdirectories
		if entry.IsDir() && depth > 0 {
			// Recursive call with depth - 1 to prevent further recursion beyond the desired depth
			subDir, err := buildDirectoryTree(entryPath, depth-1)
			if err != nil {
				continue // Ignore errors for subdirectories
			}
			dir.Dirs = append(dir.Dirs, subDir)
			totalSize += subDir.Size
			continue
		}

		// Add file to files list if it's a file (not a directory)
		if !entry.IsDir() {
			dir.Files = append(dir.Files, File{
				Name: entry.Name(),
				Path: entryPath,
				Size: entryInfo.Size(),
			})
			totalSize += entryInfo.Size()
		}
	}

	dir.Size = totalSize
	return dir, nil
}

func buildDirectoryListing(args map[string]any) (*Directory, error) {
	// dir
	root, ok := args["dir"].(string)
	if !ok || root == "" {
		root = "."
	}

	// depth
	depth, ok := args["--depth"].(int)
	if !ok || depth < 0 {
		return nil, errors.New(WrongArg).With("arg", "--depth").With("value", args["--depth"]).With("message", "depth must be a non-negative integer")
	}

	// sort by
	// sortBy, ok := args["--sort-by"].(string)

	// build info
	rootInfo, err := buildDirectoryTree(root, depth)
	if err != nil {
		return nil, err
	}

	return rootInfo, nil
}
