package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/envchain/envchain/internal/env"
	"github.com/envchain/envchain/internal/store"
)

// Commander holds the store and resolver for CLI operations.
type Commander struct {
	store    *store.Store
	resolver *env.Resolver
}

// NewCommander creates a Commander with the given store path.
func NewCommander(storePath string) (*Commander, error) {
	s, err := store.NewStore(storePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open store: %w", err)
	}
	return &Commander{
		store:    s,
		resolver: env.NewResolver(s),
	}, nil
}

// SetVar sets a key-value pair in the named profile.
func (c *Commander) SetVar(profileName, key, value string) error {
	p, err := c.store.GetProfile(profileName)
	if err != nil {
		p, err = store.NewProfile(profileName)
		if err != nil {
			return err
		}
		if err := c.store.AddProfile(p); err != nil {
			return err
		}
	}
	p.Set(key, value)
	return c.store.Save()
}

// ListVars prints all variables in the given profiles (chained).
func (c *Commander) ListVars(profiles []string) error {
	vars, err := c.resolver.Resolve(profiles)
	if err != nil {
		return err
	}
	for k, v := range vars {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
	}
	return nil
}

// DeleteVar removes a key from the named profile.
func (c *Commander) DeleteVar(profileName, key string) error {
	p, err := c.store.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("profile %q not found", profileName)
	}
	p.Delete(key)
	return c.store.Save()
}

// ListProfiles prints all profile names to stdout.
func (c *Commander) ListProfiles() error {
	profiles := c.store.Profiles()
	if len(profiles) == 0 {
		fmt.Fprintln(os.Stdout, "(no profiles)")
		return nil
	}
	fmt.Fprintln(os.Stdout, strings.Join(profiles, "\n"))
	return nil
}
