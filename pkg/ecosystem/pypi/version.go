package pypi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	// versionPattern matches PyPI version strings according to PEP 440
	versionPattern = regexp.MustCompile(`^(?:([0-9]+)!)?([0-9]+(?:\.[0-9]+)*?)(?:\.?(a|b|rc|alpha|beta|c)([0-9]+))?(?:\.?(post|rev|r)([0-9]+))?(?:\.?(dev)([0-9]+))?(?:\+([a-zA-Z0-9]+(?:[-_.][a-zA-Z0-9]+)*))?$`)
)

// Version represents a PyPI package version following PEP 440
type Version struct {
	epoch       int
	release     []int
	prerelease  string
	preNumber   int
	postrelease int
	dev         int
	local       string
	original    string
}

// newVersion creates a new PyPI version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	version = strings.TrimSpace(version)
	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Parse using regex
	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid PyPI version format: %s", version)
	}

	pv := &Version{
		epoch:       0,
		postrelease: -1,
		dev:         -1,
		original:    version,
	}

	// Parse epoch (group 1)
	if matches[1] != "" {
		epoch, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid epoch: %s", matches[1])
		}
		pv.epoch = epoch
	}

	// Parse release version (group 2)
	if matches[2] == "" {
		return nil, fmt.Errorf("missing release version")
	}
	releaseParts := strings.Split(matches[2], ".")
	pv.release = make([]int, len(releaseParts))
	for i, part := range releaseParts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid release part: %s", part)
		}
		pv.release[i] = num
	}

	// Parse prerelease (groups 3, 4)
	if matches[3] != "" {
		pv.prerelease = matches[3]
		if matches[4] != "" {
			preNum, err := strconv.Atoi(matches[4])
			if err != nil {
				return nil, fmt.Errorf("invalid prerelease number: %s", matches[4])
			}
			pv.preNumber = preNum
		}
	}

	// Parse post-release (groups 5, 6)
	if matches[5] != "" {
		if matches[6] != "" {
			postNum, err := strconv.Atoi(matches[6])
			if err != nil {
				return nil, fmt.Errorf("invalid post number: %s", matches[6])
			}
			pv.postrelease = postNum
		} else {
			pv.postrelease = 0
		}
	}

	// Parse dev release (groups 7, 8)
	if matches[7] != "" {
		if matches[8] != "" {
			devNum, err := strconv.Atoi(matches[8])
			if err != nil {
				return nil, fmt.Errorf("invalid dev number: %s", matches[8])
			}
			pv.dev = devNum
		} else {
			pv.dev = 0
		}
	}

	// Parse local version (group 9)
	if matches[9] != "" {
		pv.local = matches[9]
	}

	return pv, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another PyPI version according to PEP 440
func (v *Version) Compare(other *Version) int {
	if v.epoch != other.epoch {
		return compareInt(v.epoch, other.epoch)
	}

	releaseComparison := compareReleaseVersions(v.release, other.release)
	if releaseComparison != 0 {
		return releaseComparison
	}

	preComparison := comparePrereleases(v.prerelease, v.preNumber, other.prerelease, other.preNumber)
	if preComparison != 0 {
		return preComparison
	}

	postComparison := comparePostReleases(v.postrelease, other.postrelease)
	if postComparison != 0 {
		return postComparison
	}

	return compareDevReleases(v.dev, other.dev)
}
