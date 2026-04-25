// Package cli provides command-line interface commands for envchain.
//
// # Lint Command
//
// The lint command validates one or more profiles for common issues:
//
//   - Invalid environment variable names (must match [A-Za-z_][A-Za-z0-9_]* and
//     must not start with double underscore)
//   - Shadowed keys: a variable defined in an earlier profile that is overridden
//     by a later profile in the chain
//   - Missing profiles: referenced profile names that do not exist in the store
//
// Usage:
//
//	envchain lint <profile> [profile...]
//	envchain lint --strict <profile> [profile...]
//
// Flags:
//
//	--strict    Exit with a non-zero status code when any warnings are found.
//	            Useful in CI pipelines to enforce clean profile definitions.
//
// Examples:
//
//	# Lint a single profile
//	envchain lint production
//
//	# Lint a chain and fail on warnings
//	envchain lint --strict base staging production
package cli
