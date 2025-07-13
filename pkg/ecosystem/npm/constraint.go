package npm

import (
	"strconv"
	"strings"
)

// constraint represents a single NPM version constraint
type constraint struct {
	operator string
	version  string
}

// String returns the string representation of the constraint
func (c *constraint) String() string {
	if c.operator == "=" {
		return c.version
	}
	return c.operator + c.version
}


// matches checks if the given version matches this constraint
func (c *constraint) matches(version *Version) bool {
	if c.operator == "*" {
		return true
	}

	constraintVersion, err := NewVersion(c.version)
	if err != nil {
		return false
	}

	comparison := version.Compare(constraintVersion)

	switch c.operator {
	case "=":
		return comparison == 0
	case "!=":
		return comparison != 0
	case "<":
		return comparison < 0
	case "<=":
		return comparison <= 0
	case ">":
		return comparison > 0
	case ">=":
		return comparison >= 0
	default:
		return false
	}
}

// Helper functions

func comparePrerelease(a, b string) int {
	// No prerelease has higher precedence than prerelease
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return 1
	}
	if b == "" {
		return -1
	}

	// Split by dots and compare each part
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		var aPart, bPart string
		if i < len(aParts) {
			aPart = aParts[i]
		}
		if i < len(bParts) {
			bPart = bParts[i]
		}

		// Missing part has lower precedence
		if aPart == "" && bPart != "" {
			return -1
		}
		if aPart != "" && bPart == "" {
			return 1
		}

		// Try to parse as numbers
		aNum, aIsNum := parseNum(aPart)
		bNum, bIsNum := parseNum(bPart)

		if aIsNum && bIsNum {
			if aNum != bNum {
				return compareInt(aNum, bNum)
			}
		} else if aIsNum {
			return -1 // Numeric identifiers have lower precedence
		} else if bIsNum {
			return 1
		} else {
			// Both are strings, compare lexically
			if aPart != bPart {
				return strings.Compare(aPart, bPart)
			}
		}
	}

	return 0
}

// compareInt returns -1 if a < b, 0 if a == b, 1 if a > b
func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// parseNum returns the integer value and true if s is a valid number, otherwise 0 and false
func parseNum(s string) (int, bool) {
	if num, err := strconv.Atoi(s); err == nil {
		return num, true
	}
	return 0, false
}