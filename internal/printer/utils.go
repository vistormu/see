package printer

import (
	"bufio"
	"fmt"
	"strings"
)

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
