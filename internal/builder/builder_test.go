package builder

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
)

func TestBuildCommand_UsesEnvFallbackWhenPathMissing(t *testing.T) {
	tmpDir := t.TempDir()
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalWD) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	const key = "SEE_TEST_ENV_FALLBACK"
	const value = "from-env"
	t.Setenv(key, value)

	command, err := BuildCommand(Args{Element: key, Sort: "name", Depth: 1})
	if err != nil {
		t.Fatalf("BuildCommand returned error: %v", err)
	}

	envVar, ok := command.(*EnvVariable)
	if !ok {
		t.Fatalf("expected *EnvVariable, got %T", command)
	}

	if envVar.Name != key || envVar.Value != value {
		t.Fatalf("unexpected env result: %#v", envVar)
	}
}

func TestBuildCommand_PrefersExistingPathOverEnv(t *testing.T) {
	tmpDir := t.TempDir()
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalWD) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	const key = "SEE_TEST_PATH_PRECEDENCE"
	t.Setenv(key, "from-env")

	if err := os.WriteFile(filepath.Join(tmpDir, key), []byte("file-content"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	command, err := BuildCommand(Args{Element: key, Sort: "name", Depth: 1})
	if err != nil {
		t.Fatalf("BuildCommand returned error: %v", err)
	}

	if _, ok := command.(*EnvVariable); ok {
		t.Fatalf("expected file command, got env variable")
	}
	if _, ok := command.(*FileContent); !ok {
		t.Fatalf("expected *FileContent, got %T", command)
	}
}

func TestBuildCommand_UsesEnvFallbackWhenArgIsExpandedValue(t *testing.T) {
	tmpDir := t.TempDir()
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalWD) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	const key = "SEE_TEST_DISPLAY"
	const value = ":42"
	t.Setenv(key, value)

	command, err := BuildCommand(Args{Element: value, Sort: "name", Depth: 1})
	if err != nil {
		t.Fatalf("BuildCommand returned error: %v", err)
	}

	envVar, ok := command.(*EnvVariable)
	if !ok {
		t.Fatalf("expected *EnvVariable, got %T", command)
	}
	if envVar.Name != key || envVar.Value != value {
		t.Fatalf("unexpected env result: %#v", envVar)
	}
}

func TestBuildCommand_UsesEnvFallbackForNonRegularPathValue(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("mkfifo not supported on windows")
	}

	tmpDir := t.TempDir()
	fifoPath := filepath.Join(tmpDir, "display.sock")
	if err := syscall.Mkfifo(fifoPath, 0o600); err != nil {
		t.Fatalf("mkfifo failed: %v", err)
	}

	const key = "SEE_TEST_FIFO_DISPLAY"
	t.Setenv(key, fifoPath)

	command, err := BuildCommand(Args{Element: fifoPath, Sort: "name", Depth: 1})
	if err != nil {
		t.Fatalf("BuildCommand returned error: %v", err)
	}

	envVar, ok := command.(*EnvVariable)
	if !ok {
		t.Fatalf("expected *EnvVariable, got %T", command)
	}
	if envVar.Value != fifoPath {
		t.Fatalf("unexpected env result: %#v", envVar)
	}
}

func TestBuildDirectoryTree_DepthLimitAndUnlimited(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, "a", "b"), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("r"), 0o644); err != nil {
		t.Fatalf("write root file failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "a", "child.txt"), []byte("cc"), 0o644); err != nil {
		t.Fatalf("write child file failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "a", "b", "deep.txt"), []byte("ddd"), 0o644); err != nil {
		t.Fatalf("write deep file failed: %v", err)
	}

	depthOne, err := buildDirectoryTree(tmpDir, 1)
	if err != nil {
		t.Fatalf("buildDirectoryTree depth 1 failed: %v", err)
	}
	if len(depthOne.Dirs) != 1 {
		t.Fatalf("expected one subdir at depth 1, got %d", len(depthOne.Dirs))
	}
	if len(depthOne.Files) != 1 {
		t.Fatalf("expected one root file at depth 1, got %d", len(depthOne.Files))
	}
	if len(depthOne.Dirs[0].Dirs) != 0 || len(depthOne.Dirs[0].Files) != 0 {
		t.Fatalf("expected no nested entries at depth 1, got dirs=%d files=%d", len(depthOne.Dirs[0].Dirs), len(depthOne.Dirs[0].Files))
	}
	if depthOne.Size != 1 {
		t.Fatalf("expected depth-limited size 1, got %d", depthOne.Size)
	}

	unlimited, err := buildDirectoryTree(tmpDir, 0)
	if err != nil {
		t.Fatalf("buildDirectoryTree unlimited failed: %v", err)
	}
	if len(unlimited.Dirs) != 1 || len(unlimited.Dirs[0].Dirs) != 1 {
		t.Fatalf("expected full nested tree, got %#v", unlimited.Dirs)
	}
	if unlimited.Size != 6 {
		t.Fatalf("expected full recursive size 6, got %d", unlimited.Size)
	}
}

func TestBuildDirectoryListing_UsesDepthPlusOneForMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, "a", "b"), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "a", "child.txt"), []byte("cc"), 0o644); err != nil {
		t.Fatalf("write child file failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "a", "b", "deep.txt"), []byte("ddd"), 0o644); err != nil {
		t.Fatalf("write deep file failed: %v", err)
	}

	listing, err := buildDirectoryListing(Args{
		Element: tmpDir,
		Depth:   1,
	})
	if err != nil {
		t.Fatalf("buildDirectoryListing failed: %v", err)
	}

	if len(listing.Dirs) != 1 {
		t.Fatalf("expected one root child, got %d", len(listing.Dirs))
	}

	a := listing.Dirs[0]
	if len(a.Files) != 1 || len(a.Dirs) != 1 {
		t.Fatalf("expected depth+1 metadata on displayed dir, got files=%d dirs=%d", len(a.Files), len(a.Dirs))
	}
	if len(a.Dirs[0].Files) != 0 {
		t.Fatalf("expected metadata traversal to stop after depth+1, got nested files=%d", len(a.Dirs[0].Files))
	}
}
