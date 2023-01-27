package helpers

import "strings"

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func Every(slice []int, item int) bool {
	for _, s := range slice {
		if s != item {
			return false
		}
	}
	return true
}

func SneakCase(s string) string {
	var result string
	for i, r := range s {
		if i > 1 && (r >= 'A' && r <= 'Z') {
			result += "_"
		}
		result += string(r)
	}
	return strings.ToLower(result)
}
