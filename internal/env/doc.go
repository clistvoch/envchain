// Package env provides utilities for resolving and applying environment
// variable profiles managed by envchain.
//
// The central type is [Resolver], which accepts a [store.Store] and can
// merge one or more named profiles into a flat key-value map, export them
// as "KEY=VALUE" strings for subprocess execution, or apply them directly
// to the current process via os.Setenv.
//
// Profile chaining is supported by passing multiple profile names to
// Resolve or Environ; variables from later profiles override those from
// earlier ones, enabling composable configuration layering.
//
// Example usage:
//
//	s, _ := store.NewStore("/home/user/.envchain.json")
//	r := env.NewResolver(s)
//	envVars, _ := r.Environ([]string{"base", "production"})
//	cmd := exec.Command("make", "deploy")
//	cmd.Env = append(os.Environ(), envVars...)
package env
