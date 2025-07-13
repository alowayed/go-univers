package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/gomod"
	"github.com/alowayed/go-univers/pkg/ecosystem/npm"
	"github.com/alowayed/go-univers/pkg/ecosystem/pypi"
)

// compare compares two versions and outputs -1, 0, or 1
func compare(ecosystem string, args []string) int {
	if len(args) != 2 {
		fmt.Println("compare requires exactly 2 version arguments")
		return 1
	}

	switch ecosystem {
	case "npm":
		return compareNPM(args[0], args[1])
	case "pypi":
		return comparePyPI(args[0], args[1])
	case "go":
		return compareGo(args[0], args[1])
	default:
		fmt.Printf("Unknown ecosystem: %s\n", ecosystem)
		return 1
	}
}

// sort sorts multiple versions and outputs them in ascending order
func sort(ecosystem string, args []string) int {
	if len(args) == 0 {
		fmt.Println("sort requires at least 1 version argument")
		return 1
	}

	switch ecosystem {
	case "npm":
		return sortNPM(args)
	case "pypi":
		return sortPyPI(args)
	case "go":
		return sortGo(args)
	default:
		fmt.Printf("Unknown ecosystem: %s\n", ecosystem)
		return 1
	}
}

// satisfies checks if a version satisfies a range (exit code 0/1)
func satisfies(ecosystem string, args []string) int {
	if len(args) != 2 {
		fmt.Println("satisfies requires exactly 2 arguments: <version> <range>")
		return 1
	}

	version := args[0]
	rangeStr := args[1]

	switch ecosystem {
	case "npm":
		return satisfiesNPM(version, rangeStr)
	case "pypi":
		return satisfiesPyPI(version, rangeStr)
	case "go":
		return satisfiesGo(version, rangeStr)
	default:
		fmt.Printf("Unknown ecosystem: %s\n", ecosystem)
		return 1
	}
}

// NPM command implementations

func compareNPM(v1, v2 string) int {
	ver1, err := npm.NewVersion(v1)
	if err != nil {
		fmt.Printf("Invalid NPM version '%s': %v\n", v1, err)
		return 1
	}

	ver2, err := npm.NewVersion(v2)
	if err != nil {
		fmt.Printf("Invalid NPM version '%s': %v\n", v2, err)
		return 1
	}

	fmt.Println(ver1.Compare(ver2))
	return 0
}

func sortNPM(versionStrings []string) int {
	versions := make([]*npm.Version, 0, len(versionStrings))

	for _, vStr := range versionStrings {
		ver, err := npm.NewVersion(vStr)
		if err != nil {
			fmt.Printf("Invalid NPM version '%s': %v\n", vStr, err)
			return 1
		}
		versions = append(versions, ver)
	}

	slices.SortFunc(versions, (*npm.Version).Compare)

	// Output sorted versions
	sortedStrings := make([]string, len(versions))
	for i, ver := range versions {
		sortedStrings[i] = ver.String()
	}
	fmt.Println(strings.Join(sortedStrings, ", "))
	return 0
}

func satisfiesNPM(version, rangeStr string) int {
	ver, err := npm.NewVersion(version)
	if err != nil {
		fmt.Printf("Invalid NPM version '%s': %v\n", version, err)
		return 1
	}

	vrange, err := npm.NewVersionRange(rangeStr)
	if err != nil {
		fmt.Printf("Invalid NPM range '%s': %v\n", rangeStr, err)
		return 1
	}

	if vrange.Contains(ver) {
		return 0 // Success: version satisfies range
	}
	return 1 // Failure: version does not satisfy range
}

// PyPI command implementations

func comparePyPI(v1, v2 string) int {
	ver1, err := pypi.NewVersion(v1)
	if err != nil {
		fmt.Printf("Invalid PyPI version '%s': %v\n", v1, err)
		return 1
	}

	ver2, err := pypi.NewVersion(v2)
	if err != nil {
		fmt.Printf("Invalid PyPI version '%s': %v\n", v2, err)
		return 1
	}

	fmt.Println(ver1.Compare(ver2))
	return 0
}

func sortPyPI(versionStrings []string) int {
	versions := make([]*pypi.Version, 0, len(versionStrings))

	for _, vStr := range versionStrings {
		ver, err := pypi.NewVersion(vStr)
		if err != nil {
			fmt.Printf("Invalid PyPI version '%s': %v\n", vStr, err)
			return 1
		}
		versions = append(versions, ver)
	}

	slices.SortFunc(versions, (*pypi.Version).Compare)

	// Output sorted versions
	sortedStrings := make([]string, len(versions))
	for i, ver := range versions {
		sortedStrings[i] = ver.String()
	}
	fmt.Println(strings.Join(sortedStrings, ", "))
	return 0
}

func satisfiesPyPI(version, rangeStr string) int {
	ver, err := pypi.NewVersion(version)
	if err != nil {
		fmt.Printf("Invalid PyPI version '%s': %v\n", version, err)
		return 1
	}

	vrange, err := pypi.NewVersionRange(rangeStr)
	if err != nil {
		fmt.Printf("Invalid PyPI range '%s': %v\n", rangeStr, err)
		return 1
	}

	if vrange.Contains(ver) {
		return 0 // Success: version satisfies range
	}
	return 1 // Failure: version does not satisfy range
}

// Go command implementations

func compareGo(v1, v2 string) int {
	ver1, err := gomod.NewVersion(v1)
	if err != nil {
		fmt.Printf("Invalid Go version '%s': %v\n", v1, err)
		return 1
	}

	ver2, err := gomod.NewVersion(v2)
	if err != nil {
		fmt.Printf("Invalid Go version '%s': %v\n", v2, err)
		return 1
	}

	fmt.Println(ver1.Compare(ver2))
	return 0
}

func sortGo(versionStrings []string) int {
	versions := make([]*gomod.Version, 0, len(versionStrings))

	for _, vStr := range versionStrings {
		ver, err := gomod.NewVersion(vStr)
		if err != nil {
			fmt.Printf("Invalid Go version '%s': %v\n", vStr, err)
			return 1
		}
		versions = append(versions, ver)
	}

	slices.SortFunc(versions, (*gomod.Version).Compare)

	// Output sorted versions
	sortedStrings := make([]string, len(versions))
	for i, ver := range versions {
		sortedStrings[i] = ver.String()
	}
	fmt.Println(strings.Join(sortedStrings, ", "))
	return 0
}

func satisfiesGo(version, rangeStr string) int {
	ver, err := gomod.NewVersion(version)
	if err != nil {
		fmt.Printf("Invalid Go version '%s': %v\n", version, err)
		return 1
	}

	vrange, err := gomod.NewVersionRange(rangeStr)
	if err != nil {
		fmt.Printf("Invalid Go range '%s': %v\n", rangeStr, err)
		return 1
	}

	if vrange.Contains(ver) {
		return 0 // Success: version satisfies range
	}
	return 1 // Failure: version does not satisfy range
}