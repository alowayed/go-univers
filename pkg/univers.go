package univers

// Version represents a software version in any ecosystem
type Version interface {
	// String returns the string representation of the version
	String() string
	
	// IsValid checks if the version string is valid for this ecosystem
	IsValid() bool
	
	// Normalize returns the normalized form of the version
	Normalize() string
	
	// Compare compares this version with another version
	// Returns -1 if this < other, 0 if this == other, 1 if this > other
	Compare(other Version) int
	
	// Satisfies checks if this version satisfies the given constraint
	Satisfies(constraint VersionConstraint) bool
}

// VersionConstraint represents a version constraint (e.g., >=1.0.0)
type VersionConstraint interface {
	// String returns the string representation of the constraint
	String() string
	
	// Operator returns the constraint operator (=, !=, <, <=, >, >=, *)
	Operator() string
	
	// Version returns the version part of the constraint
	Version() string
	
	// Matches checks if the given version matches this constraint
	Matches(version Version) bool
}

// VersionRange represents a collection of version constraints
type VersionRange interface {
	// String returns the string representation of the range
	String() string
	
	// Contains checks if a version is within this range
	Contains(version Version) bool
	
	// Constraints returns all constraints in this range
	Constraints() []VersionConstraint
	
	// IsEmpty returns true if the range contains no valid versions
	IsEmpty() bool
}