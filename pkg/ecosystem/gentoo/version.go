package gentoo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionPattern matches Gentoo version strings
// Format: numbers, optional letter, optional suffixes (_alpha, _beta, _pre, _rc, _p), optional revision
var versionPattern = regexp.MustCompile(`^(\d+(?:\.\d+){0,10})([a-zA-Z])?(?:_(alpha|beta|pre|rc|p)(\d*))?(?:-r(\d+))?$`)

// suffixValues defines the ordering of version suffixes
var suffixValues = map[string]int{
	"alpha": -4,
	"beta":  -3,
	"pre":   -2,
	"rc":    -1,
	"p":     1,
}

// Version represents a Gentoo package version
type Version struct {
	numbers   []int
	letter    string
	suffix    string
	suffixNum int
	revision  int
	original  string
}

// NewVersion creates a new Gentoo version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(version)

	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid Gentoo version: %s", original)
	}

	// Parse numeric components
	numbersPart := matches[1]
	var numbers []int
	for _, numStr := range strings.Split(numbersPart, ".") {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("invalid numeric component: %s", numStr)
		}
		numbers = append(numbers, num)
	}

	// Parse letter component
	letter := matches[2]

	// Parse suffix
	suffix := matches[3]
	suffixNumStr := matches[4]
	suffixNum := 0
	if suffixNumStr != "" {
		var err error
		suffixNum, err = strconv.Atoi(suffixNumStr)
		if err != nil {
			return nil, fmt.Errorf("invalid suffix number: %s", suffixNumStr)
		}
	}

	// Parse revision
	revisionStr := matches[5]
	revision := 0
	if revisionStr != "" {
		var err error
		revision, err = strconv.Atoi(revisionStr)
		if err != nil {
			return nil, fmt.Errorf("invalid revision: %s", revisionStr)
		}
	}

	return &Version{
		numbers:   numbers,
		letter:    letter,
		suffix:    suffix,
		suffixNum: suffixNum,
		revision:  revision,
		original:  original,
	}, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another Gentoo version
func (v *Version) Compare(other *Version) int {
	// Compare numeric components
	maxLen := len(v.numbers)
	if len(other.numbers) > maxLen {
		maxLen = len(other.numbers)
	}

	for i := 0; i < maxLen; i++ {
		vNum := 0
		otherNum := 0

		if i < len(v.numbers) {
			vNum = v.numbers[i]
		}
		if i < len(other.numbers) {
			otherNum = other.numbers[i]
		}

		if vNum != otherNum {
			return compareInt(vNum, otherNum)
		}
	}

	// Compare letter components
	if letterCmp := strings.Compare(v.letter, other.letter); letterCmp != 0 {
		return letterCmp
	}

	// Compare suffixes
	vSuffixValue := 0
	otherSuffixValue := 0

	if v.suffix != "" {
		vSuffixValue = suffixValues[v.suffix]
	}
	if other.suffix != "" {
		otherSuffixValue = suffixValues[other.suffix]
	}

	if vSuffixValue != otherSuffixValue {
		return compareInt(vSuffixValue, otherSuffixValue)
	}

	// If same suffix type, compare suffix numbers
	if v.suffix != "" && other.suffix != "" && v.suffix == other.suffix {
		if v.suffixNum != other.suffixNum {
			return compareInt(v.suffixNum, other.suffixNum)
		}
	}

	// Compare revisions
	return compareInt(v.revision, other.revision)
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
