package printer

import (
	"os"
	"strings"
	"testing"

	"see/internal/builder"
)

func TestRenderDirectoryListing_ShowsPermissionsColumnByDefault(t *testing.T) {
	root := &builder.Directory{
		Name: "root",
		Mode: os.FileMode(0o755) | os.ModeDir,
		Dirs: []*builder.Directory{
			{
				Name: "sub",
				Mode: os.FileMode(0o700) | os.ModeDir,
			},
		},
		Files: []builder.File{
			{
				Name: "main.go",
				Size: 12,
				Mode: 0o644,
			},
		},
	}

	output, copyOutput := renderDirectoryListing(root, builder.Args{})
	if !strings.Contains(output, "perms") {
		t.Fatalf("expected permissions header in output: %q", output)
	}
	if !strings.Contains(copyOutput, "drwx------") {
		t.Fatalf("expected directory permissions in copy output: %q", copyOutput)
	}
	if !strings.Contains(copyOutput, "-rw-r--r--") {
		t.Fatalf("expected file permissions in copy output: %q", copyOutput)
	}
	if strings.Contains(copyOutput, "\x1b[") {
		t.Fatalf("expected ANSI-free copy output, got %q", copyOutput)
	}
}

func TestRenderDirectoryListing_DirsOnlyIsNotEmpty(t *testing.T) {
	root := &builder.Directory{
		Name: "root",
		Mode: os.FileMode(0o755) | os.ModeDir,
		Dirs: []*builder.Directory{
			{
				Name: "a",
				Mode: os.FileMode(0o755) | os.ModeDir,
				Dirs: []*builder.Directory{
					{
						Name: "b",
						Mode: os.FileMode(0o755) | os.ModeDir,
					},
				},
			},
		},
	}

	_, copyOutput := renderDirectoryListing(root, builder.Args{})
	if !strings.Contains(copyOutput, "0 B") {
		t.Fatalf("expected zero-byte size marker for non-empty dirs, got %q", copyOutput)
	}

	lines := strings.Split(copyOutput, "\n")
	for _, line := range lines {
		if strings.Contains(line, " a/") && strings.Contains(line, "empty") {
			t.Fatalf("expected directory with child dirs to avoid empty marker, got line %q", line)
		}
	}
}

func TestRenderDirectoryListing_TruncatesLongNames(t *testing.T) {
	previousWidth := termWidth
	termWidth = 52
	t.Cleanup(func() { termWidth = previousWidth })

	root := &builder.Directory{
		Name: "very-very-long-root-directory-name",
		Mode: os.FileMode(0o755) | os.ModeDir,
		Files: []builder.File{
			{
				Name: "this-is-a-very-long-filename-that-should-be-truncated.txt",
				Size: 1,
				Mode: 0o644,
			},
		},
	}

	output, _ := renderDirectoryListing(root, builder.Args{})
	if !strings.Contains(output, "…") {
		t.Fatalf("expected ellipsis in output, got %q", output)
	}

	for _, line := range strings.Split(output, "\n") {
		if visibleWidth(stripAnsi(line)) > termWidth {
			t.Fatalf("line exceeds terminal width: %q", line)
		}
	}
}

func TestRenderDirectoryListing_RootIsInsideFrame(t *testing.T) {
	root := &builder.Directory{
		Name: "root",
		Mode: os.FileMode(0o755) | os.ModeDir,
	}

	output, copyOutput := renderDirectoryListing(root, builder.Args{})
	lines := strings.Split(copyOutput, "\n")
	if len(lines) < 4 {
		t.Fatalf("unexpected short table output: %q", copyOutput)
	}
	if !strings.HasPrefix(lines[0], "┌") {
		t.Fatalf("expected table to start with frame top, got %q", lines[0])
	}
	if !strings.Contains(copyOutput, "root/") {
		t.Fatalf("expected root row inside frame: %q", copyOutput)
	}
	if strings.Contains(output, "\nroot/\n") {
		t.Fatalf("root should not be printed outside frame: %q", output)
	}
}

func TestRenderDirectoryListing_DirectItemsCounts(t *testing.T) {
	root := &builder.Directory{
		Name: "root",
		Mode: os.FileMode(0o755) | os.ModeDir,
		Dirs: []*builder.Directory{
			{
				Name: "a",
				Mode: os.FileMode(0o755) | os.ModeDir,
				Files: []builder.File{
					{Name: "x.txt", Size: 1, Mode: 0o644},
				},
				Dirs: []*builder.Directory{
					{
						Name: "b",
						Mode: os.FileMode(0o755) | os.ModeDir,
						Files: []builder.File{
							{Name: "y.txt", Size: 2, Mode: 0o644},
						},
					},
				},
			},
		},
	}

	_, copyOutput := renderDirectoryListing(root, builder.Args{})
	if !strings.Contains(copyOutput, "1"+DirIcon) {
		t.Fatalf("expected root direct item count, got %q", copyOutput)
	}
	if !strings.Contains(copyOutput, "1"+FileIcon+" 1"+DirIcon) {
		t.Fatalf("expected direct subdir item counts, got %q", copyOutput)
	}
}

func TestRenderDirectoryListing_TreeDepthUsesDepthButMetadataUsesDepthPlusOne(t *testing.T) {
	root := &builder.Directory{
		Name: "root",
		Mode: os.FileMode(0o755) | os.ModeDir,
		Dirs: []*builder.Directory{
			{
				Name: "a",
				Mode: os.FileMode(0o755) | os.ModeDir,
				Size: 3,
				Files: []builder.File{
					{Name: "x.txt", Size: 1, Mode: 0o644},
				},
				Dirs: []*builder.Directory{
					{
						Name: "b",
						Mode: os.FileMode(0o755) | os.ModeDir,
						Size: 2,
						Files: []builder.File{
							{Name: "y.txt", Size: 2, Mode: 0o644},
						},
					},
				},
			},
		},
	}

	_, copyOutput := renderDirectoryListing(root, builder.Args{Depth: 1})
	if !strings.Contains(copyOutput, "a/") {
		t.Fatalf("expected first-level dir to be rendered, got %q", copyOutput)
	}
	if strings.Contains(copyOutput, "b/") {
		t.Fatalf("expected second-level dir to be hidden at depth 1, got %q", copyOutput)
	}
	if !strings.Contains(copyOutput, "1"+FileIcon+" 1"+DirIcon) {
		t.Fatalf("expected displayed dir metadata to include depth+1 children counts, got %q", copyOutput)
	}
	if !strings.Contains(copyOutput, "3 B") {
		t.Fatalf("expected displayed dir size metadata to be present, got %q", copyOutput)
	}
}
