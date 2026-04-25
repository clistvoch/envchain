package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// RegisterCloneCmd registers the `clone` subcommand which deep-copies a profile
// into a new profile name, optionally prefixing every key.
func RegisterCloneCmd(app *cli.App, cmdr *Commander) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:      "clone",
		Usage:     "Clone a profile into a new profile",
		ArgsUsage: "<source> <destination>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "prefix",
				Usage: "Prefix to prepend to every key in the cloned profile",
			},
			&cli.BoolFlag{
				Name:  "overwrite",
				Usage: "Overwrite destination profile if it already exists",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
				return fmt.Errorf("usage: envchain clone <source> <destination>")
			}

			src := c.Args().Get(0)
			dst := c.Args().Get(1)
			prefix := c.String("prefix")
			overwrite := c.Bool("overwrite")

			if src == dst {
				return fmt.Errorf("source and destination profiles must differ")
			}

			if err := cmdr.CloneProfile(src, dst, prefix, overwrite); err != nil {
				return err
			}

			fmt.Printf("Cloned profile %q → %q\n", src, dst)
			if prefix != "" {
				fmt.Printf("Keys prefixed with %q\n", prefix)
			}
			return nil
		},
	})
}

// CloneProfile copies all variables from src into a new profile dst.
// If prefix is non-empty it is prepended to every key.
// Returns an error if dst already exists and overwrite is false.
func (cmdr *Commander) CloneProfile(src, dst, prefix string, overwrite bool) error {
	srcProfile, err := cmdr.store.GetProfile(src)
	if err != nil {
		return fmt.Errorf("source profile %q not found: %w", src, err)
	}

	if !overwrite {
		if _, err := cmdr.store.GetProfile(dst); err == nil {
			return fmt.Errorf("destination profile %q already exists; use --overwrite to replace", dst)
		}
	}

	dstProfile, err := cmdr.store.AddProfile(dst)
	if err != nil {
		// Profile may already exist when overwrite=true; fetch it instead.
		dstProfile, err = cmdr.store.GetProfile(dst)
		if err != nil {
			return fmt.Errorf("could not create destination profile %q: %w", dst, err)
		}
	}

	for k, v := range srcProfile.Vars() {
		key := k
		if prefix != "" {
			key = prefix + k
		}
		dstProfile.Set(key, v)
	}

	return cmdr.store.Save()
}
