package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Store manages persistence of profiles on disk.
type Store struct {
	path     string
	profiles map[string]*Profile
}

type storeData struct {
	Profiles map[string]*Profile `json:"profiles"`
}

// NewStore creates a Store backed by the given file path.
func NewStore(path string) (*Store, error) {
	s := &Store{
		path:     path,
		profiles: make(map[string]*Profile),
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	var sd storeData
	if err := json.Unmarshal(data, &sd); err != nil {
		return err
	}
	if sd.Profiles != nil {
		s.profiles = sd.Profiles
	}
	return nil
}

// Save persists all profiles to disk.
func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	sd := storeData{Profiles: s.profiles}
	data, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

// AddProfile adds a new profile to the store.
func (s *Store) AddProfile(p *Profile) error {
	if _, exists := s.profiles[p.Name]; exists {
		return ErrProfileAlreadyExists
	}
	s.profiles[p.Name] = p
	return nil
}

// GetProfile retrieves a profile by name.
func (s *Store) GetProfile(name string) (*Profile, error) {
	p, ok := s.profiles[name]
	if !ok {
		return nil, ErrProfileNotFound
	}
	return p, nil
}

// DeleteProfile removes a profile by name.
func (s *Store) DeleteProfile(name string) error {
	if _, ok := s.profiles[name]; !ok {
		return ErrProfileNotFound
	}
	delete(s.profiles, name)
	return nil
}

// ListProfiles returns all profile names.
func (s *Store) ListProfiles() []string {
	names := make([]string, 0, len(s.profiles))
	for name := range s.profiles {
		names = append(names, name)
	}
	return names
}
