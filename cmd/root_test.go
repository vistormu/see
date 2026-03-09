package cmd

import (
	"io"
	"strings"
	"testing"
)

func resetTestFlags() {
	sortBy = "name"
	version = false
	help = false
	filter = ""
	depth = 1
	nerd = false
	head = -1
	tail = -1
	copyOutput = false
}

func TestNewFlagShorthands(t *testing.T) {
	resetTestFlags()

	tests := []struct {
		name      string
		shorthand string
	}{
		{name: "head", shorthand: "H"},
		{name: "tail", shorthand: "t"},
		{name: "copy", shorthand: "c"},
	}

	for _, tc := range tests {
		flag := rootCmd.Flags().Lookup(tc.name)
		if flag == nil {
			t.Fatalf("flag %q not found", tc.name)
		}
		if flag.Shorthand != tc.shorthand {
			t.Fatalf("flag %q shorthand mismatch: got %q want %q", tc.name, flag.Shorthand, tc.shorthand)
		}
	}
}

func TestNormalizeOptionalIntFlags(t *testing.T) {
	got := normalizeOptionalIntFlags([]string{"--head", "README.md"})
	want := []string{"--head", "10", "README.md"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("unexpected normalized args: got %#v want %#v", got, want)
	}
}

func TestExecute_HeadTailConflictReturnsSingleError(t *testing.T) {
	resetTestFlags()

	rootCmd.SetOut(io.Discard)
	rootCmd.SetErr(io.Discard)
	rootCmd.SetArgs([]string{"--head=2", "--tail=2", "README.md"})

	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Fatal("expected error for --head and --tail conflict")
	}
	if !strings.Contains(err.Error(), "cannot be used together") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_HeadOnDirectoryReturnsError(t *testing.T) {
	resetTestFlags()

	rootCmd.SetOut(io.Discard)
	rootCmd.SetErr(io.Discard)
	rootCmd.SetArgs([]string{"--head=2", "."})

	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Fatal("expected error when --head is used with directory output")
	}
	if !strings.Contains(err.Error(), "only be used when showing file or env content") {
		t.Fatalf("unexpected error: %v", err)
	}
}
