package spaces

import (
	"errors"
	"regexp"
	"strings"
)

var (
	
	slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)


func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	if len(name) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}
	return nil
}


func ValidateSlug(slug string) error {
	slug = strings.TrimSpace(slug)
	if len(slug) < 2 {
		return errors.New("slug must be at least 2 characters long")
	}
	if len(slug) > 100 {
		return errors.New("slug cannot exceed 100 characters")
	}
	if !slugPattern.MatchString(slug) {
		return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}
	return nil
}
