package main

import (
	"fmt"
	"slices"

	"github.com/alowayed/go-univers/ecosystem/npm"
)

func main() {
	fmt.Println("NPM Version Sorting Example")
	fmt.Println("==========================")

	// Create some versions in random order
	versionStrings := []string{
		"2.1.0",
		"1.0.0",
		"2.0.0-alpha",
		"1.2.3",
		"2.0.0",
		"1.2.3-beta.1",
		"1.2.3-alpha",
		"1.2.3-beta.2",
		"3.0.0-rc.1",
		"1.0.0-alpha.1",
	}

	// Parse versions
	versions := make([]*npm.Version, 0, len(versionStrings))
	for _, vStr := range versionStrings {
		v, err := npm.NewVersion(vStr)
		if err != nil {
			fmt.Printf("Error parsing version %s: %v\n", vStr, err)
			continue
		}
		versions = append(versions, v)
	}

	fmt.Println("\nOriginal order:")
	for i, v := range versions {
		fmt.Printf("%d. %s\n", i+1, v)
	}

	// Sort versions using slices.SortFunc with the existing Compare method
	slices.SortFunc(versions, (*npm.Version).Compare)

	fmt.Println("\nSorted (ascending):")
	for i, v := range versions {
		fmt.Printf("%d. %s\n", i+1, v)
	}

	// Reverse sort for descending order
	slices.SortFunc(versions, func(a, b *npm.Version) int {
		return b.Compare(a) // Note: reversed order
	})

	fmt.Println("\nSorted (descending):")
	for i, v := range versions {
		fmt.Printf("%d. %s\n", i+1, v)
	}

	// Demonstrate stable sort with equal versions
	fmt.Println("\nStable Sort Example:")
	equalVersions := []*npm.Version{}
	
	// Create versions with different build metadata (equal for sorting)
	buildVersions := []string{"1.0.0+build.1", "1.0.0+build.2", "1.0.0"}
	for _, vStr := range buildVersions {
		v, _ := npm.NewVersion(vStr)
		equalVersions = append(equalVersions, v)
	}

	fmt.Println("Before stable sort:")
	for i, v := range equalVersions {
		fmt.Printf("%d. %s\n", i+1, v)
	}

	slices.SortStableFunc(equalVersions, (*npm.Version).Compare)

	fmt.Println("After stable sort (build metadata ignored in comparison):")
	for i, v := range equalVersions {
		fmt.Printf("%d. %s\n", i+1, v)
	}

	// Show semantic version ordering rules
	fmt.Println("\nSemantic Version Ordering Rules Demonstration:")
	ruleVersions := []string{
		"1.0.0",           // Normal version
		"1.0.0-alpha",     // Pre-release (lower than normal)
		"1.0.0-alpha.1",   // Pre-release with number
		"1.0.0-alpha.beta", // Pre-release with string
		"1.0.0-beta",      // Pre-release (higher than alpha)
		"1.0.0-beta.2",    // Pre-release with number
		"1.0.0-beta.11",   // Pre-release (numeric comparison)
		"1.0.0-rc.1",      // Release candidate
	}

	ruleParsed := make([]*npm.Version, 0, len(ruleVersions))
	for _, vStr := range ruleVersions {
		v, _ := npm.NewVersion(vStr)
		ruleParsed = append(ruleParsed, v)
	}

	slices.SortFunc(ruleParsed, (*npm.Version).Compare)

	fmt.Println("Properly sorted according to semver rules:")
	for i, v := range ruleParsed {
		fmt.Printf("%d. %s\n", i+1, v)
	}
}