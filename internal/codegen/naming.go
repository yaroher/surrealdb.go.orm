package codegen

import (
	"strings"
	"unicode"
)

func applyRenameAll(name, renameAll string) string {
	switch renameAll {
	case "camelCase":
		return toCamel(name)
	case "PascalCase":
		return toPascal(name)
	case "snake_case":
		return toSnake(name)
	case "SCREAMING_SNAKE_CASE":
		return strings.ToUpper(toSnake(name))
	case "lowercase":
		return strings.ToLower(name)
	case "UPPERCASE":
		return strings.ToUpper(name)
	default:
		return name
	}
}

func ApplyRenameAll(name, renameAll string) string {
	return applyRenameAll(name, renameAll)
}

func ToSnake(s string) string {
	return toSnake(s)
}

func toSnake(s string) string {
	words := splitWords(s)
	return strings.ToLower(strings.Join(words, "_"))
}

func toPascal(s string) string {
	words := splitWords(s)
	for i := range words {
		words[i] = capitalize(strings.ToLower(words[i]))
	}
	return strings.Join(words, "")
}

func toCamel(s string) string {
	p := toPascal(s)
	if p == "" {
		return p
	}
	return strings.ToLower(p[:1]) + p[1:]
}

func splitWords(s string) []string {
	var out []string
	var current []rune
	flush := func() {
		if len(current) == 0 {
			return
		}
		out = append(out, string(current))
		current = current[:0]
	}
	runes := []rune(s)
	for i, r := range runes {
		if r == '_' || r == '-' || r == ' ' {
			flush()
			continue
		}
		if i > 0 {
			prev := runes[i-1]
			var next rune
			if i+1 < len(runes) {
				next = runes[i+1]
			}
			switch {
			case unicode.IsUpper(r) && (unicode.IsLower(prev) || unicode.IsDigit(prev)):
				flush()
			case unicode.IsUpper(r) && unicode.IsUpper(prev) && next != 0 && unicode.IsLower(next):
				flush()
			case unicode.IsDigit(r) && !unicode.IsDigit(prev):
				flush()
			case !unicode.IsDigit(r) && unicode.IsDigit(prev):
				flush()
			}
		}
		current = append(current, r)
	}
	flush()
	return out
}

func charClass(r rune) int {
	switch {
	case unicode.IsDigit(r):
		return 1
	case unicode.IsUpper(r):
		return 2
	case unicode.IsLower(r):
		return 3
	default:
		return 4
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
