// Package env provides environment variable resolution and export utilities.
package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportFormat represents the shell export format.
type ExportFormat int

const (
	// FormatBash generates bash/sh compatible export statements.
	FormatBash ExportFormat = iota
	// FormatFish generates fish shell compatible set statements.
	FormatFish
	// FormatDotenv generates .env file format.
	FormatDotenv
)

// Exporter writes resolved environment variables to an output stream.
type Exporter struct {
	resolver *Resolver
}

// NewExporter creates a new Exporter backed by the given Resolver.
func NewExporter(r *Resolver) *Exporter {
	return &Exporter{resolver: r}
}

// Export writes environment variables for the given profiles to w in the specified format.
func (e *Exporter) Export(w io.Writer, format ExportFormat, profiles ...string) error {
	vars, err := e.resolver.Resolve(profiles...)
	if err != nil {
		return fmt.Errorf("export: resolve failed: %w", err)
	}

	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := vars[k]
		var line string
		switch format {
		case FormatFish:
			line = fmt.Sprintf("set -x %s %s\n", k, shellEscape(v))
		case FormatDotenv:
			line = fmt.Sprintf("%s=%s\n", k, v)
		default: // FormatBash
			line = fmt.Sprintf("export %s=%s\n", k, shellEscape(v))
		}
		if _, err := io.WriteString(w, line); err != nil {
			return fmt.Errorf("export: write failed: %w", err)
		}
	}
	return nil
}

// shellEscape wraps a value in single quotes, escaping any existing single quotes.
func shellEscape(v string) string {
	escaped := strings.ReplaceAll(v, "'", "'\\''")
	return "'" + escaped + "'"
}
