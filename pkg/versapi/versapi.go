// Package versapi provides stateless API functions for VERS operations.
// This package acts as a bridge between the vers parsing utilities and ecosystem implementations.
package versapi

import (
	"fmt"

	"github.com/alowayed/go-univers/pkg/ecosystem/maven"
	"github.com/alowayed/go-univers/pkg/vers"
)

// Contains checks if a version satisfies a VERS range using the stateless API.
// Example: Contains("vers:maven/>=1.0.0|<=2.0.0", "1.5.0") returns true.
func Contains(versRange, version string) (bool, error) {
	// Parse the VERS range to get ecosystem and constraints
	ecosystem, _, err := vers.ParseVersString(versRange)
	if err != nil {
		return false, fmt.Errorf("invalid VERS range: %w", err)
	}

	// Route to appropriate ecosystem implementation
	switch ecosystem {
	case "maven":
		return containsMaven(versRange, version)
	default:
		return false, fmt.Errorf("ecosystem '%s' not yet implemented for stateless API", ecosystem)
	}
}

// Compare compares two versions within a VERS range context using the stateless API.
// Returns -1 if version1 < version2, 0 if equal, 1 if version1 > version2.
func Compare(versRange, version1, version2 string) (int, error) {
	// Parse the VERS range to get ecosystem
	ecosystem, _, err := vers.ParseVersString(versRange)
	if err != nil {
		return 0, fmt.Errorf("invalid VERS range: %w", err)
	}

	// Route to appropriate ecosystem implementation
	switch ecosystem {
	case "maven":
		return compareMaven(version1, version2)
	default:
		return 0, fmt.Errorf("ecosystem '%s' not yet implemented for stateless API", ecosystem)
	}
}

// containsMaven implements Contains for Maven ecosystem
func containsMaven(versRange, version string) (bool, error) {
	// Parse the VERS range using Maven's ParseVers
	vr, err := maven.ParseVers(versRange)
	if err != nil {
		return false, fmt.Errorf("failed to parse Maven VERS range: %w", err)
	}

	// Parse the version using Maven's NewVersion
	e := &maven.Ecosystem{}
	v, err := e.NewVersion(version)
	if err != nil {
		return false, fmt.Errorf("failed to parse Maven version '%s': %w", version, err)
	}

	// Check if the version is contained in the range
	return vr.Contains(v), nil
}

// compareMaven implements Compare for Maven ecosystem  
func compareMaven(version1, version2 string) (int, error) {
	// Parse both versions using Maven's NewVersion
	e := &maven.Ecosystem{}
	
	v1, err := e.NewVersion(version1)
	if err != nil {
		return 0, fmt.Errorf("failed to parse Maven version '%s': %w", version1, err)
	}

	v2, err := e.NewVersion(version2)
	if err != nil {
		return 0, fmt.Errorf("failed to parse Maven version '%s': %w", version2, err)
	}

	// Compare the versions
	return v1.Compare(v2), nil
}