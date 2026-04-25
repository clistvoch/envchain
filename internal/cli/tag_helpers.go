package cli

import (
	"fmt"
	"sort"
	"strings"
)

const tagMetaKey = "__envchain_tags__"

// AddTag adds a tag to the named profile by storing it as a special meta key.
func (c *Commander) AddTag(profileName, tag string) error {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return fmt.Errorf("tag must not be empty")
	}
	p, err := c.store.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("profile %q not found", profileName)
	}
	existing := parseTags(p.Get(tagMetaKey))
	for _, t := range existing {
		if t == tag {
			return fmt.Errorf("tag %q already exists on profile %q", tag, profileName)
		}
	}
	existing = append(existing, tag)
	sort.Strings(existing)
	p.Set(tagMetaKey, strings.Join(existing, ","))
	return c.store.Save()
}

// RemoveTag removes a tag from the named profile.
func (c *Commander) RemoveTag(profileName, tag string) error {
	p, err := c.store.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("profile %q not found", profileName)
	}
	existing := parseTags(p.Get(tagMetaKey))
	updated := existing[:0]
	found := false
	for _, t := range existing {
		if t == tag {
			found = true
			continue
		}
		updated = append(updated, t)
	}
	if !found {
		return fmt.Errorf("tag %q not found on profile %q", tag, profileName)
	}
	if len(updated) == 0 {
		p.Delete(tagMetaKey)
	} else {
		p.Set(tagMetaKey, strings.Join(updated, ","))
	}
	return c.store.Save()
}

// ListByTag returns all profile names that carry the given tag.
func (c *Commander) ListByTag(tag string) ([]string, error) {
	names := c.store.ProfileNames()
	var matched []string
	for _, name := range names {
		if p, err := c.store.GetProfile(name); err == nil {
			for _, t := range parseTags(p.Get(tagMetaKey)) {
				if t == tag {
					matched = append(matched, name)
					break
				}
			}
		}
	}
	sort.Strings(matched)
	return matched, nil
}

// GetTags returns the tags for a profile (best-effort, empty on error).
func (c *Commander) GetTags(profileName string) []string {
	p, err := c.store.GetProfile(profileName)
	if err != nil {
		return nil
	}
	return parseTags(p.Get(tagMetaKey))
}

func parseTags(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := parts[:0]
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
