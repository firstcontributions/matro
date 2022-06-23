package utils

import (
	"strings"
)

// ToTitleCase converts a string to title case
// for eg: title -> Title, created_at -> CreatedAt etc.
func ToTitleCase(s string) string {
	parts := strings.Split(s, "_")
	for i, v := range parts {
		parts[i] = strings.Title(v)
	}
	return strings.Join(parts, "")
}

func ToCamelCase(s string) string {
	title := ToTitleCase(s)
	return strings.ToLower(title[:1]) + title[1:]
}
