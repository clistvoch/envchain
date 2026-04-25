package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ImportDotenv reads a .env file and imports all key=value pairs into the given profile.
// Lines starting with '#' and empty lines are ignored.
// Values may optionally be quoted with single or double quotes.
func (c *Commander) ImportDotenv(profileName, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("import: cannot open file %q: %w", filePath, err)
	}
	defer f.Close()

	var pairs [][2]string
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Strip optional leading "export "
		line = strings.TrimPrefix(line, "export ")

		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			return fmt.Errorf("import: invalid syntax on line %d: %q", lineNum, line)
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = unquote(val)
		pairs = append(pairs, [2]string{key, val})
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("import: read error: %w", err)
	}

	for _, kv := range pairs {
		if err := c.SetVar(profileName, kv[0], kv[1]); err != nil {
			return fmt.Errorf("import: failed to set %q: %w", kv[0], err)
		}
	}
	return nil
}

// unquote removes surrounding single or double quotes from a value string.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
