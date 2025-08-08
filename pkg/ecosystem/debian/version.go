package debian

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// versionPattern matches Debian version strings
// Format: [epoch:]upstream_version[-debian_revision]
var versionPattern = regexp.MustCompile(`^(?:(\d+):)?(.+?)(?:-([^-]+))?$`)

// Version represents a Debian package version
type Version struct {
	epoch    int    // optional epoch (defaults to 0)
	upstream string // upstream version part
	revision string // optional Debian revision (empty for native packages)
	original string // original version string
}

// NewVersion creates a new Debian version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(version)

	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Parse using regex
	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid Debian version format: %s", original)
	}

	epochStr := matches[1]
	upstream := matches[2]
	revision := matches[3]

	// Parse epoch (default to 0)
	epoch := 0
	if epochStr != "" {
		var err error
		epoch, err = strconv.Atoi(epochStr)
		if err != nil {
			return nil, fmt.Errorf("invalid epoch in version %s: %v", original, err)
		}
	}

	// Validate upstream version
	if upstream == "" {
		return nil, fmt.Errorf("upstream version cannot be empty in version %s", original)
	}

	// Upstream version must start with a digit per Debian policy
	if !unicode.IsDigit(rune(upstream[0])) {
		return nil, fmt.Errorf("upstream version must start with a digit in version %s", original)
	}

	// Validate characters in upstream version
	if err := validateVersionString(upstream, "upstream version"); err != nil {
		return nil, fmt.Errorf("invalid upstream version in %s: %v", original, err)
	}

	// Validate characters in revision (if present)
	if revision != "" {
		if err := validateVersionString(revision, "debian revision"); err != nil {
			return nil, fmt.Errorf("invalid debian revision in %s: %v", original, err)
		}
	}

	return &Version{
		epoch:    epoch,
		upstream: upstream,
		revision: revision,
		original: original,
	}, nil
}

// validateVersionString validates that a version string contains only allowed characters
func validateVersionString(s, part string) error {
	for _, r := range s {
		if !isValidVersionChar(r) {
			return fmt.Errorf("invalid character %q in %s", r, part)
		}
	}
	return nil
}

// isValidVersionChar checks if a character is valid in a Debian version string
// Per Debian policy: alphanumerics and the characters . + - ~ (and : for epochs)
func isValidVersionChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '+' || r == '-' || r == '~'
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another Debian version
// Follows dpkg --compare-versions behavior
func (v *Version) Compare(other *Version) int {
	// 1. Compare epochs first
	if v.epoch != other.epoch {
		if v.epoch < other.epoch {
			return -1
		}
		return 1
	}

	// 2. Compare upstream versions
	upstreamCmp := compareDebianVersionString(v.upstream, other.upstream)
	if upstreamCmp != 0 {
		return upstreamCmp
	}

	// 3. Compare revisions (empty revision is treated as "0")
	vRevision := v.revision
	otherRevision := other.revision

	// Native packages (no revision) have implicit revision "0"
	if vRevision == "" {
		vRevision = "0"
	}
	if otherRevision == "" {
		otherRevision = "0"
	}

	return compareDebianVersionString(vRevision, otherRevision)
}

// compareDebianVersionString compares two Debian version strings using Debian's rules
// This implements the dpkg version comparison algorithm
func compareDebianVersionString(a, b string) int {
	i, j := 0, 0

	for i < len(a) || j < len(b) {
		// Extract non-digit prefix
		iStart := i
		for i < len(a) && !unicode.IsDigit(rune(a[i])) {
			i++
		}
		aNonDigit := a[iStart:i]

		jStart := j
		for j < len(b) && !unicode.IsDigit(rune(b[j])) {
			j++
		}
		bNonDigit := b[jStart:j]

		// Compare non-digit parts lexicographically with special tilde handling
		nonDigitCmp := compareDebianNonDigits(aNonDigit, bNonDigit)
		if nonDigitCmp != 0 {
			return nonDigitCmp
		}

		// Extract digit prefix
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

		// Compare digit parts numerically
		digitCmp := compareDebianDigits(aDigit, bDigit)
		if digitCmp != 0 {
			return digitCmp
		}
	}

	return 0
}

// compareDebianNonDigits compares non-digit parts with Debian-specific rules
func compareDebianNonDigits(a, b string) int {
	maxLen := max(len(a), len(b))

	for i := range maxLen {
		var aChar, bChar rune

		// Get character or treat missing as null (sorts before anything)
		if i < len(a) {
			aChar = rune(a[i])
		} else {
			aChar = 0 // null character
		}
		if i < len(b) {
			bChar = rune(b[i])
		} else {
			bChar = 0 // null character
		}

		// Apply Debian character weights
		aWeight := getDebianCharWeight(aChar)
		bWeight := getDebianCharWeight(bChar)

		if aWeight != bWeight {
			if aWeight < bWeight {
				return -1
			}
			return 1
		}
	}

	return 0
}

// getDebianCharWeight returns the sort weight for a character per Debian rules
// Tilde (~) sorts earliest, then null, then letters/other chars
func getDebianCharWeight(r rune) int {
	switch r {
	case '~':
		return -1 // Tilde sorts before everything else
	case 0:
		return 0 // Null/missing character
	default:
		return int(r) // Use Unicode value for other characters
	}
}

// compareDebianDigits compares digit strings numerically
func compareDebianDigits(a, b string) int {
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

	// Fallback for very large numbers that don't fit in uint64.
	// Compare by length first.
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}

	// If lengths are equal, a string comparison is correct.
	return strings.Compare(a, b)
}
