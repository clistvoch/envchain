package store

import (
	"errors"
	"regexp"
)

var validName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// Profile holds a named set of environment variables.
type Profile struct {
	name string
	vars map[string]string
}

// NewProfile creates a new empty Profile with the given name.
// Returns an error if the name contains invalid characters.
func NewProfile(name string) (*Profile, error) {
	if !validName.MatchString(name) {
		return nil, errors.New("profile name must match [a-zA-Z0-9_-]+")
	}
	return &Profile{name: name, vars: make(map[string]string)}, nil
}

// Name returns the profile's name.
func (p *Profile) Name() string { return p.name }

// Set stores a key/value pair in the profile.
func (p *Profile) Set(key, value string) { p.vars[key] = value }

// Get retrieves a value by key. Returns ("", false) if not found.
func (p *Profile) Get(key string) (string, bool) {
	v, ok := p.vars[key]
	return v, ok
}

// Delete removes a key from the profile.
func (p *Profile) Delete(key string) { delete(p.vars, key) }

// Vars returns a shallow copy of all variables in the profile.
func (p *Profile) Vars() map[string]string {
	copy := make(map[string]string, len(p.vars))
	for k, v := range p.vars {
		copy[k] = v
	}
	return copy
}

// Merge copies all variables from src into p, overwriting existing keys.
func (p *Profile) Merge(src *Profile) {
	for k, v := range src.vars {
		p.vars[k] = v
	}
}
