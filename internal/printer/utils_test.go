package printer

import (
	"testing"

	"github.com/vistormu/go-dsa/ansi"
)

func TestVisibleWidthIgnoresAnsi(t *testing.T) {
	colored := ansi.Red + "hello" + ansi.Reset
	if visibleWidth(colored) != 5 {
		t.Fatalf("expected width 5, got %d", visibleWidth(colored))
	}
}

func TestTruncateWithEllipsis(t *testing.T) {
	got := truncateWithEllipsis("abcdefghijklmnopqrstuvwxyz", 8)
	if got != "abcdefg…" {
		t.Fatalf("unexpected truncation: %q", got)
	}
}

func TestStripAnsi(t *testing.T) {
	colored := ansi.Green + "value" + ansi.Reset
	if stripAnsi(colored) != "value" {
		t.Fatalf("unexpected stripped value: %q", stripAnsi(colored))
	}
}
