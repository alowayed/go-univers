package vers

import (
	"testing"
)

// TestContains_SemVer tests VERS functionality specifically for the generic SemVer ecosystem
func TestContains_SemVer(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		{
			name:      "semver simple range - contained",
			versRange: "vers:semver/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver simple range - not contained",
			versRange: "vers:semver/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "semver exact match",
			versRange: "vers:semver/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver exact match - not equal",
			versRange: "vers:semver/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "semver lower bound only",
			versRange: "vers:semver/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver upper bound only",
			versRange: "vers:semver/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver greater than",
			versRange: "vers:semver/>1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver greater than - not satisfied",
			versRange: "vers:semver/>1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "semver less than",
			versRange: "vers:semver/<2.0.0",
			version:   "1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver less than - not satisfied",
			versRange: "vers:semver/<2.0.0",
			version:   "2.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "semver not equal - satisfied",
			versRange: "vers:semver/!=1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver not equal - not satisfied",
			versRange: "vers:semver/!=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "semver multiple constraints - AND logic",
			versRange: "vers:semver/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver multiple constraints - excluded",
			versRange: "vers:semver/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "semver prerelease version",
			versRange: "vers:semver/>=1.0.0-alpha.1",
			version:   "1.0.0-alpha.2",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver prerelease vs release",
			versRange: "vers:semver/>=1.0.0",
			version:   "1.0.0-alpha.1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "semver release vs prerelease",
			versRange: "vers:semver/>=1.0.0-alpha.1",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver complex prerelease",
			versRange: "vers:semver/>=1.0.0-alpha.1|<2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver complex prerelease - contained",
			versRange: "vers:semver/>=1.0.0-alpha.1|<2.0.0",
			version:   "1.0.0-alpha.5",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver complex prerelease - not contained",
			versRange: "vers:semver/>=1.0.0-alpha.5|<2.0.0",
			version:   "1.0.0-alpha.1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "semver with build metadata - equal",
			versRange: "vers:semver/>=1.0.0+build123",
			version:   "1.0.0+build456",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver with build metadata - contained",
			versRange: "vers:semver/>=1.0.0+build123|<2.0.0",
			version:   "1.5.0+other",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver prerelease with build metadata",
			versRange: "vers:semver/>=1.0.0-alpha.1+build|<2.0.0",
			version:   "1.0.0-beta.1+another",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver comprehensive range",
			versRange: "vers:semver/>=1.0.0-alpha|<=2.0.0-beta+build",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver star constraint - matches all",
			versRange: "vers:semver/*",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver star constraint - matches prerelease",
			versRange: "vers:semver/*",
			version:   "1.0.0-alpha.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "semver star constraint - matches build metadata",
			versRange: "vers:semver/*",
			version:   "1.0.0+build123",
			want:      true,
			wantErr:   false,
		},
		// Error cases
		{
			name:      "semver invalid version",
			versRange: "vers:semver/>=1.0.0",
			version:   "invalid-version",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "semver invalid constraint version",
			versRange: "vers:semver/>=invalid",
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
