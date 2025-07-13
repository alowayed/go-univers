// Package main demonstrates complex NPM version range operations and OR logic.
package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/alowayed/go-univers/pkg/ecosystem/npm"
)

func main() {
	fmt.Println("=== Complex NPM Version Range Operations ===")

	// Complex range parsing and evaluation
	fmt.Println("\n1. Complex Range Patterns:")
	complexRanges := []string{
		"^1.2.3",           // Caret range
		"~1.2.3",           // Tilde range  
		"1.x",              // X-range
		"1.2.3 - 2.0.0",    // Hyphen range
		">=1.0.0 <2.0.0",   // Multiple constraints
		"1.x || 2.x",       // OR logic
		"^0.1.2",           // Zero major caret
		"^0.0.1",           // Zero major/minor caret
		">=1.0.0-alpha <2.0.0", // Prerelease constraints
		"1.x || >=2.0.0-alpha <3.0.0", // Complex OR with prerelease
		"(>=1.0.0 <1.5.0) || (>=2.0.0 <3.0.0)", // Parentheses grouping
	}

	for _, rangeStr := range complexRanges {
		vr, err := npm.NewVersionRange(rangeStr)
		if err != nil {
			log.Printf("Failed to parse range %s: %v", rangeStr, err)
			continue
		}
		fmt.Printf("  ✅ Parsed: %s\n", rangeStr)
		
		// Test with sample versions
		testVersions := []string{"0.9.0", "1.2.5", "1.9.9", "2.0.0", "2.5.0-beta"}
		matching := []string{}
		
		for _, versionStr := range testVersions {
			version, err := npm.NewVersion(versionStr)
			if err != nil {
				continue
			}
			if vr.Contains(version) {
				matching = append(matching, versionStr)
			}
		}
		
		if len(matching) > 0 {
			fmt.Printf("     Matches: %v\n", matching)
		} else {
			fmt.Printf("     No matches from test set\n")
		}
	}

	// OR logic demonstration
	fmt.Println("\n2. OR Logic Examples:")
	orTests := []struct {
		range_   string
		versions []string
	}{
		{
			range_:   "1.x || 2.x",
			versions: []string{"0.9.0", "1.2.3", "1.9.9", "2.0.0", "2.5.0", "3.0.0"},
		},
		{
			range_:   "^1.0.0 || ^2.0.0",
			versions: []string{"0.9.0", "1.2.3", "1.9.9", "2.0.0", "2.5.0", "3.0.0"},
		},
		{
			range_:   "~1.2.0 || >=2.0.0-alpha <3.0.0",
			versions: []string{"1.1.9", "1.2.5", "1.3.0", "2.0.0-alpha", "2.5.0-beta", "2.9.9", "3.0.0"},
		},
	}

	for _, test := range orTests {
		fmt.Printf("  Range: %s\n", test.range_)
		vr, err := npm.NewVersionRange(test.range_)
		if err != nil {
			log.Printf("Failed to parse range: %v", err)
			continue
		}

		for _, versionStr := range test.versions {
			version, err := npm.NewVersion(versionStr)
			if err != nil {
				continue
			}
			
			match := vr.Contains(version)
			status := "❌"
			if match {
				status = "✅"
			}
			fmt.Printf("    %s %s\n", status, versionStr)
		}
		fmt.Println()
	}

	// Prerelease version handling
	fmt.Println("3. Prerelease Version Handling:")
	prereleaseTests := []struct {
		range_   string
		versions []string
	}{
		{
			range_:   "^1.0.0-alpha",
			versions: []string{"1.0.0-alpha", "1.0.0-beta", "1.0.0", "1.1.0-alpha", "1.1.0", "2.0.0"},
		},
		{
			range_:   ">=1.0.0-alpha <2.0.0",
			versions: []string{"0.9.0", "1.0.0-alpha", "1.0.0-beta", "1.0.0", "1.5.0", "2.0.0-alpha", "2.0.0"},
		},
	}

	for _, test := range prereleaseTests {
		fmt.Printf("  Range: %s\n", test.range_)
		vr, err := npm.NewVersionRange(test.range_)
		if err != nil {
			log.Printf("Failed to parse range: %v", err)
			continue
		}

		for _, versionStr := range test.versions {
			version, err := npm.NewVersion(versionStr)
			if err != nil {
				continue
			}
			
			match := vr.Contains(version)
			status := "❌"
			if match {
				status = "✅"
			}
			fmt.Printf("    %s %s\n", status, versionStr)
		}
		fmt.Println()
	}

	// Version sorting demonstration
	fmt.Println("4. Version Sorting:")
	unsortedVersions := []string{
		"2.0.0", "1.0.0-alpha", "1.0.0", "1.0.0-beta", "1.0.0-alpha.1", 
		"1.0.0-alpha.2", "10.0.0", "1.10.0", "1.2.0", "1.2.0-alpha",
	}

	versions := make([]*npm.Version, 0, len(unsortedVersions))
	for _, vStr := range unsortedVersions {
		v, err := npm.NewVersion(vStr)
		if err != nil {
			continue
		}
		versions = append(versions, v)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Compare(versions[j]) < 0
	})

	fmt.Printf("  Unsorted: %v\n", unsortedVersions)
	sortedStrs := make([]string, len(versions))
	for i, v := range versions {
		sortedStrs[i] = v.String()
	}
	fmt.Printf("  Sorted:   %v\n", sortedStrs)

	fmt.Println("\n=== Complex Example Complete ===")
}