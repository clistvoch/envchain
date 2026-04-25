package cli

import (
	"fmt"
	"sort"

	"github.com/urfave/cli/v2"
)

// RegisterInspectCmd registers the inspect command which shows resolved
// environment variables for a profile chain, including the source profile
// for each variable.
func RegisterInspectCmd(app *cli.App, cmdr *Commander) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:      "inspect",
		Usage:     "Show resolved env vars and their source profiles for a chain",
		ArgsUsage: "<profile> [profile...]",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("at least one profile name is required")
			}
			profiles := c.Args().Slice()
			return cmdr.RunInspect(profiles)
		},
	})
}

// RunInspect resolves the given profile chain and prints each variable
// alongside the profile it originates from.
func (cmdr *Commander) RunInspect(profiles []string) error {
	// Build a map: varName -> sourceProfile
	source := map[string]string{}
	order := []string{}

	for _, name := range profiles {
		p, err := cmdr.store.GetProfile(name)
		if err != nil {
			return fmt.Errorf("profile %q not found", name)
		}
		for k := range p.Vars {
			if _, seen := source[k]; !seen {
				order = append(order, k)
			}
			source[k] = name
		}
	}

	sort.Strings(order)

	for _, k := range order {
		profileName := source[k]
		p, _ := cmdr.store.GetProfile(profileName)
		fmt.Printf("%-30s %-20s %s\n", k, "["+profileName+"]", p.Vars[k])
	}
	return nil
}
