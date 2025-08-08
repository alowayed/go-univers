package rpm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// versionPattern matches RPM version strings
// Format: [epoch:]version[-release]
var versionPattern = regexp.MustCompile(`^(?:(\d+):)?(.+)$`)

// Version represents an RPM package version
type Version struct {
	epoch    int    // optional epoch (defaults to 0)
	version  string // version part (required)
	release  string // optional release part
	original string // original version string
}

// NewVersion creates a new RPM version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(version)

	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Parse using regex to extract epoch
	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid RPM version format: %s", original)
	}

	epochStr := matches[1]
	versionReleasePart := matches[2]

	// Split version and release on the last hyphen (if any)
	var versionPart, releasePart string
	if lastHyphen := strings.LastIndex(versionReleasePart, "-"); lastHyphen != -1 {
		versionPart = versionReleasePart[:lastHyphen]
		releasePart = versionReleasePart[lastHyphen+1:]
	} else {
		versionPart = versionReleasePart
		releasePart = ""
	}

	// Parse epoch (default to 0)
	epoch := 0
	if epochStr != "" {
		var err error
		epoch, err = strconv.Atoi(epochStr)
		if err != nil {
			return nil, fmt.Errorf("invalid epoch in version %s: %v", original, err)
		}
		if epoch < 0 {
			return nil, fmt.Errorf("epoch cannot be negative in version %s", original)
		}
	}

	// Validate version part (cannot be empty)
	if versionPart == "" {
		return nil, fmt.Errorf("version part cannot be empty in version %s", original)
	}

	// Validate characters in version and release parts
	if err := validateRPMVersionString(versionPart, "version"); err != nil {
		return nil, fmt.Errorf("invalid version part in %s: %v", original, err)
	}

	if releasePart != "" {
		if err := validateRPMVersionString(releasePart, "release"); err != nil {
			return nil, fmt.Errorf("invalid release part in %s: %v", original, err)
		}
	}

	return &Version{
		epoch:    epoch,
		version:  versionPart,
		release:  releasePart,
		original: original,
	}, nil
}

// validateRPMVersionString validates that a version string contains only allowed characters
// RPM allows alphanumerics and . + - ~ ^ (and : for epochs, but that's handled separately)
func validateRPMVersionString(s, part string) error {
	for _, r := range s {
		if !isValidRPMVersionChar(r) {
			return fmt.Errorf("invalid character %q in %s", r, part)
		}
	}
	return nil
}

// isValidRPMVersionChar checks if a character is valid in an RPM version string
// Based on RPM specification: alphanumerics and the characters . + - ~ ^ _
func isValidRPMVersionChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '+' || r == '-' || r == '~' || r == '^' || r == '_'
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another RPM version
// Follows RPM's version comparison algorithm:
// 1. Compare epochs first (higher epoch wins)
// 2. Compare version parts using RPM's lexicographic/numeric rules
// 3. Compare release parts using the same rules
func (v *Version) Compare(other *Version) int {
	// 1. Compare epochs first
	if v.epoch != other.epoch {
		if v.epoch < other.epoch {
			return -1
		}
		return 1
	}

	// 2. Compare version parts
	versionCmp := compareRPMVersionString(v.version, other.version)
	if versionCmp != 0 {
		return versionCmp
	}

	// 3. Compare release parts (empty release is treated as empty string)
	return compareRPMVersionString(v.release, other.release)
}

// compareRPMVersionString compares two RPM version strings using RPM's rules
// This implements RPM's version comparison algorithm which alternates between
// comparing non-numeric and numeric segments
func compareRPMVersionString(a, b string) int {
	i, j := 0, 0

	for i < len(a) || j < len(b) {
		// Skip separators (. + - ~ ^)
		for i < len(a) && isSeparator(rune(a[i])) {
			i++
		}
		for j < len(b) && isSeparator(rune(b[j])) {
			j++
		}

		// Extract non-digit segments
		iStart := i
		for i < len(a) && !unicode.IsDigit(rune(a[i])) && !isSeparator(rune(a[i])) {
			i++
		}
		aNonDigit := a[iStart:i]

		jStart := j
		for j < len(b) && !unicode.IsDigit(rune(b[j])) && !isSeparator(rune(b[j])) {
			j++
		}
		bNonDigit := b[jStart:j]

		// Compare non-digit segments lexicographically
		// Special case: tilde (~) sorts before anything (including empty string)
		nonDigitCmp := compareRPMNonDigits(aNonDigit, bNonDigit)
		if nonDigitCmp != 0 {
			return nonDigitCmp
		}

		// Extract digit segments
		iStart = i
		for i < len(a) && unicode.IsDigit(rune(a[i])) {
			i++
		}
		aDigit := a[iStart:i]

		jStart = j
		for j < len(b) && unicode.IsDigit(rune(b[j])) {
			j++
		}
		bDigit := b[jStart:j]

		// Compare digit segments numerically
		digitCmp := compareRPMDigits(aDigit, bDigit)
		if digitCmp != 0 {
			return digitCmp
		}
	}

	return 0
}

// isSeparator checks if a character is a separator in RPM versions
func isSeparator(r rune) bool {
	return r == '.' || r == '+' || r == '-' || r == '^'
}

// compareRPMNonDigits compares non-digit segments with RPM-specific rules
func compareRPMNonDigits(a, b string) int {
	// Handle special case where one string starts with tilde
	aHasTilde := strings.HasPrefix(a, "~")
	bHasTilde := strings.HasPrefix(b, "~")

	if aHasTilde && !bHasTilde {
		return -1 // Tilde sorts before non-tilde
	}
	if !aHasTilde && bHasTilde {
		return 1 // Non-tilde sorts after tilde
	}

	// Both have tilde or both don't have tilde - compare lexicographically
	return strings.Compare(a, b)
}

// compareRPMDigits compares digit strings numerically (leading zeros ignored)
func compareRPMDigits(a, b string) int {
	// Empty string is treated as 0
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return -1
	}
	if b == "" {
		return 1
	}

	// Convert to integers for comparison (this handles leading zeros correctly)
	aNum, aErr := strconv.ParseUint(a, 10, 64)
	bNum, bErr := strconv.ParseUint(b, 10, 64)

	if aErr == nil && bErr == nil {
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
		return 0
	}

	// Fallback for very large numbers that don't fit in uint64
	// Compare by length first (longer number is larger)
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}

	// If lengths are equal, string comparison works for digits
	return strings.Compare(a, b)
}
