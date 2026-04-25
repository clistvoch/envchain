package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// addRenameCommand registers the "rename" subcommand on the root cobra command.
// Usage: envchain rename <old-profile> <new-profile>
func (c *Commander) addRenameCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "rename <old-profile> <new-profile>",
		Short: "Rename an existing profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName := args[0]
			newName := args[1]
			return c.RenameProfile(oldName, newName)
		},
	}
	root.AddCommand(cmd)
}

// RenameProfile copies all variables from oldName into a new profile called
// newName and then removes the old profile. It returns an error when the
// source profile does not exist, the destination profile already exists, or
// either name is invalid.
func (c *Commander) RenameProfile(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return errors.New("profile names must not be empty")
	}
	if oldName == newName {
		return errors.New("old and new profile names are the same")
	}

	src, err := c.store.GetProfile(oldName)
	if err != nil {
		return fmt.Errorf("source profile %q not found: %w", oldName, err)
	}

	if _, err := c.store.GetProfile(newName); err == nil {
		return fmt.Errorf("destination profile %q already exists", newName)
	}

	if err := c.store.AddProfile(newName); err != nil {
		return fmt.Errorf("could not create profile %q: %w", newName, err)
	}

	dst, err := c.store.GetProfile(newName)
	if err != nil {
		return fmt.Errorf("could not retrieve new profile %q: %w", newName, err)
	}

	for k, v := range src.Vars() {
		dst.Set(k, v)
	}

	if err := c.store.Save(); err != nil {
		return fmt.Errorf("could not save store: %w", err)
	}

	if err := c.store.DeleteProfile(oldName); err != nil {
		return fmt.Errorf("could not remove old profile %q: %w", oldName, err)
	}

	if err := c.store.Save(); err != nil {
		return fmt.Errorf("could not save store after removing old profile: %w", err)
	}

	fmt.Fprintf(c.out, "renamed profile %q to %q\n", oldName, newName)
	return nil
}
