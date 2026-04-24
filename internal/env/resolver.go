package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/envchain/envchain/internal/store"
)

// Resolver resolves environment variables from a profile chain.
type Resolver struct {
	store *store.Store
}

// NewResolver creates a new Resolver backed by the given store.
func NewResolver(s *store.Store) *Resolver {
	return &Resolver{store: s}
}

// Resolve returns a merged map of environment variables for the given profile
// names, in order. Later profiles override earlier ones.
func (r *Resolver) Resolve(profileNames []string) (map[string]string, error) {
	result := make(map[string]string)

	for _, name := range profileNames {
		p, err := r.store.GetProfile(name)
		if err != nil {
			return nil, fmt.Errorf("resolver: profile %q not found: %w", name, err)
		}
		for k, v := range p.Vars() {
			result[k] = v
		}
	}

	return result, nil
}

// Environ returns the resolved variables as a slice of "KEY=VALUE" strings,
// suitable for use with os/exec.Cmd.Env.
func (r *Resolver) Environ(profileNames []string) ([]string, error) {
	vars, err := r.Resolve(profileNames)
	if err != nil {
		return nil, err
	}

	env := make([]string, 0, len(vars))
	for k, v := range vars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env, nil
}

// ApplyToProcess sets the resolved variables in the current process environment.
func (r *Resolver) ApplyToProcess(profileNames []string) error {
	vars, err := r.Resolve(profileNames)
	if err != nil {
		return err
	}

	for k, v := range vars {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("resolver: failed to set %s: %w", k, err)
		}
	}
	return nil
}

// ExpandVars replaces ${VAR} or $VAR references in s using the resolved vars.
func ExpandVars(vars map[string]string, s string) string {
	return os.Expand(s, func(key string) string {
		if v, ok := vars[key]; ok {
			return v
		}
		return strings.TrimSpace(os.Getenv(key))
	})
}
