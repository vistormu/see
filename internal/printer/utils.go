package printer

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

func humanizeSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func repeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	result := make([]byte, len(s)*count)
	for i := range result {
		result[i] = s[i%len(s)]
	}
	return string(result)
}

func visibleWidth(content string) int {
	return lipgloss.Width(content)
}

func padRight(content string, width int) string {
	padding := width - visibleWidth(content)
	return content + repeat(" ", padding)
}

func padLeft(content string, width int) string {
	padding := width - visibleWidth(content)
	return repeat(" ", padding) + content
}

func truncateWithEllipsis(content string, width int) string {
	if width <= 0 {
		return ""
	}
	if visibleWidth(content) <= width {
		return content
	}
	if width == 1 {
		return "…"
	}

	runes := []rune(content)
	for i := len(runes); i >= 0; i-- {
		candidate := string(runes[:i]) + "…"
		if visibleWidth(candidate) <= width {
			return candidate
		}
	}

	return "…"
}

func stripAnsi(content string) string {
	return ansiPattern.ReplaceAllString(content, "")
}

func filterLines(content, keyword string) string {
	var b strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, keyword) {
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}

	return b.String()
}
