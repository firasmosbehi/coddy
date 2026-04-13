package sandbox

import (
	"fmt"
	"strings"
)

// ParseMemory converts memory limit string to bytes.
func ParseMemory(limit string) int64 {
	limit = strings.ToLower(strings.TrimSpace(limit))

	multiplier := int64(1)
	switch {
	case strings.HasSuffix(limit, "gb"), strings.HasSuffix(limit, "g"):
		multiplier = 1024 * 1024 * 1024
		limit = strings.TrimSuffix(strings.TrimSuffix(limit, "b"), "g")
	case strings.HasSuffix(limit, "mb"), strings.HasSuffix(limit, "m"):
		multiplier = 1024 * 1024
		limit = strings.TrimSuffix(strings.TrimSuffix(limit, "b"), "m")
	case strings.HasSuffix(limit, "kb"), strings.HasSuffix(limit, "k"):
		multiplier = 1024
		limit = strings.TrimSuffix(strings.TrimSuffix(limit, "b"), "k")
	}

	var value int64
	fmt.Sscanf(limit, "%d", &value)
	return value * multiplier
}

// DiffStringSlices returns elements in after but not in before.
func DiffStringSlices(before, after []string) []string {
	beforeSet := make(map[string]bool)
	for _, s := range before {
		beforeSet[s] = true
	}

	var result []string
	for _, s := range after {
		if !beforeSet[s] {
			result = append(result, s)
		}
	}
	return result
}
