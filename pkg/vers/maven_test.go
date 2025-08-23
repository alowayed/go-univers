package vers

import (
	"testing"
)

// TestContains_Maven tests VERS functionality specifically for the Maven ecosystem
func TestContains_Maven(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		{
			name:      "maven simple range - contained",
			versRange: "vers:maven/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven simple range - not contained",
			versRange: "vers:maven/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "maven exact match",
			versRange: "vers:maven/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven exact match - not equal",
			versRange: "vers:maven/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "maven lower bound only",
			versRange: "vers:maven/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven upper bound only",
			versRange: "vers:maven/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven complex range from issue",
			versRange: "vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7",
			version:   "1.1.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven complex range - second interval",
			versRange: "vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7",
			version:   "7.0.5",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven complex range - outside intervals",
			versRange: "vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7",
			version:   "5.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "invalid version",
			versRange: "vers:maven/>=1.0.0",
			version:   "invalid-version",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "invalid constraint format",
			versRange: "vers:maven/~1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "maven unordered constraints - should handle correctly",
			versRange: "vers:maven/>=2.0.0|>=1.0.0|<=3.0.0",
			version:   "2.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven unordered constraints - outside range",
			versRange: "vers:maven/>=2.0.0|>=1.0.0|<=3.0.0",
			version:   "0.5.0",
			want:      false, // Fixed: Multiple lower bounds with single upper should merge to most restrictive (>=2.0.0, <=3.0.0)
			wantErr:   false,
		},
		{
			name:      "maven redundant lower bounds - should use more restrictive",
			versRange: "vers:maven/>=1.0.0|>=2.0.0|<=3.0.0",
			version:   "1.5.0",
			want:      false, // Fixed: Should use most restrictive lower bound (>=2.0.0), so 1.5.0 is excluded
			wantErr:   false,
		},
		{
			name:      "maven redundant upper bounds - should use more restrictive",
			versRange: "vers:maven/>=1.0.0|<=3.0.0|<=2.0.0",
			version:   "2.5.0",
			want:      false, // Fixed: Should use most restrictive upper bound (<=2.0.0), so 2.5.0 is excluded
			wantErr:   false,
		},
		// Test exclude functionality
		{
			name:      "maven exclude - version should be excluded",
			versRange: "vers:maven/!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "maven exclude - version below excluded should be included",
			versRange: "vers:maven/!=1.5.0",
			version:   "1.4.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven exclude - version above excluded should be included",
			versRange: "vers:maven/!=1.5.0",
			version:   "1.6.0",
			want:      true,
			wantErr:   false,
		},
		// Test star wildcard
		{
			name:      "maven star wildcard - should match any version",
			versRange: "vers:maven/*",
			version:   "1.2.3",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven star wildcard - should match any version including prerelease",
			versRange: "vers:maven/*",
			version:   "2.0.0-SNAPSHOT",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "whitespace in constraints should be normalized",
			versRange: "vers:maven/ >= 1.0.0 | <= 2.0.0 ",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "duplicate constraints should be normalized",
			versRange: "vers:maven/>=1.0.0|>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "unsorted constraints should work (normalization)",
			versRange: "vers:maven/<=2.0.0|>=1.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		// Complex constraint grouping scenarios
		{
			name:      "conflicting constraints - impossible range",
			versRange: "vers:maven/>=2.0.0|<=1.0.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "mixed inclusive and exclusive bounds",
			versRange: "vers:maven/>1.0.0|<2.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "mixed inclusive and exclusive bounds - edge case",
			versRange: "vers:maven/>1.0.0|<2.0.0",
			version:   "2.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "three intervals scenario - gap between intervals",
			versRange: "vers:maven/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0|>=3.0.0|<=3.5.0",
			version:   "1.8.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "three intervals scenario - in first interval",
			versRange: "vers:maven/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0|>=3.0.0|<=3.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "exclude with range combination",
			versRange: "vers:maven/>=1.0.0|<=3.0.0|!=2.0.0",
			version:   "2.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "exclude with range combination - allowed version",
			versRange: "vers:maven/>=1.0.0|<=3.0.0|!=2.0.0",
			version:   "1.9.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "multiple excludes",
			versRange: "vers:maven/!=1.0.0|!=2.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "multiple excludes - allowed version",
			versRange: "vers:maven/!=1.0.0|!=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "redundant constraints with different operators",
			versRange: "vers:maven/>=1.0.0|>1.0.0|<=2.0.0",
			version:   "1.0.0",
			want:      false, // Fixed: Most restrictive lower bound is >1.0.0, so 1.0.0 is excluded
			wantErr:   false,
		},
		{
			name:      "constraints with only spaces should be normalized out",
			versRange: "vers:maven/>=1.0.0|   |<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "malformed constraint without operator should fail",
			versRange: "vers:maven/1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "constraint with operator but no version should fail",
			versRange: "vers:maven/>=",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
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
