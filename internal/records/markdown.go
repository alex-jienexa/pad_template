package records

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alex-jienexa/pad_template/internal/yamlmini"
)

type Document struct {
	Metadata map[string]any
	Body     string
}

func ParseMarkdown(data []byte) (Document, error) {
	s := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return Document{}, fmt.Errorf("missing YAML front matter opener")
	}
	rest := s[4:]
	idx := strings.Index(rest, "\n---\n")
	if idx < 0 {
		if strings.HasSuffix(rest, "\n---") {
			idx = len(rest) - 4
		} else {
			return Document{}, fmt.Errorf("missing YAML front matter closer")
		}
	}
	meta, err := yamlmini.Parse(rest[:idx])
	if err != nil {
		return Document{}, err
	}
	body := ""
	end := idx + 5
	if end < len(rest) {
		body = rest[end:]
	}
	return Document{meta, body}, nil
}
func Read(path string) (Document, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return Document{}, e
	}
	return ParseMarkdown(b)
}
func Write(path string, metadata map[string]any, body string) error {
	var b strings.Builder
	b.WriteString("---\n")
	// Stable scalar/list YAML emitter; PAD metadata intentionally avoids nested structures in front matter.
	keys := []string{"id", "title", "type", "status", "priority", "assignee", "owner", "created_at", "updated_at", "estimate", "due_date", "parent", "depends_on", "requirement_ids", "related_requirements", "related_tasks", "labels", "files_owned", "verification", "date", "deciders"}
	seen := map[string]bool{}
	for _, k := range keys {
		v, ok := metadata[k]
		if !ok {
			continue
		}
		seen[k] = true
		writeField(&b, k, v)
	}
	for k, v := range metadata {
		if !seen[k] {
			writeField(&b, k, v)
		}
	}
	b.WriteString("---\n\n")
	b.WriteString(strings.TrimLeft(body, "\n"))
	if !strings.HasSuffix(b.String(), "\n") {
		b.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
func writeField(b *strings.Builder, k string, v any) {
	switch x := v.(type) {
	case nil:
		fmt.Fprintf(b, "%s: null\n", k)
	case []string:
		if len(x) == 0 {
			fmt.Fprintf(b, "%s: []\n", k)
			break
		}
		fmt.Fprintf(b, "%s:\n", k)
		for _, s := range x {
			fmt.Fprintf(b, "  - %s\n", quote(s))
		}
	case []any:
		if len(x) == 0 {
			fmt.Fprintf(b, "%s: []\n", k)
			break
		}
		fmt.Fprintf(b, "%s:\n", k)
		for _, s := range x {
			fmt.Fprintf(b, "  - %s\n", quote(fmt.Sprint(s)))
		}
	case bool:
		fmt.Fprintf(b, "%s: %t\n", k, x)
	case int, int64, float64:
		fmt.Fprintf(b, "%s: %v\n", k, x)
	default:
		fmt.Fprintf(b, "%s: %s\n", k, quote(fmt.Sprint(v)))
	}
}
func quote(s string) string {
	if s == "" {
		return "\"\""
	}
	if strings.ContainsAny(s, ":#[]{}\n\r\t") || strings.HasPrefix(s, " ") || strings.HasSuffix(s, " ") || s == "null" || s == "true" || s == "false" || s == "~" {
		return fmt.Sprintf("%q", s)
	}
	return s
}
func Date() string                             { return time.Now().Format("2006-01-02") }
func Timestamp() string                        { return time.Now().Format(time.RFC3339) }
func String(m map[string]any, k string) string { v, _ := m[k].(string); return v }
func Strings(m map[string]any, k string) []string {
	v, ok := m[k].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, fmt.Sprint(x))
	}
	return out
}
