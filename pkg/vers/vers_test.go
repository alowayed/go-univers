package vers

import (
	"testing"
)

func TestContains(t *testing.T) {
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
			name:      "invalid VERS format",
			versRange: "not-vers-format",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "unsupported ecosystem",
			versRange: "vers:unsupported/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
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
		// VERS validation edge cases
		{
			name:      "empty ecosystem name",
			versRange: "vers:/>1.0.0",
			version:   "1.1.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "empty constraints",
			versRange: "vers:maven/",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "missing colon",
			versRange: "versmaven/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "missing slash separator",
			versRange: "vers:maven>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "uppercase ecosystem should fail per VERS spec",
			versRange: "vers:MAVEN/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true, // VERS spec requires lowercase
		},
		{
			name:      "whitespace in constraints should be normalized",
			versRange: "vers:maven/ >= 1.0.0 | <= 2.0.0 ",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "multiple stars should fail",
			versRange: "vers:maven/*|*",
			version:   "1.0.0",
			want:      false,
			wantErr:   true, // VERS spec allows only one star
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
			name:      "multiple lower bounds - should take most restrictive",
			versRange: "vers:maven/>=1.0.0|>=2.0.0|<=3.0.0",
			version:   "1.5.0",
			want:      false, // Current implementation may fail this
			wantErr:   false,
		},
		{
			name:      "multiple upper bounds - should take most restrictive",
			versRange: "vers:maven/>=1.0.0|<=3.0.0|<=2.0.0",
			version:   "2.5.0",
			want:      false, // Current implementation may fail this
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
		// Additional VERS specification compliance tests
		{
			name:      "non-printable ASCII character should fail",
			versRange: "vers:maven/>=1.0.0\x00",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "tab character should fail",
			versRange: "vers:maven/>=1.0.0\t",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "newline character should fail",
			versRange: "vers:maven/>=1.0.0\n",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "star with other non-empty constraints should fail",
			versRange: "vers:maven/*|>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "star with whitespace-only constraints should pass",
			versRange: "vers:maven/*| | ",
			version:   "1.0.0",
			want:      true,
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
		{
			name:      "mixed case ecosystem should fail",
			versRange: "vers:Maven/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		// Additional ecosystem name validation tests
		{
			name:      "ecosystem with hyphen should fail",
			versRange: "vers:my-ecosystem/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "ecosystem with underscore should fail",
			versRange: "vers:my_ecosystem/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "ecosystem with period should fail",
			versRange: "vers:my.ecosystem/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "ecosystem with special characters should fail",
			versRange: "vers:my@ecosystem/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "ecosystem with numbers should pass validation but fail as unsupported",
			versRange: "vers:ecosystem2/>=1.0.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   true, // Should fail due to unsupported ecosystem, not validation
		},
		{
			name:      "ecosystem with only numbers should pass validation but fail as unsupported",
			versRange: "vers:123/>=1.0.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   true, // Should fail due to unsupported ecosystem, not validation
		},
		{
			name:      "ecosystem with mixed letters and numbers should pass validation but fail as unsupported",
			versRange: "vers:maven2/>=1.0.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   true, // Should fail due to unsupported ecosystem, not validation
		},
		// NPM VERS tests
		{
			name:      "npm simple range - contained",
			versRange: "vers:npm/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm simple range - not contained",
			versRange: "vers:npm/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "npm exact match",
			versRange: "vers:npm/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm exact match - not equal",
			versRange: "vers:npm/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "npm lower bound only",
			versRange: "vers:npm/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm upper bound only",
			versRange: "vers:npm/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm with prerelease versions",
			versRange: "vers:npm/>=1.0.0-alpha|<=2.0.0",
			version:   "1.5.0-beta",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm complex range - multiple intervals",
			versRange: "vers:npm/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm complex range - second interval",
			versRange: "vers:npm/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "2.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm complex range - outside intervals",
			versRange: "vers:npm/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "1.8.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "npm exclude constraint",
			versRange: "vers:npm/!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "npm exclude constraint - allowed version",
			versRange: "vers:npm/!=1.5.0",
			version:   "1.6.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm star constraint",
			versRange: "vers:npm/*",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm multiple lower bounds - should use most restrictive",
			versRange: "vers:npm/>=1.0.0|>=2.0.0|<=3.0.0",
			version:   "1.5.0",
			want:      false, // Should use most restrictive lower bound (>=2.0.0)
			wantErr:   false,
		},
		{
			name:      "npm multiple upper bounds - should use most restrictive",
			versRange: "vers:npm/>=1.0.0|<=3.0.0|<=2.0.0",
			version:   "2.5.0",
			want:      false, // Should use most restrictive upper bound (<=2.0.0)
			wantErr:   false,
		},
		{
			name:      "npm greater than constraint",
			versRange: "vers:npm/>1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm less than constraint",
			versRange: "vers:npm/<2.0.0",
			version:   "1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm invalid version in NPM",
			versRange: "vers:npm/>=1.0.0",
			version:   "invalid-npm-version",
			want:      false,
			wantErr:   true,
		},
		// PyPI VERS tests
		{
			name:      "pypi simple range - contained",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi simple range - not contained",
			versRange: "vers:pypi/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi exact match",
			versRange: "vers:pypi/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi exact match - not equal",
			versRange: "vers:pypi/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi lower bound only",
			versRange: "vers:pypi/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi upper bound only",
			versRange: "vers:pypi/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi with epoch versions",
			versRange: "vers:pypi/>=1!1.0.0|<=1!2.0.0",
			version:   "1!1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi with prerelease versions - alpha",
			versRange: "vers:pypi/>=1.0.0a1|<=2.0.0",
			version:   "1.5.0b1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi with prerelease versions - beta",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0b1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi with prerelease versions - rc",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0rc1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi with post-release versions",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0.post1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi with development release versions",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0.dev1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi with local version identifiers",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0+local.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi complex range - multiple intervals",
			versRange: "vers:pypi/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi complex range - second interval",
			versRange: "vers:pypi/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "2.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi complex range - outside intervals",
			versRange: "vers:pypi/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "1.8.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi exclude constraint",
			versRange: "vers:pypi/!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi exclude constraint - allowed version",
			versRange: "vers:pypi/!=1.5.0",
			version:   "1.6.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi star constraint",
			versRange: "vers:pypi/*",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi multiple lower bounds - should use most restrictive",
			versRange: "vers:pypi/>=1.0.0|>=2.0.0|<=3.0.0",
			version:   "1.5.0",
			want:      false, // Should use most restrictive lower bound (>=2.0.0)
			wantErr:   false,
		},
		{
			name:      "pypi multiple upper bounds - should use most restrictive",
			versRange: "vers:pypi/>=1.0.0|<=3.0.0|<=2.0.0",
			version:   "2.5.0",
			want:      false, // Should use most restrictive upper bound (<=2.0.0)
			wantErr:   false,
		},
		{
			name:      "pypi greater than constraint",
			versRange: "vers:pypi/>1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi less than constraint",
			versRange: "vers:pypi/<2.0.0",
			version:   "1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi epoch precedence test",
			versRange: "vers:pypi/>=1!1.0.0",
			version:   "2.0.0",
			want:      false, // Epoch 1 > Epoch 0 (default)
			wantErr:   false,
		},
		{
			name:      "pypi epoch precedence test - satisfied",
			versRange: "vers:pypi/>=1!1.0.0",
			version:   "1!1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi prerelease vs release precedence",
			versRange: "vers:pypi/>=1.0.0",
			version:   "1.0.0a1",
			want:      false, // Prerelease < release
			wantErr:   false,
		},
		{
			name:      "pypi prerelease vs release precedence - satisfied",
			versRange: "vers:pypi/>=1.0.0a1",
			version:   "1.0.0",
			want:      true, // Release > prerelease
			wantErr:   false,
		},
		{
			name:      "pypi invalid version",
			versRange: "vers:pypi/>=1.0.0",
			version:   "invalid-pypi-version",
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
