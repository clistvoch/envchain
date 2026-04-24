// Package cli provides the Commander type which bridges the CLI layer
// with the underlying store and environment resolver.
//
// Commander exposes high-level operations such as setting and deleting
// variables within named profiles, listing profiles, and resolving
// chained profile variables for process injection.
//
// Typical usage:
//
//	cmd, err := cli.NewCommander("/path/to/store.json")
//	if err != nil { ... }
//
//	// Set a variable in a profile (creates the profile if absent)
//	cmd.SetVar("dev", "DATABASE_URL", "postgres://localhost/dev")
//
//	// List resolved variables from a chain of profiles
//	cmd.ListVars([]string{"base", "dev"})
package cli
