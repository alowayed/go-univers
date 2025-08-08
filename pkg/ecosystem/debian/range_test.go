package debian

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
		wantErr  bool
	}{
		// Basic ranges
		{
			name:     "exact match",
			rangeStr: "1.0",
			wantErr:  false,
		},
		{
			name:     "greater than",
			rangeStr: "> 1.0",
			wantErr:  false,
		},
		{
			name:     "greater than or equal",
			rangeStr: ">= 1.0",
			wantErr:  false,
		},
		{
			name:     "less than",
			rangeStr: "< 2.0",
			wantErr:  false,
		},
		{
			name:     "less than or equal",
			rangeStr: "<= 2.0",
			wantErr:  false,
		},
		{
			name:     "not equal",
			rangeStr: "!= 1.5",
			wantErr:  false,
		},

		// Debian-specific operators
		{
			name:     "much greater than",
			rangeStr: ">> 1.0",
			wantErr:  false,
		},
		{
			name:     "much less than",
			rangeStr: "<< 2.0",
			wantErr:  false,
		},

		// Multiple constraints
		{
			name:     "range with AND",
			rangeStr: ">= 1.0, < 2.0",
			wantErr:  false,
		},
		{
			name:     "complex range",
			rangeStr: ">= 1.0, < 2.0, != 1.5",
			wantErr:  false,
		},

		// Complex versions in ranges
		{
			name:     "range with epoch",
			rangeStr: ">= 1:1.0",
			wantErr:  false,
		},
		{
			name:     "range with revision",
			rangeStr: ">= 1.0-1",
			wantErr:  false,
		},

		// Error cases
		{
			name:     "empty range",
			rangeStr: "",
			wantErr:  true,
		},
		{
			name:     "whitespace only",
			rangeStr: "   ",
			wantErr:  true,
		},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ecosystem.NewVersionRange(tt.rangeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecosystem.NewVersionRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.rangeStr {
				t.Errorf("Ecosystem.NewVersionRange().String() = %v, want %v", got.String(), tt.rangeStr)
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	ecosystem := &Ecosystem{}

	tests := []struct {
		name     string
		rangeStr string
		version  string
		expected bool
	}{
		// Basic exact matches
		{
			name:     "exact match true",
			rangeStr: "1.0",
			version:  "1.0",
			expected: true,
		},
		{
			name:     "exact match false",
			rangeStr: "1.0",
			version:  "1.1",
			expected: false,
		},

		// Greater than
		{
			name:     "greater than true",
			rangeStr: "> 1.0",
			version:  "1.1",
			expected: true,
		},
		{
			name:     "greater than false",
			rangeStr: "> 1.0",
			version:  "1.0",
			expected: false,
		},

		// Greater than or equal
		{
			name:     "greater than or equal true (greater)",
			rangeStr: ">= 1.0",
			version:  "1.1",
			expected: true,
		},
		{
			name:     "greater than or equal true (equal)",
			rangeStr: ">= 1.0",
			version:  "1.0",
			expected: true,
		},
		{
			name:     "greater than or equal false",
			rangeStr: ">= 1.0",
			version:  "0.9",
			expected: false,
		},

		// Less than
		{
			name:     "less than true",
			rangeStr: "< 2.0",
			version:  "1.9",
			expected: true,
		},
		{
			name:     "less than false",
			rangeStr: "< 2.0",
			version:  "2.0",
			expected: false,
		},

		// Less than or equal
		{
			name:     "less than or equal true (less)",
			rangeStr: "<= 2.0",
			version:  "1.9",
			expected: true,
		},
		{
			name:     "less than or equal true (equal)",
			rangeStr: "<= 2.0",
			version:  "2.0",
			expected: true,
		},
		{
			name:     "less than or equal false",
			rangeStr: "<= 2.0",
			version:  "2.1",
			expected: false,
		},

		// Not equal
		{
			name:     "not equal true",
			rangeStr: "!= 1.5",
			version:  "1.4",
			expected: true,
		},
		{
			name:     "not equal false",
			rangeStr: "!= 1.5",
			version:  "1.5",
			expected: false,
		},

		// Debian-specific operators
		{
			name:     "much greater than true",
			rangeStr: ">> 1.0",
			version:  "1.1",
			expected: true,
		},
		{
			name:     "much greater than false",
			rangeStr: ">> 1.0",
			version:  "1.0",
			expected: false,
		},
		{
			name:     "much less than true",
			rangeStr: "<< 2.0",
			version:  "1.9",
			expected: true,
		},
		{
			name:     "much less than false",
			rangeStr: "<< 2.0",
			version:  "2.0",
			expected: false,
		},

		// Multiple constraints (AND logic)
		{
			name:     "range satisfied",
			rangeStr: ">= 1.0, < 2.0",
			version:  "1.5",
			expected: true,
		},
		{
			name:     "range not satisfied (too low)",
			rangeStr: ">= 1.0, < 2.0",
			version:  "0.9",
			expected: false,
		},
		{
			name:     "range not satisfied (too high)",
			rangeStr: ">= 1.0, < 2.0",
			version:  "2.0",
			expected: false,
		},

		// Complex versions
		{
			name:     "epoch comparison in range",
			rangeStr: ">= 1:1.0",
			version:  "1:1.5",
			expected: true,
		},
		{
			name:     "epoch comparison in range false",
			rangeStr: ">= 1:1.0",
			version:  "0:2.0",
			expected: false,
		},
		{
			name:     "revision comparison in range",
			rangeStr: ">= 1.0-1",
			version:  "1.0-2",
			expected: true,
		},

		// Tilde handling in ranges
		{
			name:     "tilde version in range",
			rangeStr: ">= 1.0~beta",
			version:  "1.0",
			expected: true,
		},
		{
			name:     "tilde version not in range",
			rangeStr: "> 1.0",
			version:  "1.0~beta",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr, err := ecosystem.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("Failed to parse range %s: %v", tt.rangeStr, err)
			}

			v, err := ecosystem.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt.version, err)
			}

			result := vr.Contains(v)
			if result != tt.expected {
				t.Errorf("VersionRange.Contains(%s in %s) = %v, want %v", tt.version, tt.rangeStr, result, tt.expected)
			}
		})
	}
}

func TestVersionRange_String(t *testing.T) {
	ecosystem := &Ecosystem{}

	tests := []string{
		"1.0",
		">= 1.0",
		"< 2.0",
		">= 1.0, < 2.0",
		">> 1:1.0-1",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			vr, err := ecosystem.NewVersionRange(tt)
			if err != nil {
				t.Fatalf("Failed to parse range %s: %v", tt, err)
			}

			if vr.String() != tt {
				t.Errorf("VersionRange.String() = %v, want %v", vr.String(), tt)
			}
		})
	}
}
