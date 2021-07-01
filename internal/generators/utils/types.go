package utils

import (
	"github.com/firstcontributions/matro/internal/parser"
)

func IsPrimitiveType(t string) bool {
	switch t {
	case parser.Int, parser.Bool, parser.Float, parser.ID, parser.String, parser.Time:
		return true
	}
	return false
}
