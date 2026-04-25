// Package cli provides command-line interface commands for envchain.
//
// # Tag Command
//
// The tag command allows you to label profiles with arbitrary string tags,
// making it easier to group, filter, and discover related profiles.
//
// Tags are stored as a special internal metadata key within each profile and
// are excluded from environment variable resolution and export.
//
// Usage:
//
//	envchain tag add <profile> <tag>     # Attach a tag to a profile
//	envchain tag remove <profile> <tag>  # Detach a tag from a profile
//	envchain tag list <tag>              # List all profiles carrying a tag
//
// Examples:
//
//	# Mark profiles with environment tier tags
//	envchain tag add prod-aws production
//	envchain tag add staging-aws staging
//
//	# Find all production profiles
//	envchain tag list production
//
// Tags are sorted alphabetically and deduplicated automatically.
package cli
