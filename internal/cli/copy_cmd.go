package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newCopyCmd returns a command that copies all variables from one profile to another.
func (c *Commander) newCopyCmd() *cobra.Command {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "copy <src-profile> <dst-profile>",
		Short: "Copy all variables from one profile to another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]

			srcProfile, err := c.store.GetProfile(src)
			if err != nil {
				return fmt.Errorf("source profile %q not found: %w", src, err)
			}

			dstProfile, err := c.store.GetProfile(dst)
			if err != nil {
				// Destination does not exist yet; create it via AddProfile.
				if addErr := c.store.AddProfile(dst); addErr != nil {
					return fmt.Errorf("could not create destination profile %q: %w", dst, addErr)
				}
				dstProfile, err = c.store.GetProfile(dst)
				if err != nil {
					return err
				}
			}

			vars := srcProfile.Vars()
			copied := 0
			skipped := 0
			for k, v := range vars {
				if _, exists := dstProfile.Get(k); exists && !overwrite {
					skipped++
					continue
				}
				dstProfile.Set(k, v)
				copied++
			}

			if err := c.store.Save(); err != nil {
				return fmt.Errorf("failed to save store: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Copied %d variable(s) from %q to %q", copied, src, dst)
			if skipped > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), " (%d skipped, use --overwrite to replace)", skipped)
			}
			fmt.Fprintln(cmd.OutOrStdout())
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing variables in the destination profile")
	return cmd
}
