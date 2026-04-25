// Package cli — inspect command
//
// The inspect command provides a detailed view of the resolved environment
// variable chain for one or more profiles. Unlike `export`, which simply
// emits shell-ready output, `inspect` annotates each variable with the
// profile it originates from.
//
// Usage:
//
//	envchain inspect <profile> [profile...]
//
// Example:
//
//	envchain inspect base production
//
// Output format:
//
//	VAR_NAME                       [source-profile]     value
//
// When multiple profiles are specified, the first profile that defines a
// variable is considered the authoritative source (left-to-right precedence).
// This mirrors the resolution order used by the `run` and `export` commands.
package cli
