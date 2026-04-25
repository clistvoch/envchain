package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

// RegisterLintCmd registers the lint command which validates profile variable names
// and detects potential issues such as overridden keys in a chain.
func RegisterLintCmd(app *cli.App, cmdr *Commander) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "lint",
		Usage: "Validate profiles and detect shadowed or invalid variables",
		ArgsUsage: "<profile> [profile...]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "strict",
				Usage: "Exit with non-zero status if any warnings are found",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("at least one profile name is required")
			}

			profileNames := c.Args().Slice()
			warnings := lintProfiles(cmdr, profileNames)

			if len(warnings) == 0 {
				fmt.Fprintln(c.App.Writer, "OK: no issues found")
				return nil
			}

			for _, w := range warnings {
				fmt.Fprintln(c.App.Writer, "WARN:", w)
			}

			if c.Bool("strict") {
				return fmt.Errorf("%d issue(s) found", len(warnings))
			}
			return nil
		},
	})
}

// lintProfiles checks each profile for invalid variable names and shadowed keys
// when multiple profiles are provided (chain order).
func lintProfiles(cmdr *Commander, names []string) []string {
	var warnings []string
	seen := map[string]string{} // key -> first profile that defined it

	for _, name := range names {
		p, err := cmdr.store.GetProfile(name)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("profile %q not found", name))
			continue
		}

		for key := range p.Vars {
			if !isValidEnvKey(key) {
				warnings = append(warnings, fmt.Sprintf("[%s] invalid variable name: %q", name, key))
			}
			if prev, ok := seen[key]; ok {
				warnings = append(warnings, fmt.Sprintf("[%s] shadows %q already defined in %q", name, key, prev))
			} else {
				seen[key] = name
			}
		}
	}
	return warnings
}

// isValidEnvKey returns true if the key is a valid POSIX environment variable name.
func isValidEnvKey(key string) bool {
	if len(key) == 0 {
		return false
	}
	for i, ch := range key {
		switch {
		case ch >= 'A' && ch <= 'Z':
		case ch >= 'a' && ch <= 'z':
		case ch == '_':
		case ch >= '0' && ch <= '9' && i > 0:
		default:
			return false
		}
	}
	return !strings.HasPrefix(key, "__")
}
