// Package maven provides functionality for working with Maven package versions.
package maven

import (
	"fmt"

	"github.com/alowayed/go-univers/pkg/vers"
)

const (
	Name = "maven"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}

// ParseVers parses a VERS string into a Maven VersionRange.
// Example: "vers:maven/>=1.0.0|<=2.0.0" creates a range from 1.0.0 to 2.0.0.
func ParseVers(versString string) (*VersionRange, error) {
	ecosystem, constraints, err := vers.ParseVersString(versString)
	if err != nil {
		return nil, fmt.Errorf("invalid VERS string: %w", err)
	}

	if ecosystem != Name {
		return nil, fmt.Errorf("VERS string is for ecosystem '%s', expected '%s'", ecosystem, Name)
	}

	// Parse constraints into VERS constraint objects
	versConstraints, err := vers.ParseVersConstraints(constraints)
	if err != nil {
		return nil, fmt.Errorf("invalid VERS constraints: %w", err)
	}

	return &VersionRange{
		original:       versString,
		versConstraints: versConstraints, // Store VERS constraints for proper algorithm
	}, nil
}

