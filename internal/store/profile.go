package store

import (
	"errors"
	"regexp"
)

var profileNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// ErrInvalidProfileName is returned when a profile name contains invalid characters.
var ErrInvalidProfileName = errors.New("profile name must contain only alphanumeric characters, hyphens, or underscores")

// ErrProfileNotFound is returned when a requested profile does not exist.
var ErrProfileNotFound = errors.New("profile not found")

// ErrProfileAlreadyExists is returned when trying to create a duplicate profile.
var ErrProfileAlreadyExists = errors.New("profile already exists")

// Profile represents a named set of environment variables.
type Profile struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// NewProfile creates a new Profile with the given name.
func NewProfile(name string) (*Profile, error) {
	if !profileNameRegex.MatchString(name) {
		return nil, ErrInvalidProfileName
	}
	return &Profile{
		Name: name,
		Vars: make(map[string]string),
	}, nil
}

// Set adds or updates an environment variable in the profile.
func (p *Profile) Set(key, value string) {
	p.Vars[key] = value
}

// Get retrieves an environment variable from the profile.
func (p *Profile) Get(key string) (string, bool) {
	v, ok := p.Vars[key]
	return v, ok
}

// Delete removes an environment variable from the profile.
func (p *Profile) Delete(key string) {
	delete(p.Vars, key)
}

// Merge combines another profile's variables into this profile.
// Variables from other take precedence over existing ones.
func (p *Profile) Merge(other *Profile) {
	for k, v := range other.Vars {
		p.Vars[k] = v
	}
}
