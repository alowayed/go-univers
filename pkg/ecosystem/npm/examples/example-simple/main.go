// Package main demonstrates basic NPM version parsing and comparison functionality.
package main

import (
	"fmt"
	"log"

	"github.com/alowayed/go-univers/pkg/ecosystem/npm"
)

func main() {
	fmt.Println("=== Basic NPM Version Operations ===")

	// Basic version parsing
	versions := []string{
		"1.2.3",
		"v2.0.0",
		"=3.1.4",
		"1.0.0-alpha.1",
		"2.5.0+build.123",
		"1.2.3-beta.2+build.456",
	}

	fmt.Println("\n1. Version Parsing:")
	parsedVersions := make([]*npm.Version, 0, len(versions))
	for _, versionStr := range versions {
		version, err := npm.NewVersion(versionStr)
		if err != nil {
			log.Printf("Failed to parse version %s: %v", versionStr, err)
			continue
		}
		parsedVersions = append(parsedVersions, version)
		fmt.Printf("  Input: %-20s | Parsed: %-15s | Normalized: %s\n", 
			versionStr, version.String(), version.Normalize())
	}

	// Version comparison
	fmt.Println("\n2. Version Comparison:")
	comparisons := []struct {
		v1, v2 string
	}{
		{"1.2.3", "1.2.4"},
		{"2.0.0", "1.9.9"},
		{"1.0.0-alpha", "1.0.0"},
		{"1.0.0-alpha.1", "1.0.0-alpha.2"},
		{"1.0.0+build1", "1.0.0+build2"},
		{"1.0.0-beta", "1.0.0-alpha"},
	}

	for _, comp := range comparisons {
		v1, err1 := npm.NewVersion(comp.v1)
		v2, err2 := npm.NewVersion(comp.v2)
		if err1 != nil || err2 != nil {
			continue
		}

		result := v1.Compare(v2)
		var relation string
		switch {
		case result < 0:
			relation = "<"
		case result > 0:
			relation = ">"
		default:
			relation = "="
		}
		fmt.Printf("  %s %s %s\n", comp.v1, relation, comp.v2)
	}

	// Simple range matching
	fmt.Println("\n3. Simple Range Matching:")
	rangeTests := []struct {
		version string
		ranges  []string
	}{
		{
			version: "1.2.5",
			ranges:  []string{"^1.2.0", "~1.2.3", ">=1.0.0", "1.x"},
		},
		{
			version: "2.0.0-beta.1",
			ranges:  []string{"^2.0.0", ">=2.0.0-alpha", "2.x"},
		},
	}

	for _, test := range rangeTests {
		version, err := npm.NewVersion(test.version)
		if err != nil {
			continue
		}
		
		fmt.Printf("  Version %s satisfies:\n", test.version)
		for _, rangeStr := range test.ranges {
			vr, err := npm.NewVersionRange(rangeStr)
			if err != nil {
				continue
			}
			
			satisfies := vr.Contains(version)
			status := "❌"
			if satisfies {
				status = "✅"
			}
			fmt.Printf("    %s %s\n", status, rangeStr)
		}
		fmt.Println()
	}

	fmt.Println("=== Example Complete ===")
}