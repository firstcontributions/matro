package utils

import "strings"

// ToTitleCase converts a string to title case
// for eg: title -> Title, created_at -> CreatedAt etc.
func ToTitleCase(s string) string {
	parts := strings.Split(s, "_")
	for i, v := range parts {
		parts[i] = strings.Title(v)
	}
	return strings.Join(parts, "")
}

func IsElementOfStringArray(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}
