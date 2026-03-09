package printer

import (
	"io"
	"os"
	"strings"
	"testing"

	"see/internal/builder"
)

func TestSelectContentWindow(t *testing.T) {
	content := "a\nb\nc\nd"

	head, headStart := selectContentWindow(content, 2, -1)
	if head != "a\nb" || headStart != 1 {
		t.Fatalf("unexpected head selection: %q start=%d", head, headStart)
	}

	tail, tailStart := selectContentWindow(content, -1, 2)
	if tail != "c\nd" || tailStart != 3 {
		t.Fatalf("unexpected tail selection: %q start=%d", tail, tailStart)
	}
}

func TestAddLineNumbers_Alignment(t *testing.T) {
	got := addLineNumbers("x\ny\nz", 98)
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], " 98 │ ") {
		t.Fatalf("expected aligned start line, got %q", lines[0])
	}
	if !strings.Contains(lines[2], "100 │ ") {
		t.Fatalf("expected aligned third line, got %q", lines[2])
	}
}

func TestPrintFileContent_CopyOmitsLineNumbers(t *testing.T) {
	originalCopyFn := copyFn
	var copied string
	copyFn = func(content string) error {
		copied = content
		return nil
	}
	t.Cleanup(func() { copyFn = originalCopyFn })

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe failed: %v", err)
	}
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = oldStdout })

	fileContent := &builder.FileContent{
		File: &builder.File{
			Name: "test.txt",
			Path: "test.txt",
			Size: 5,
		},
		Content: "line1\nline2\nline3",
		NLines:  3,
	}

	if err := printFileContent(fileContent, builder.Args{Copy: true, Head: 2}); err != nil {
		t.Fatalf("printFileContent failed: %v", err)
	}

	_ = w.Close()
	_, _ = io.ReadAll(r)
	_ = r.Close()

	if strings.Contains(copied, "│") {
		t.Fatalf("copy payload should not include line-number separators: %q", copied)
	}
	if copied != "line1\nline2" {
		t.Fatalf("unexpected copy payload: %q", copied)
	}
}

func TestPrintFileContent_UsesFrameAndWrapsLongLines(t *testing.T) {
	originalCopyFn := copyFn
	copyFn = func(content string) error { return nil }
	t.Cleanup(func() { copyFn = originalCopyFn })

	previousWidth := termWidth
	termWidth = 40
	t.Cleanup(func() { termWidth = previousWidth })

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe failed: %v", err)
	}
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = oldStdout })

	fileContent := &builder.FileContent{
		File: &builder.File{
			Name: "test.txt",
			Path: "test.txt",
			Size: 40,
		},
		Content: "this is a very long line that should wrap in the framed output",
		NLines:  1,
	}

	if err := printFileContent(fileContent, builder.Args{}); err != nil {
		t.Fatalf("printFileContent failed: %v", err)
	}

	_ = w.Close()
	outBytes, _ := io.ReadAll(r)
	_ = r.Close()

	out := stripAnsi(string(outBytes))
	if !strings.Contains(out, "┌") || !strings.Contains(out, "└") {
		t.Fatalf("expected framed output, got %q", out)
	}

	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if visibleWidth(line) > termWidth {
			t.Fatalf("wrapped line exceeds terminal width: %q", line)
		}
	}
}
