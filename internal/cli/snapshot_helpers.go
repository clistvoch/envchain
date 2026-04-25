package cli

import (
	"fmt"
	"strings"
)

// snapshotPrefix returns the prefix used to identify snapshot profiles.
func snapshotPrefix(profileName string) string {
	return fmt.Sprintf("%s.snap.", profileName)
}

// ListSnapshots returns all snapshot names associated with the given profile.
func (c *Commander) ListSnapshots(profileName string) ([]string, error) {
	profiles, err := c.store.ListProfiles()
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}
	prefix := snapshotPrefix(profileName)
	var snapshots []string
	for _, p := range profiles {
		if strings.HasPrefix(p, prefix) {
			snapshotTag := strings.TrimPrefix(p, prefix)
			snapshots = append(snapshots, snapshotTag)
		}
	}
	return snapshots, nil
}

// DeleteSnapshot removes a specific snapshot for the given profile.
func (c *Commander) DeleteSnapshot(profileName, snapshotName string) error {
	srcName := fmt.Sprintf("%s.snap.%s", profileName, snapshotName)
	profiles, err := c.store.ListProfiles()
	if err != nil {
		return fmt.Errorf("failed to list profiles: %w", err)
	}
	for _, p := range profiles {
		if p == srcName {
			return c.store.DeleteProfile(srcName)
		}
	}
	return fmt.Errorf("snapshot %q not found for profile %q", snapshotName, profileName)
}
