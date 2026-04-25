package cli

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

// RegisterSnapshotCmd registers the snapshot command which saves a named
// point-in-time copy of a profile for later restoration.
func RegisterSnapshotCmd(app *cli.App, cmdr *Commander) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:      "snapshot",
		Usage:     "Save or restore a named snapshot of a profile",
		ArgsUsage: "<profile>",
		Subcommands: []*cli.Command{
			{
				Name:      "save",
				Usage:     "Save current profile state as a snapshot",
				ArgsUsage: "<profile> [snapshot-name]",
				Action: func(c *cli.Context) error {
					profileName := c.Args().Get(0)
					if profileName == "" {
						return fmt.Errorf("profile name is required")
					}
					snapshotName := c.Args().Get(1)
					if snapshotName == "" {
						snapshotName = time.Now().Format("20060102-150405")
					}
					destName := fmt.Sprintf("%s.snap.%s", profileName, snapshotName)
					if err := cmdr.CopyProfile(profileName, destName, false); err != nil {
						return fmt.Errorf("snapshot save failed: %w", err)
					}
					fmt.Fprintf(c.App.Writer, "Snapshot saved: %s\n", destName)
					return nil
				},
			},
			{
				Name:      "restore",
				Usage:     "Restore a profile from a snapshot",
				ArgsUsage: "<profile> <snapshot-name>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "overwrite",
						Usage: "Overwrite the target profile if it already exists",
					},
				},
				Action: func(c *cli.Context) error {
					profileName := c.Args().Get(0)
					snapshotName := c.Args().Get(1)
					if profileName == "" || snapshotName == "" {
						return fmt.Errorf("profile and snapshot name are required")
					}
					srcName := fmt.Sprintf("%s.snap.%s", profileName, snapshotName)
					overwrite := c.Bool("overwrite")
					if err := cmdr.CopyProfile(srcName, profileName, overwrite); err != nil {
						return fmt.Errorf("snapshot restore failed: %w", err)
					}
					fmt.Fprintf(c.App.Writer, "Profile %q restored from snapshot %q\n", profileName, snapshotName)
					return nil
				},
			},
			{
				Name:      "list",
				Usage:     "List all snapshots for a profile",
				ArgsUsage: "<profile>",
				Action: func(c *cli.Context) error {
					profileName := c.Args().Get(0)
					if profileName == "" {
						return fmt.Errorf("profile name is required")
					}
					snaps, err := cmdr.ListSnapshots(profileName)
					if err != nil {
						return err
					}
					if len(snaps) == 0 {
						fmt.Fprintf(c.App.Writer, "No snapshots found for profile %q\n", profileName)
						return nil
					}
					for _, s := range snaps {
						fmt.Fprintln(c.App.Writer, s)
					}
					return nil
				},
			},
		},
	})
}
