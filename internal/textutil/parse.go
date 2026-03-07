package textutil

import "strings"

func ParseKeyValue(entry string) (string, string, bool) {
	parts := strings.SplitN(entry, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" || value == "" {
		return "", "", false
	}
	return key, value, true
}

func SplitCSV(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	items := strings.Split(input, ",")
	out := make([]string, 0, len(items))
	for _, item := range items {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}

func SplitNonEmptyLines(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	lines := strings.Split(input, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		value := strings.TrimSpace(line)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}
