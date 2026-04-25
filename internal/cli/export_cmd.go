package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envchain/internal/env"
)

// ExportFormat maps flag string values to env.ExportFormat constants.
var exportFormats = map[string]env.ExportFormat{
	"bash":   env.FormatBash,
	"fish":   env.FormatFish,
	"dotenv": env.FormatDotenv,
}

// ExportVars writes environment variables for the given profiles to stdout
// in the requested shell format. Profiles are resolved left-to-right with
// later profiles overriding earlier ones.
//
// Usage: envchain export [--format bash|fish|dotenv] <profile> [profile...]
func (c *Commander) ExportVars(format string, profiles []string) error {
	if len(profiles) == 0 {
		return fmt.Errorf("export: at least one profile name is required")
	}

	fmt_lower := strings.ToLower(format)
	exFmt, ok := exportFormats[fmt_lower]
	if !ok {
		return fmt.Errorf("export: unknown format %q (supported: bash, fish, dotenv)", format)
	}

	r := env.NewResolver(c.store)
	ex := env.NewExporter(r)

	if err := ex.Export(os.Stdout, exFmt, profiles...); err != nil {
		return fmt.Errorf("export: %w", err)
	}
	return nil
}
