package cli

import (
	"fmt"
	"sort"

	"github.com/urfave/cli/v2"
)

// RegisterDiffCmd registers the diff subcommand which shows differences
// between two profiles' environment variable sets.
func RegisterDiffCmd(app *cli.App, cmdr *Commander) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:      "diff",
		Usage:     "Show differences between two profiles",
		ArgsUsage: "<profile-a> <profile-b>",
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
				return fmt.Errorf("usage: envchain diff <profile-a> <profile-b>")
			}
			nameA := c.Args().Get(0)
			nameB := c.Args().Get(1)
			return cmdr.DiffProfiles(nameA, nameB)
		},
	})
}

// DiffProfiles prints the variable differences between two profiles.
// Lines prefixed with '-' exist only in profile A, '+' only in profile B,
// and '~' exist in both but with different values.
func (cmdr *Commander) DiffProfiles(nameA, nameB string) error {
	profA, err := cmdr.store.GetProfile(nameA)
	if err != nil {
		return fmt.Errorf("profile %q not found", nameA)
	}
	profB, err := cmdr.store.GetProfile(nameB)
	if err != nil {
		return fmt.Errorf("profile %q not found", nameB)
	}

	varsA := profA.Vars()
	varsB := profB.Vars()

	keys := make(map[string]struct{})
	for k := range varsA {
		keys[k] = struct{}{}
	}
	for k := range varsB {
		keys[k] = struct{}{}
	}

	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	hasDiff := false
	for _, k := range sorted {
		valA, inA := varsA[k]
		valB, inB := varsB[k]
		switch {
		case inA && !inB:
			fmt.Fprintf(cmdr.out, "- %s=%s\n", k, valA)
			hasDiff = true
		case !inA && inB:
			fmt.Fprintf(cmdr.out, "+ %s=%s\n", k, valB)
			hasDiff = true
		case valA != valB:
			fmt.Fprintf(cmdr.out, "~ %s: %s -> %s\n", k, valA, valB)
			hasDiff = true
		}
	}

	if !hasDiff {
		fmt.Fprintf(cmdr.out, "Profiles %q and %q are identical.\n", nameA, nameB)
	}
	return nil
}
