package vers

import "testing"

func TestAlpineContains(t *testing.T) {
	tests := []struct {
		name        string
		constraints []string
		version     string
		expected    bool
		expectError bool
	}{
		// Basic comparison operators
		{
			name:        "exact_match",
			constraints: []string{"=1.2.3"},
			version:     "1.2.3",
			expected:    true,
		},
		{
			name:        "exact_match_false",
			constraints: []string{"=1.2.3"},
			version:     "1.2.4",
			expected:    false,
		},
		{
			name:        "greater_than_equal",
			constraints: []string{">=1.2.0"},
			version:     "1.2.3",
			expected:    true,
		},
		{
			name:        "greater_than_equal_false",
			constraints: []string{">=1.2.3"},
			version:     "1.2.0",
			expected:    false,
		},
		{
			name:        "less_than_equal",
			constraints: []string{"<=1.2.3"},
			version:     "1.2.0",
			expected:    true,
		},
		{
			name:        "less_than_equal_false",
			constraints: []string{"<=1.2.0"},
			version:     "1.2.3",
			expected:    false,
		},
		{
			name:        "greater_than",
			constraints: []string{">1.2.0"},
			version:     "1.2.3",
			expected:    true,
		},
		{
			name:        "greater_than_false",
			constraints: []string{">1.2.3"},
			version:     "1.2.3",
			expected:    false,
		},
		{
			name:        "less_than",
			constraints: []string{"<1.2.3"},
			version:     "1.2.0",
			expected:    true,
		},
		{
			name:        "less_than_false",
			constraints: []string{"<1.2.0"},
			version:     "1.2.3",
			expected:    false,
		},
		{
			name:        "not_equal",
			constraints: []string{"!=1.2.0"},
			version:     "1.2.3",
			expected:    true,
		},
		{
			name:        "not_equal_false",
			constraints: []string{"!=1.2.3"},
			version:     "1.2.3",
			expected:    false,
		},

		// Alpine-specific features - revision numbers
		{
			name:        "revision_comparison_greater",
			constraints: []string{">=1.2.0-r5"},
			version:     "1.2.0-r7",
			expected:    true,
		},
		{
			name:        "revision_comparison_lesser",
			constraints: []string{">=1.2.0-r5"},
			version:     "1.2.0-r3",
			expected:    false,
		},
		{
			name:        "revision_range",
			constraints: []string{">=1.2.0-r1", "<=1.2.0-r10"},
			version:     "1.2.0-r5",
			expected:    true,
		},
		{
			name:        "revision_range_outside",
			constraints: []string{">=1.2.0-r1", "<=1.2.0-r10"},
			version:     "1.2.0-r15",
			expected:    false,
		},

		// Alpine suffixes
		{
			name:        "suffix_alpha",
			constraints: []string{">=1.2.0_alpha"},
			version:     "1.2.0_beta",
			expected:    true,
		},
		{
			name:        "suffix_alpha_false",
			constraints: []string{">1.2.0_beta"},
			version:     "1.2.0_alpha",
			expected:    false,
		},

		// Multiple constraints (AND logic)
		{
			name:        "multiple_constraints_and",
			constraints: []string{">=1.2.0", "<=2.0.0"},
			version:     "1.5.0",
			expected:    true,
		},
		{
			name:        "multiple_constraints_and_false",
			constraints: []string{">=1.2.0", "<=2.0.0"},
			version:     "2.5.0",
			expected:    false,
		},

		// Complex Alpine versions
		{
			name:        "complex_version_with_letter",
			constraints: []string{">=1.2.3a"},
			version:     "1.2.3b",
			expected:    true,
		},
		{
			name:        "complex_version_with_hash",
			constraints: []string{"=1.2.3~abc123"},
			version:     "1.2.3~abc123",
			expected:    true,
		},
		{
			name:        "complex_version_full",
			constraints: []string{">=1.2.3a_alpha1~abc123-r1"},
			version:     "1.2.3a_alpha1~abc123-r2",
			expected:    true,
		},

		// Error cases
		{
			name:        "invalid_version",
			constraints: []string{">=1.2.3"},
			version:     "not-a-version",
			expected:    false,
			expectError: true,
		},
		{
			name:        "invalid_constraint",
			constraints: []string{"invalid"},
			version:     "1.2.3",
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := alpineContains(tt.constraints, tt.version)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !tt.expectError && result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIntervalToAlpineRanges(t *testing.T) {
	tests := []struct {
		name     string
		interval interval
		expected []string
	}{
		{
			name: "exact_match",
			interval: interval{
				exact: "1.2.3",
			},
			expected: []string{"=1.2.3"},
		},
		{
			name: "lower_bound_inclusive",
			interval: interval{
				lower:          "1.2.0",
				lowerInclusive: true,
			},
			expected: []string{">=1.2.0"},
		},
		{
			name: "lower_bound_exclusive",
			interval: interval{
				lower:          "1.2.0",
				lowerInclusive: false,
			},
			expected: []string{">1.2.0"},
		},
		{
			name: "upper_bound_inclusive",
			interval: interval{
				upper:          "2.0.0",
				upperInclusive: true,
			},
			expected: []string{"<=2.0.0"},
		},
		{
			name: "upper_bound_exclusive",
			interval: interval{
				upper:          "2.0.0",
				upperInclusive: false,
			},
			expected: []string{"<2.0.0"},
		},
		{
			name: "both_bounds_inclusive",
			interval: interval{
				lower:          "1.2.0",
				lowerInclusive: true,
				upper:          "2.0.0",
				upperInclusive: true,
			},
			expected: []string{">=1.2.0 <=2.0.0"},
		},
		{
			name: "both_bounds_exclusive",
			interval: interval{
				lower:          "1.2.0",
				lowerInclusive: false,
				upper:          "2.0.0",
				upperInclusive: false,
			},
			expected: []string{">1.2.0 <2.0.0"},
		},
		{
			name: "exclude_handled_separately",
			interval: interval{
				exclude: "1.2.3",
			},
			expected: []string{},
		},
		{
			name:     "empty_interval",
			interval: interval{},
			expected: []string{},
		},
		{
			name: "alpine_revision_range",
			interval: interval{
				lower:          "1.2.0-r1",
				lowerInclusive: true,
				upper:          "1.2.0-r10",
				upperInclusive: true,
			},
			expected: []string{">=1.2.0-r1 <=1.2.0-r10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intervalToAlpineRanges(tt.interval)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d ranges, got %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if i >= len(result) || result[i] != expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
					return
				}
			}
		})
	}
}

// Test VERS integration
func TestVERSAlpineIntegration(t *testing.T) {
	tests := []struct {
		name     string
		versStr  string
		version  string
		expected bool
	}{
		{
			name:     "basic_range",
			versStr:  "vers:alpine/>=1.2.0",
			version:  "1.2.3",
			expected: true,
		},
		{
			name:     "revision_range",
			versStr:  "vers:alpine/>=1.2.0-r5",
			version:  "1.2.0-r7",
			expected: true,
		},
		{
			name:     "combined_constraints",
			versStr:  "vers:alpine/>=1.2.0|<=2.0.0",
			version:  "1.5.0",
			expected: true,
		},
		{
			name:     "exact_match",
			versStr:  "vers:alpine/=1.2.3-r5",
			version:  "1.2.3-r5",
			expected: true,
		},
		{
			name:     "exclusion",
			versStr:  "vers:alpine/!=1.2.3",
			version:  "1.2.4",
			expected: true,
		},
		{
			name:     "exclusion_false",
			versStr:  "vers:alpine/!=1.2.3",
			version:  "1.2.3",
			expected: false,
		},
		{
			name:     "star_constraint",
			versStr:  "vers:alpine/*",
			version:  "1.2.3-r5",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Contains(tt.versStr, tt.version)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for VERS %s with version %s",
					tt.expected, result, tt.versStr, tt.version)
			}
		})
	}
}
