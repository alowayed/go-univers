package alpm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Version represents an ALMP package version
type Version struct {
	epoch     int    // optional epoch (defaults to 0)
	pkgver    string // package version (upstream software version)
	pkgrel    int    // optional package release number (defaults to 0)
	hasPkgrel bool   // whether pkgrel was explicitly provided
	original  string // original version string
}

// NewVersion creates a new ALMP version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(version)

	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Split on epoch first
	var epochStr string
	var versionPart string

	if colonIndex := strings.Index(version, ":"); colonIndex != -1 {
		epochStr = version[:colonIndex]
		versionPart = version[colonIndex+1:]
	} else {
		epochStr = ""
		versionPart = version
	}

	// Split version part on last hyphen followed by digits (pkgrel)
	var pkgver string
	var pkgrelStr string

	// Find the last hyphen followed by only digits
	lastValidHyphen := -1
	for i := len(versionPart) - 1; i >= 0; i-- {
		if versionPart[i] == '-' && i+1 < len(versionPart) {
			// Check if everything after this hyphen is digits
			afterHyphen := versionPart[i+1:]
			if len(afterHyphen) > 0 && isAllDigits(afterHyphen) {
				lastValidHyphen = i
				break
			}
		}
	}

	if lastValidHyphen != -1 {
		pkgver = versionPart[:lastValidHyphen]
		pkgrelStr = versionPart[lastValidHyphen+1:]
	} else {
		pkgver = versionPart
		pkgrelStr = ""
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

	// Validate pkgver (cannot be empty)
	if pkgver == "" {
		return nil, fmt.Errorf("pkgver cannot be empty in version %s", original)
	}

	// Validate characters in pkgver - ALMP allows alphanumerics, periods, underscores
	if err := validateALMPVersionString(pkgver, "pkgver"); err != nil {
		return nil, fmt.Errorf("invalid pkgver in %s: %v", original, err)
	}

	// Parse pkgrel (track if it was explicitly provided)
	pkgrel := 0
	hasPkgrel := pkgrelStr != ""
	if hasPkgrel {
		var err error
		pkgrel, err = strconv.Atoi(pkgrelStr)
		if err != nil {
			return nil, fmt.Errorf("invalid pkgrel in version %s: %v", original, err)
		}
		if pkgrel < 0 {
			return nil, fmt.Errorf("pkgrel cannot be negative in version %s", original)
		}
	}

	return &Version{
		epoch:     epoch,
		pkgver:    pkgver,
		pkgrel:    pkgrel,
		hasPkgrel: hasPkgrel,
		original:  original,
	}, nil
}

// validateALMPVersionString validates that a version string contains only allowed characters
// ALMP allows: alphanumerics, periods, underscores, plus signs, hyphens
func validateALMPVersionString(s, part string) error {
	for _, r := range s {
		if !isValidALMPVersionChar(r) {
			return fmt.Errorf("invalid character %q in %s", r, part)
		}
	}
	return nil
}

// isValidALMPVersionChar checks if a character is valid in an ALMP version string
// Based on Arch Linux PKGBUILD specification
func isValidALMPVersionChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '+' || r == '-'
}

// isAllDigits checks if a string contains only digits
func isAllDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another ALMP version using vercmp rules
// Follows Arch Linux vercmp(8) algorithm:
// 1. Compare epochs first (higher epoch wins)
// 2. Compare pkgver parts using vercmp's alphanumeric rules
// 3. Compare pkgrel parts numerically
func (v *Version) Compare(other *Version) int {
	// 1. Compare epochs first
	if v.epoch != other.epoch {
		if v.epoch < other.epoch {
			return -1
		}
		return 1
	}

	// 2. Compare pkgver parts using vercmp rules
	pkgverCmp := compareALMPVersionString(v.pkgver, other.pkgver)
	if pkgverCmp != 0 {
		return pkgverCmp
	}

	// 3. Compare pkgrel parts with special handling for missing pkgrel
	// According to vercmp: comparing "1.5-1" and "1.5" yields 0
	if !v.hasPkgrel && !other.hasPkgrel {
		return 0 // Both have no pkgrel
	}
	if !v.hasPkgrel && other.hasPkgrel {
		return 0 // Special case: no pkgrel == with pkgrel when versions match
	}
	if v.hasPkgrel && !other.hasPkgrel {
		return 0 // Special case: with pkgrel == no pkgrel when versions match
	}

	// Both have pkgrel, compare numerically
	if v.pkgrel < other.pkgrel {
		return -1
	}
	if v.pkgrel > other.pkgrel {
		return 1
	}

	return 0
}

// compareALMPVersionString compares two ALMP version strings using vercmp rules
// This implements the Arch Linux vercmp algorithm based on the precedence:
// 1.0a < 1.0b < 1.0beta < 1.0p < 1.0pre < 1.0rc < 1.0 < 1.0.a < 1.0.1
func compareALMPVersionString(a, b string) int {
	// Handle the specific documented precedence cases first
	if a == b {
		return 0
	}

	// Check if this is a direct suffix comparison (no dots separating)
	if isDirectSuffixComparison(a, b) {
		return compareDirectSuffixes(a, b)
	}

	// Otherwise use standard segment-by-segment comparison
	return compareSegmentBySegment(a, b)
}

// isDirectSuffixComparison checks if we're comparing like "1.0" vs "1.0rc"
func isDirectSuffixComparison(a, b string) bool {
	// Simple heuristic: if one is a prefix of the other without separators
	if len(a) < len(b) && b[:len(a)] == a {
		// Check if remainder is alpha (no separators)
		remainder := b[len(a):]
		return len(remainder) > 0 && unicode.IsLetter(rune(remainder[0])) &&
			!strings.ContainsAny(remainder[:1], ".+-_")
	}
	if len(b) < len(a) && a[:len(b)] == b {
		// Check if remainder is alpha (no separators)
		remainder := a[len(b):]
		return len(remainder) > 0 && unicode.IsLetter(rune(remainder[0])) &&
			!strings.ContainsAny(remainder[:1], ".+-_")
	}
	return false
}

// compareDirectSuffixes handles cases like "1.0" vs "1.0rc"
func compareDirectSuffixes(a, b string) int {
	if len(a) < len(b) && b[:len(a)] == a {
		// a is prefix of b, b has direct suffix -> a wins (1.0 > 1.0rc)
		return 1
	}
	if len(b) < len(a) && a[:len(b)] == b {
		// b is prefix of a, a has direct suffix -> b wins
		return -1
	}
	// Both have suffixes, compare lexicographically
	return strings.Compare(a, b)
}

// compareSegmentBySegment does standard version segment comparison
// This implements a more accurate vercmp-style algorithm
func compareSegmentBySegment(a, b string) int {
	// Convert to segments first, handling delimiters properly
	aSegments := splitToSegments(a)
	bSegments := splitToSegments(b)

	// Compare segment by segment
	maxLen := len(aSegments)
	if len(bSegments) > maxLen {
		maxLen = len(bSegments)
	}

	for i := 0; i < maxLen; i++ {
		var aSeg, bSeg string
		if i < len(aSegments) {
			aSeg = aSegments[i]
		}
		if i < len(bSegments) {
			bSeg = bSegments[i]
		}

		// Compare segments
		cmp := compareSegments(aSeg, bSeg)
		if cmp != 0 {
			return cmp
		}
	}

	return 0
}

// splitToSegments splits a version string into segments, preserving empty segments
func splitToSegments(version string) []string {
	var segments []string
	var current strings.Builder

	for _, r := range version {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else {
			// Delimiter found - end current segment
			segments = append(segments, current.String())
			current.Reset()
		}
	}

	// Add final segment
	segments = append(segments, current.String())

	return segments
}

// compareSegments compares individual segments using vercmp rules
func compareSegments(a, b string) int {
	// Handle empty segments according to vercmp "final showdown" rules
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		// Empty segment vs non-empty segment
		// Based on vercmp results: more segments (even empty) = greater version
		return 1 // empty > non-empty (empty segments add "structure")
	}
	if b == "" {
		// Non-empty vs empty segment
		return -1 // non-empty < empty
	}

	// Both non-empty segments
	aIsNum := len(a) > 0 && unicode.IsDigit(rune(a[0]))
	bIsNum := len(b) > 0 && unicode.IsDigit(rune(b[0]))

	if aIsNum && bIsNum {
		return compareALMPDigits(a, b)
	} else if aIsNum {
		return 1 // numeric > alpha
	} else if bIsNum {
		return -1 // alpha < numeric
	} else {
		return strings.Compare(a, b) // both alpha
	}
}

// compareALMPDigits compares digit strings numerically
func compareALMPDigits(a, b string) int {
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

	// Convert to integers for comparison
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
