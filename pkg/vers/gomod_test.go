package vers

import (
	"testing"
)

// TestContains_Gomod tests VERS constraint checking for Go modules ecosystem using 'golang' scheme
func TestContains_Gomod(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		// Basic v-prefixed versions
		{
			name:      "greater_than_equal_v_prefixed",
			versRange: "vers:golang/>=v1.2.3",
			version:   "v1.2.4",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "less_than_equal_v_prefixed",
			versRange: "vers:golang/<=v2.0.0",
			version:   "v1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "exact_match_v_prefixed",
			versRange: "vers:golang/=v1.2.3",
			version:   "v1.2.3",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "not_equal_v_prefixed",
			versRange: "vers:golang/!=v1.2.3",
			version:   "v1.2.4",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "not_equal_match_v_prefixed",
			versRange: "vers:golang/!=v1.2.3",
			version:   "v1.2.3",
			want:      false,
			wantErr:   false,
		},
		// Version without v-prefix (should work)
		{
			name:      "greater_than_no_v_prefix",
			versRange: "vers:golang/>1.2.3",
			version:   "1.2.4",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "mixed_v_prefix_constraints",
			versRange: "vers:golang/>=v1.2.3|<=2.0.0",
			version:   "v1.5.0",
			want:      true,
			wantErr:   false,
		},
		// Pre-release versions
		{
			name:      "prerelease_greater_than",
			versRange: "vers:golang/>v1.2.3-alpha",
			version:   "v1.2.3-beta",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "prerelease_less_than_release",
			versRange: "vers:golang/<v1.2.3",
			version:   "v1.2.3-alpha",
			want:      true,
			wantErr:   false,
		},
		// Pseudo-versions
		{
			name:      "pseudo_version_pattern1",
			versRange: "vers:golang/>=v0.0.0-20170915032832-14c0d48ead0c",
			version:   "v0.0.0-20170915032833-14c0d48ead0c",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pseudo_version_pattern2",
			versRange: "vers:golang/>=v1.2.3-pre.0.20170915032832-14c0d48ead0c",
			version:   "v1.2.3-pre.0.20170915032833-14c0d48ead0c",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pseudo_version_pattern3",
			versRange: "vers:golang/>=v1.2.4-0.20170915032832-14c0d48ead0c",
			version:   "v1.2.4-0.20170915032833-14c0d48ead0c",
			want:      true,
			wantErr:   false,
		},
		// Build metadata (should be ignored in comparison)
		{
			name:      "build_metadata_ignored",
			versRange: "vers:golang/=v1.2.3+meta",
			version:   "v1.2.3+different",
			want:      true,
			wantErr:   false,
		},
		// Range constraints (multiple operators)
		{
			name:      "range_within_bounds",
			versRange: "vers:golang/>=v1.2.3|<=v2.0.0",
			version:   "v1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "range_outside_lower_bound",
			versRange: "vers:golang/>=v1.2.3|<=v2.0.0",
			version:   "v1.2.2",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "range_outside_upper_bound",
			versRange: "vers:golang/>=v1.2.3|<=v2.0.0",
			version:   "v2.0.1",
			want:      false,
			wantErr:   false,
		},
		// Major version boundaries
		{
			name:      "v0_unstable_version",
			versRange: "vers:golang/>=v0.1.0|<=v0.9.0",
			version:   "v0.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "v1_stable_version",
			versRange: "vers:golang/>=v1.0.0|<v2.0.0",
			version:   "v1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "v2_major_version",
			versRange: "vers:golang/>=v2.0.0",
			version:   "v2.1.0",
			want:      true,
			wantErr:   false,
		},
		// Complex range expressions
		{
			name:      "multiple_constraints_all_match",
			versRange: "vers:golang/>=v1.0.0|<=v2.0.0|!=v1.5.0",
			version:   "v1.4.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "multiple_constraints_exclusion_match",
			versRange: "vers:golang/>=v1.0.0|<=v2.0.0|!=v1.5.0",
			version:   "v1.5.0",
			want:      false,
			wantErr:   false,
		},
		// Star constraint (all versions)
		{
			name:      "star_constraint_matches_all",
			versRange: "vers:golang/*",
			version:   "v1.2.3",
			want:      true,
			wantErr:   false,
		},
		// Error cases
		{
			name:      "invalid_version_format",
			versRange: "vers:golang/>=v1.2.3",
			version:   "invalid-version",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "invalid_constraint_version",
			versRange: "vers:golang/>=invalid",
			version:   "v1.2.3",
			want:      false,
			wantErr:   true,
		},
		// Edge cases for Go-specific patterns
		{
			name:      "zero_major_version",
			versRange: "vers:golang/>=v0.0.1",
			version:   "v0.0.2",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "large_major_version",
			versRange: "vers:golang/>=v10.0.0",
			version:   "v10.0.1",
			want:      true,
			wantErr:   false,
		},
		// v-prefix handling scenarios
		{
			name:      "constraint_with_v_version_without_v",
			versRange: "vers:golang/>=v1.2.3",
			version:   "1.2.4",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "constraint_without_v_version_with_v",
			versRange: "vers:golang/>=1.2.3",
			version:   "v1.2.4",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "both_without_v_prefix",
			versRange: "vers:golang/>=1.2.3",
			version:   "1.2.4",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "both_with_v_prefix",
			versRange: "vers:golang/>=v1.2.3",
			version:   "v1.2.4",
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
