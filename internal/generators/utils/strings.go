package utils

import (
	"regexp"
	"strings"
)

// ToTitleCase converts a string to title case
// for eg: title -> Title, created_at -> CreatedAt etc.
func ToTitleCase(s string) string {
	var re = regexp.MustCompile(`(?m)[-_ \.]`)
	parts := re.Split(s, -1)
	for i, v := range parts {
		parts[i] = ToTitle(v)
	}
	return strings.Join(parts, "")
}

func ToTitle(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func ToCamelCase(s string) string {
	title := ToTitleCase(s)
	return strings.ToLower(title[:1]) + title[1:]
}
