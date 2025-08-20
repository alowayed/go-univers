package maven

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alowayed/go-univers/pkg/vers"
)

type VersionRange struct {
	original        string
	constraints     []constraint              // For native Maven ranges
	versConstraints []vers.Constraint        // For VERS ranges
}

type constraint struct {
	version   *Version
	inclusive bool
	isLower   bool // true for lower bound, false for upper bound
}

func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	if rangeStr == "" {
		return nil, fmt.Errorf("range string cannot be empty")
	}

	// Trim whitespace
	trimmed := strings.TrimSpace(rangeStr)
	if trimmed == "" {
		return nil, fmt.Errorf("range string cannot be empty or only whitespace")
	}

	constraints, err := parseVersionRange(trimmed, e)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		original:    rangeStr,
		constraints: constraints,
	}, nil
}

func (vr *VersionRange) Contains(version *Version) bool {
	// If this is a VERS range, use the VERS algorithm
	if vr.versConstraints != nil {
		return vr.containsVers(version)
	}

	// Otherwise use the traditional Maven algorithm
	if len(vr.constraints) == 0 {
		return false
	}

	// All constraints must be satisfied
	for _, constraint := range vr.constraints {
		if !satisfiesConstraint(version, constraint) {
			return false
		}
	}
	return true
}

// containsVers uses the VERS algorithm for checking version containment
func (vr *VersionRange) containsVers(version *Version) bool {
	// Create a parse function for Maven versions
	parseVersion := func(versionStr string) (vers.VersionComparator, error) {
		e := &Ecosystem{}
		v, err := e.NewVersion(versionStr)
		if err != nil {
			return nil, err
		}
		return &versionComparatorWrapper{v}, nil
	}

	// Use the VERS containment algorithm
	result, err := vers.ContainsVersion(&versionComparatorWrapper{version}, vr.versConstraints, parseVersion)
	if err != nil {
		return false // If there's an error, assume not contained
	}
	return result
}

// versionComparatorWrapper wraps Maven Version to implement vers.VersionComparator
type versionComparatorWrapper struct {
	*Version
}

func (w *versionComparatorWrapper) Compare(other vers.VersionComparator) int {
	if otherWrapper, ok := other.(*versionComparatorWrapper); ok {
		return w.Version.Compare(otherWrapper.Version)
	}
	// Fallback to string comparison if types don't match
	return strings.Compare(w.String(), other.String())
}

func (vr *VersionRange) String() string {
	return vr.original
}

func parseVersionRange(rangeStr string, e *Ecosystem) ([]constraint, error) {
	var constraints []constraint

	// Check if it's a bracket range: [1.0], [1.0,2.0], (1.0,2.0), etc.
	bracketRegex := regexp.MustCompile(`^[\[\(]([^,\]\)]*)(,([^,\]\)]*))?[\]\)]$`)
	matches := bracketRegex.FindStringSubmatch(rangeStr)

	if matches != nil {
		// This is a bracket range
		lowerBracket := rangeStr[0]
		upperBracket := rangeStr[len(rangeStr)-1]
		lowerInclusive := lowerBracket == '['
		upperInclusive := upperBracket == ']'

		lowerVersionStr := strings.TrimSpace(matches[1])
		upperVersionStr := ""
		if len(matches) > 3 && matches[3] != "" {
			upperVersionStr = strings.TrimSpace(matches[3])
		}

		// Check for empty exact version []
		if matches[2] == "" && lowerVersionStr == "" {
			return nil, fmt.Errorf("empty version in exact range")
		}

		// Handle exact version [1.0] (no comma in the match)
		if matches[2] == "" && lowerVersionStr != "" {
			version, err := e.NewVersion(lowerVersionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid version in range: %v", err)
			}
			// Exact version means both upper and lower bounds are the same
			constraints = append(constraints, constraint{
				version:   version,
				inclusive: true,
				isLower:   true,
			})
			constraints = append(constraints, constraint{
				version:   version,
				inclusive: true,
				isLower:   false,
			})
			return constraints, nil
		}

		// Handle lower bound
		if lowerVersionStr != "" {
			version, err := e.NewVersion(lowerVersionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid lower bound version: %v", err)
			}
			constraints = append(constraints, constraint{
				version:   version,
				inclusive: lowerInclusive,
				isLower:   true,
			})
		}

		// Handle upper bound
		if upperVersionStr != "" {
			version, err := e.NewVersion(upperVersionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid upper bound version: %v", err)
			}
			constraints = append(constraints, constraint{
				version:   version,
				inclusive: upperInclusive,
				isLower:   false,
			})
		}

		// Validate that we have at least one constraint
		if len(constraints) == 0 {
			return nil, fmt.Errorf("invalid range format")
		}

		return constraints, nil
	}

	// Check for malformed brackets (missing closing bracket)
	if strings.Contains(rangeStr, "[") || strings.Contains(rangeStr, "(") ||
		strings.Contains(rangeStr, "]") || strings.Contains(rangeStr, ")") {
		return nil, fmt.Errorf("malformed bracket range")
	}

	// If not a bracket range, treat as simple version requirement
	version, err := e.NewVersion(rangeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version: %v", err)
	}

	// Simple version is treated as exact match
	constraints = append(constraints, constraint{
		version:   version,
		inclusive: true,
		isLower:   true,
	})
	constraints = append(constraints, constraint{
		version:   version,
		inclusive: true,
		isLower:   false,
	})

	return constraints, nil
}

func satisfiesConstraint(version *Version, constraint constraint) bool {
	cmp := version.Compare(constraint.version)

	if constraint.isLower {
		// Lower bound: version >= constraint.version (if inclusive) or version > constraint.version (if exclusive)
		if constraint.inclusive {
			return cmp >= 0
		} else {
			return cmp > 0
		}
	} else {
		// Upper bound: version <= constraint.version (if inclusive) or version < constraint.version (if exclusive)
		if constraint.inclusive {
			return cmp <= 0
		} else {
			return cmp < 0
		}
	}
}
