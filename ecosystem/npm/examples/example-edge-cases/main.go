// Package main demonstrates edge cases and error handling in NPM version parsing.
package main

import (
	"fmt"
	"log"

	"github.com/alowayed/go-univers/ecosystem/npm"
)

func main() {
	fmt.Println("=== NPM Version Edge Cases and Error Handling ===")

	// Edge cases for version parsing
	fmt.Println("\n1. Version Parsing Edge Cases:")
	versionEdgeCases := []struct {
		input       string
		expectError bool
		description string
	}{
		{"1.2.3", false, "Basic version"},
		{"v1.2.3", false, "Single v prefix"},
		{"vv1.2.3", false, "Double v prefix (allowed)"},
		{"vvv1.2.3", true, "Triple v prefix (invalid)"},
		{"=1.2.3", false, "Equals prefix"},
		{"01.02.03", false, "Zero-padded version"},
		{" 1.2.3 ", false, "Whitespace handling"},
		{"1.2.3\t", false, "Tab handling"},
		{"1.2.3-alpha.1", false, "Prerelease"},
		{"1.2.3+build.1", false, "Build metadata"},
		{"1.2.3-alpha.1+build.1", false, "Both prerelease and build"},
		{"1.2.3-01", false, "Zero-padded prerelease"},
		{"1.2.3-alpha.1.beta", false, "Complex prerelease"},
		{"1.2.3+build-01.02", false, "Build with hyphens"},
		{"", true, "Empty string"},
		{"1.2", true, "Missing patch"},
		{"1.2.3.4", true, "Too many parts"},
		{"a.b.c", true, "Non-numeric"},
		{"   ", true, "Only whitespace"},
	}

	for _, test := range versionEdgeCases {
		version, err := npm.NewVersion(test.input)
		hasError := err != nil
		
		status := "✅"
		if hasError != test.expectError {
			status = "❌"
		}
		
		if !hasError && version != nil {
			fmt.Printf("  %s %-25s | %-30s | Result: %s\n", 
				status, fmt.Sprintf("'%s'", test.input), test.description, version.String())
		} else if hasError && test.expectError {
			fmt.Printf("  %s %-25s | %-30s | Expected error: %v\n", 
				status, fmt.Sprintf("'%s'", test.input), test.description, err)
		} else {
			fmt.Printf("  %s %-25s | %-30s | Unexpected result\n", 
				status, fmt.Sprintf("'%s'", test.input), test.description)
		}
	}

	// Edge cases for version range parsing
	fmt.Println("\n2. Version Range Parsing Edge Cases:")
	rangeEdgeCases := []struct {
		input       string
		expectError bool
		description string
	}{
		{"1.2.3", false, "Exact version"},
		{"^1.2.3", false, "Caret range"},
		{"~1.2.3", false, "Tilde range"},
		{"1.x", false, "X-range"},
		{"1.2.3 - 2.0.0", false, "Hyphen range"},
		{">=1.0.0 <2.0.0", false, "Multiple constraints"},
		{"1.x || 2.x", false, "OR logic"},
		{"*", false, "Wildcard"},
		{" ^1.2.3 ", false, "Range with whitespace"},
		{"1.x  ||  2.x", false, "OR with extra spaces"},
		{"(>=1.0.0 <2.0.0)", false, "Parentheses"},
		{"", true, "Empty range"},
		{"   ", true, "Only whitespace"},
		{"~~1.2.3", true, "Double tilde"},
		{"^^1.2.3", true, "Double caret"},
		{"1.2.3 -", true, "Malformed hyphen (missing end)"},
		{"- 1.2.3", true, "Malformed hyphen (missing start)"},
		{"1.2.3 - - 2.0.0", true, "Double hyphen"},
		{"1.2.3@invalid", true, "Invalid characters"},
	}

	for _, test := range rangeEdgeCases {
		vr, err := npm.NewVersionRange(test.input)
		hasError := err != nil
		
		status := "✅"
		if hasError != test.expectError {
			status = "❌"
		}
		
		if !hasError && vr != nil {
			fmt.Printf("  %s %-25s | %-30s | Result: %s\n", 
				status, fmt.Sprintf("'%s'", test.input), test.description, vr.String())
		} else if hasError && test.expectError {
			fmt.Printf("  %s %-25s | %-30s | Expected error\n", 
				status, fmt.Sprintf("'%s'", test.input), test.description)
		} else {
			fmt.Printf("  %s %-25s | %-30s | Unexpected result\n", 
				status, fmt.Sprintf("'%s'", test.input), test.description)
		}
	}

	// Zero version caret range edge cases
	fmt.Println("\n3. Zero Version Caret Range Special Cases:")
	zeroCaretTests := []struct {
		range_     string
		version    string
		shouldMatch bool
		explanation string
	}{
		{"^0.0.1", "0.0.1", true, "Exact match for ^0.0.x"},
		{"^0.0.1", "0.0.2", false, "^0.0.x only matches exact patch"},
		{"^0.1.0", "0.1.0", true, "Exact match for ^0.x.y"},
		{"^0.1.0", "0.1.5", true, "^0.x.y allows patch increments"},
		{"^0.1.0", "0.2.0", false, "^0.x.y doesn't allow minor increments"},
		{"^1.0.0", "1.5.0", true, "Normal caret behavior for non-zero major"},
		{"^1.0.0", "2.0.0", false, "Caret doesn't cross major versions"},
	}

	for _, test := range zeroCaretTests {
		vr, err := npm.NewVersionRange(test.range_)
		if err != nil {
			log.Printf("Failed to parse range %s: %v", test.range_, err)
			continue
		}
		
		version, err := npm.NewVersion(test.version)
		if err != nil {
			log.Printf("Failed to parse version %s: %v", test.version, err)
			continue
		}
		
		matches := vr.Contains(version)
		status := "✅"
		if matches != test.shouldMatch {
			status = "❌"
		}
		
		fmt.Printf("  %s %s contains %s: %v | %s\n", 
			status, test.range_, test.version, matches, test.explanation)
	}

	// Build metadata handling (should be ignored in comparisons)
	fmt.Println("\n4. Build Metadata Handling:")
	buildTests := []struct {
		v1, v2      string
		shouldEqual bool
	}{
		{"1.2.3+build1", "1.2.3+build2", true},
		{"1.2.3+build.123", "1.2.3", true},
		{"1.2.3-alpha+build1", "1.2.3-alpha+build2", true},
		{"1.2.3-alpha+build", "1.2.3-beta+build", false},
	}

	for _, test := range buildTests {
		v1, err1 := npm.NewVersion(test.v1)
		v2, err2 := npm.NewVersion(test.v2)
		if err1 != nil || err2 != nil {
			continue
		}
		
		equal := v1.Compare(v2) == 0
		status := "✅"
		if equal != test.shouldEqual {
			status = "❌"
		}
		
		fmt.Printf("  %s %s == %s: %v (build metadata ignored)\n", 
			status, test.v1, test.v2, equal)
	}

	// Prerelease boundary cases
	fmt.Println("\n5. Prerelease Boundary Cases:")
	prereleaseTests := []struct {
		range_     string
		version    string
		shouldMatch bool
		explanation string
	}{
		{">=1.0.0", "1.0.0-alpha", false, "Prerelease < release for same version"},
		{">=1.0.0-alpha", "1.0.0-beta", true, "Beta > alpha"},
		{">=1.0.0-alpha", "1.0.0", true, "Release > prerelease"},
		{"^1.0.0-alpha", "1.0.0", true, "Caret includes release from prerelease base"},
		{"^1.0.0-alpha", "1.1.0-alpha", true, "Caret includes minor prerelease"},
		{"1.x", "1.0.0-alpha", true, "X-range includes prereleases"},
	}

	for _, test := range prereleaseTests {
		vr, err := npm.NewVersionRange(test.range_)
		if err != nil {
			log.Printf("Failed to parse range %s: %v", test.range_, err)
			continue
		}
		
		version, err := npm.NewVersion(test.version)
		if err != nil {
			log.Printf("Failed to parse version %s: %v", test.version, err)
			continue
		}
		
		matches := vr.Contains(version)
		status := "✅"
		if matches != test.shouldMatch {
			status = "❌"
		}
		
		fmt.Printf("  %s %s contains %s: %v | %s\n", 
			status, test.range_, test.version, matches, test.explanation)
	}

	fmt.Println("\n=== Edge Cases Example Complete ===")
}