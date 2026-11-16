package yamlmini

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// Parse parses the deliberately small, data-only YAML subset used by PAD metadata.
// It supports indentation-based maps, scalar sequences, quoted strings, inline lists,
// booleans, nulls and numbers. Anchors, tags, block scalars and complex keys are rejected.
func Parse(src string) (map[string]any, error) {
	lines, err := tokenize(src)
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return map[string]any{}, nil
	}
	value, next, err := parseBlock(lines, 0, lines[0].indent)
	if err != nil {
		return nil, err
	}
	if next != len(lines) {
		return nil, fmt.Errorf("unexpected YAML content at line %d", lines[next].line)
	}
	m, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("YAML document must be a mapping")
	}
	return m, nil
}

type line struct {
	indent int
	text   string
	line   int
}

func tokenize(src string) ([]line, error) {
	s := bufio.NewScanner(strings.NewReader(src))
	out := []line{}
	n := 0
	for s.Scan() {
		n++
		raw := strings.TrimRight(s.Text(), " \t\r")
		if strings.TrimSpace(raw) == "" || strings.HasPrefix(strings.TrimSpace(raw), "#") {
			continue
		}
		indent := len(raw) - len(strings.TrimLeft(raw, " "))
		if strings.Contains(raw[:indent], "\t") {
			return nil, fmt.Errorf("tabs are not allowed for YAML indentation at line %d", n)
		}
		t := strings.TrimSpace(raw)
		if strings.Contains(t, "\t") {
			return nil, fmt.Errorf("tabs are not supported at line %d", n)
		}
		out = append(out, line{indent, t, n})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
func parseBlock(ls []line, i, indent int) (any, int, error) {
	if i >= len(ls) {
		return nil, i, nil
	}
	if strings.HasPrefix(ls[i].text, "- ") || ls[i].text == "-" {
		return parseSeq(ls, i, indent)
	}
	return parseMap(ls, i, indent)
}
func parseMap(ls []line, i, indent int) (map[string]any, int, error) {
	m := map[string]any{}
	for i < len(ls) && ls[i].indent == indent && !strings.HasPrefix(ls[i].text, "- ") && ls[i].text != "-" {
		l := ls[i]
		k, v, ok := strings.Cut(l.text, ":")
		if !ok {
			return nil, i, fmt.Errorf("expected key: value at line %d", l.line)
		}
		k = strings.TrimSpace(k)
		if k == "" || strings.ContainsAny(k, "{}[]") {
			return nil, i, fmt.Errorf("invalid mapping key at line %d", l.line)
		}
		if _, exists := m[k]; exists {
			return nil, i, fmt.Errorf("duplicate key %q at line %d", k, l.line)
		}
		v = strings.TrimSpace(v)
		i++
		if v != "" {
			x, e := scalar(stripComment(v))
			if e != nil {
				return nil, i, e
			}
			m[k] = x
			continue
		}
		if i < len(ls) && ls[i].indent > indent {
			x, next, e := parseBlock(ls, i, ls[i].indent)
			if e != nil {
				return nil, i, e
			}
			m[k] = x
			i = next
		} else {
			m[k] = nil
		}
	}
	return m, i, nil
}
func parseSeq(ls []line, i, indent int) ([]any, int, error) {
	a := []any{}
	for i < len(ls) && ls[i].indent == indent && (strings.HasPrefix(ls[i].text, "- ") || ls[i].text == "-") {
		t := strings.TrimSpace(strings.TrimPrefix(ls[i].text, "-"))
		ln := ls[i].line
		i++
		if t != "" {
			x, e := scalar(stripComment(t))
			if e != nil {
				return nil, i, e
			}
			a = append(a, x)
		} else if i < len(ls) && ls[i].indent > indent {
			x, next, e := parseBlock(ls, i, ls[i].indent)
			if e != nil {
				return nil, i, e
			}
			a = append(a, x)
			i = next
		} else {
			a = append(a, nil)
		}
		_ = ln
	}
	return a, i, nil
}
func stripComment(s string) string {
	q := rune(0)
	esc := false
	for i, r := range s {
		if esc {
			esc = false
			continue
		}
		if q == '"' && r == '\\' {
			esc = true
			continue
		}
		if q != 0 {
			if r == q {
				q = 0
			}
			continue
		}
		if r == '\'' || r == '"' {
			q = r
			continue
		}
		if r == '#' && (i == 0 || s[i-1] == ' ') {
			return strings.TrimSpace(s[:i])
		}
	}
	return strings.TrimSpace(s)
}
func scalar(s string) (any, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "null" || s == "~" {
		return nil, nil
	}
	if s == "true" {
		return true, nil
	}
	if s == "false" {
		return false, nil
	}
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		inner := strings.TrimSpace(s[1 : len(s)-1])
		if inner == "" {
			return []any{}, nil
		}
		parts := strings.Split(inner, ",")
		a := make([]any, 0, len(parts))
		for _, p := range parts {
			x, e := scalar(strings.TrimSpace(p))
			if e != nil {
				return nil, e
			}
			a = append(a, x)
		}
		return a, nil
	}
	if (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) || (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) {
		if s[0] == '\'' {
			return strings.ReplaceAll(s[1:len(s)-1], "''", "'"), nil
		}
		x, e := strconv.Unquote(s)
		if e != nil {
			return nil, fmt.Errorf("invalid quoted scalar %s: %w", s, e)
		}
		return x, nil
	}
	if strings.HasPrefix(s, "&") || strings.HasPrefix(s, "*") || strings.HasPrefix(s, "!") || strings.HasPrefix(s, "|") || strings.HasPrefix(s, ">") {
		return nil, fmt.Errorf("unsupported YAML feature in %q", s)
	}
	if n, e := strconv.ParseInt(s, 10, 64); e == nil {
		return n, nil
	}
	if f, e := strconv.ParseFloat(s, 64); e == nil {
		return f, nil
	}
	return s, nil
}
