package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
)

// RegisterTagCmd registers the tag subcommands for labeling profiles.
func RegisterTagCmd(app *cli.App, cmdr *Commander) {
	cmd := &cli.Command{
		Name:  "tag",
		Usage: "Manage tags on profiles",
		Subcommands: []*cli.Command{
			{
				Name:      "add",
				Usage:     "Add a tag to a profile",
				ArgsUsage: "<profile> <tag>",
				Action: func(c *cli.Context) error {
					if c.NArg() < 2 {
						return fmt.Errorf("usage: tag add <profile> <tag>")
					}
					profile := c.Args().Get(0)
					tag := c.Args().Get(1)
					return cmdr.AddTag(profile, tag)
				},
			},
			{
				Name:      "remove",
				Usage:     "Remove a tag from a profile",
				ArgsUsage: "<profile> <tag>",
				Action: func(c *cli.Context) error {
					if c.NArg() < 2 {
						return fmt.Errorf("usage: tag remove <profile> <tag>")
					}
					profile := c.Args().Get(0)
					tag := c.Args().Get(1)
					return cmdr.RemoveTag(profile, tag)
				},
			},
			{
				Name:      "list",
				Usage:     "List profiles matching a tag",
				ArgsUsage: "<tag>",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return fmt.Errorf("usage: tag list <tag>")
					}
					tag := c.Args().Get(0)
					profiles, err := cmdr.ListByTag(tag)
					if err != nil {
						return err
					}
					if len(profiles) == 0 {
						fmt.Printf("No profiles tagged %q\n", tag)
						return nil
					}
					fmt.Printf("Profiles tagged %q:\n", tag)
					for _, p := range profiles {
						fmt.Printf("  %s [%s]\n", p, strings.Join(cmdr.GetTags(p), ", "))
					}
					return nil
				},
			},
		},
	}
	app.Commands = append(app.Commands, cmd)
}
