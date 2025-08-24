package vers

import (
	"testing"
)

// TestContains_Alpine tests VERS constraint checking for Alpine ecosystem
func TestContains_Alpine(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		// Basic comparison operators
		{
			name:      "exact_match",
			versRange: "vers:alpine/=1.2.3",
			version:   "1.2.3",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "exact_match_false",
			versRange: "vers:alpine/=1.2.3",
			version:   "1.2.4",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "greater_than_equal",
			versRange: "vers:alpine/>=1.2.0",
			version:   "1.2.3",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "greater_than_equal_false",
			versRange: "vers:alpine/>=1.2.3",
			version:   "1.2.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "less_than_equal",
			versRange: "vers:alpine/<=1.2.3",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "less_than_equal_false",
			versRange: "vers:alpine/<=1.2.0",
			version:   "1.2.3",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "greater_than",
			versRange: "vers:alpine/>1.2.0",
			version:   "1.2.3",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "greater_than_false",
			versRange: "vers:alpine/>1.2.3",
			version:   "1.2.3",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "less_than",
			versRange: "vers:alpine/<1.2.3",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "less_than_false",
			versRange: "vers:alpine/<1.2.0",
			version:   "1.2.3",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "not_equal",
			versRange: "vers:alpine/!=1.2.0",
			version:   "1.2.3",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "not_equal_false",
			versRange: "vers:alpine/!=1.2.3",
			version:   "1.2.3",
			want:      false,
			wantErr:   false,
		},

		// Alpine-specific features - revision numbers
		{
			name:      "revision_comparison_greater",
			versRange: "vers:alpine/>=1.2.0-r5",
			version:   "1.2.0-r7",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "revision_comparison_lesser",
			versRange: "vers:alpine/>=1.2.0-r5",
			version:   "1.2.0-r3",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "revision_range",
			versRange: "vers:alpine/>=1.2.0-r1|<=1.2.0-r10",
			version:   "1.2.0-r5",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "revision_range_outside",
			versRange: "vers:alpine/>=1.2.0-r1|<=1.2.0-r10",
			version:   "1.2.0-r15",
			want:      false,
			wantErr:   false,
		},

		// Alpine suffixes
		{
			name:      "suffix_alpha_to_beta",
			versRange: "vers:alpine/>=1.2.0_alpha",
			version:   "1.2.0_beta",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "suffix_beta_to_alpha",
			versRange: "vers:alpine/>1.2.0_beta",
			version:   "1.2.0_alpha",
			want:      false,
			wantErr:   false,
		},

		// Multiple constraints (range)
		{
			name:      "multiple_constraints_and_within",
			versRange: "vers:alpine/>=1.2.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "multiple_constraints_and_outside",
			versRange: "vers:alpine/>=1.2.0|<=2.0.0",
			version:   "2.5.0",
			want:      false,
			wantErr:   false,
		},

		// Complex Alpine versions
		{
			name:      "complex_version_with_letter",
			versRange: "vers:alpine/>=1.2.3a",
			version:   "1.2.3b",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "complex_version_with_hash",
			versRange: "vers:alpine/=1.2.3~abc123",
			version:   "1.2.3~abc123",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "complex_version_full",
			versRange: "vers:alpine/>=1.2.3a_alpha1~abc123-r1",
			version:   "1.2.3a_alpha1~abc123-r2",
			want:      true,
			wantErr:   false,
		},

		// Star constraint (all versions)
		{
			name:      "star_constraint_matches_all",
			versRange: "vers:alpine/*",
			version:   "1.2.3-r5",
			want:      true,
			wantErr:   false,
		},

		// Error cases
		{
			name:      "invalid_version",
			versRange: "vers:alpine/>=1.2.3",
			version:   "not-a-version",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "invalid_constraint_version",
			versRange: "vers:alpine/>=invalid-version",
			version:   "1.2.3",
			want:      false,
			wantErr:   true,
		},

		// Edge cases
		{
			name:      "no_revision_vs_revision",
			versRange: "vers:alpine/>=1.2.3",
			version:   "1.2.3-r1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "revision_vs_no_revision",
			versRange: "vers:alpine/>=1.2.3-r1",
			version:   "1.2.3",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "suffix_vs_no_suffix",
			versRange: "vers:alpine/>=1.2.3_alpha",
			version:   "1.2.3",
			want:      true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Contains(tt.versRange, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
