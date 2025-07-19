// Package univers provides interfaces for package ecosystems, versions, and version ranges.
package univers

// Version represents a version within a specific ecosystem.
type Version[T any] interface {
	// Compare compares this version with another version of the same type.
	// Returns -1 if this < other, 0 if this == other, 1 if this > other.
	Compare(other T) int

	// Returns the original string representation of the version.
	String() string
}

// VersionRange represents a version range within a specific ecosystem.
type VersionRange[V Version[V]] interface {
	// Contains checks if a version is within this range.
	Contains(version V) bool

	// Returns the original string representation of the version range.
	String() string
}

// Ecosystem represents a package ecosystem that can create versions and version ranges.
type Ecosystem[V Version[V], VR VersionRange[V]] interface {
	// Name returns the name of the ecosystem.
	Name() string

	// NewVersion creates a new version instance from a string.
	NewVersion(s string) (V, error)

	// NewVersionRange creates a new version range instance from a string.
	NewVersionRange(s string) (VR, error)
}
