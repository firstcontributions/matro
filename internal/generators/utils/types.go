package utils

import (
	"fmt"

	"github.com/firstcontributions/matro/internal/parser"
)

// IsPrimitiveType says if the given type is a primitive type or not
func IsPrimitiveType(t string) bool {
	switch t {
	case parser.Int, parser.Bool, parser.Float, parser.ID, parser.String, parser.Time:
		return true
	}
	return false
}

// Counter is a closure, used in template generation, mostly for protobufs
func Counter() func(...int) string {
	count := 0
	return func(reset ...int) string {
		if len(reset) > 0 {
			count = reset[0]
			return ""
		}
		count++
		return fmt.Sprint(count)
	}
}
