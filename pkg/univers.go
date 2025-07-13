package univers

// This file contains documentation interfaces that describe the common patterns
// used across ecosystem packages. These are not meant to be implemented directly,
// but serve as documentation for the expected API surface of each ecosystem.
//
// Each ecosystem (npm, pypi, gomod) has its own concrete types that follow these
// patterns but use ecosystem-specific parameter types for type safety.

// Version describes the common interface pattern for version types.
// Actual implementations use concrete types for parameters (e.g., *npm.Version).
type Version interface {
	// String returns the string representation of the version
	String() string
	
	// Normalize returns the normalized form of the version
	Normalize() string
	
	// Compare compares this version with another version of the same ecosystem
	// Returns -1 if this < other, 0 if this == other, 1 if this > other
	// Actual signature: Compare(other *EcosystemVersion) int
	Compare(other interface{}) int
}

// VersionRange describes the common interface pattern for version range types.
// Actual implementations use concrete types for parameters (e.g., *npm.Version).
type VersionRange interface {
	// String returns the string representation of the range
	String() string
	
	// Contains checks if a version is within this range
	// Actual signature: Contains(version *EcosystemVersion) bool
	Contains(version interface{}) bool
}