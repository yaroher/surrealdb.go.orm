package codegen

import "strings"

func ParseAnnotation(line string) (Annotation, bool) {
	line = strings.TrimSpace(line)
	idx := strings.Index(line, "orm:")
	if idx == -1 {
		return Annotation{}, false
	}
	rest := strings.TrimSpace(line[idx+len("orm:"):])
	if rest == "" {
		return Annotation{}, false
	}
	parts := splitTokens(rest)
	if len(parts) == 0 {
		return Annotation{}, false
	}
	ann := Annotation{Kind: parts[0], Args: map[string]string{}}
	for _, p := range parts[1:] {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			ann.Args[kv[0]] = trimQuotes(kv[1])
			continue
		}
		ann.Args[p] = "true"
	}
	return ann, true
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func splitTokens(s string) []string {
	var out []string
	var buf strings.Builder
	var quote rune
	flush := func() {
		if buf.Len() > 0 {
			out = append(out, buf.String())
			buf.Reset()
		}
	}
	for _, r := range s {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
				continue
			}
			buf.WriteRune(r)
		case r == '"' || r == '\'':
			quote = r
		case r == ' ' || r == '\t':
			flush()
		default:
			buf.WriteRune(r)
		}
	}
	flush()
	return out
}
