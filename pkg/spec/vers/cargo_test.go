package vers

import (
	"testing"
)

// TestContains_Cargo tests VERS functionality specifically for the Cargo (Rust) ecosystem
func TestContains_Cargo(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		{
			name:      "cargo simple range - contained",
			versRange: "vers:cargo/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo simple range - not contained",
			versRange: "vers:cargo/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "cargo exact match",
			versRange: "vers:cargo/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo exact match - not equal",
			versRange: "vers:cargo/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "cargo lower bound only",
			versRange: "vers:cargo/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo upper bound only",
			versRange: "vers:cargo/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo greater than",
			versRange: "vers:cargo/>1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo greater than - not satisfied",
			versRange: "vers:cargo/>1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "cargo less than",
			versRange: "vers:cargo/<2.0.0",
			version:   "1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo less than - not satisfied",
			versRange: "vers:cargo/<2.0.0",
			version:   "2.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "cargo not equal - satisfied",
			versRange: "vers:cargo/!=1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo not equal - not satisfied",
			versRange: "vers:cargo/!=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "cargo multiple constraints - AND logic",
			versRange: "vers:cargo/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo multiple constraints - excluded",
			versRange: "vers:cargo/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "cargo prerelease version",
			versRange: "vers:cargo/>=1.0.0-alpha.1",
			version:   "1.0.0-alpha.2",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo prerelease vs release",
			versRange: "vers:cargo/>=1.0.0",
			version:   "1.0.0-alpha.1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "cargo release vs prerelease",
			versRange: "vers:cargo/>=1.0.0-alpha.1",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo complex prerelease",
			versRange: "vers:cargo/>=1.0.0-alpha.1|<2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo complex prerelease - contained",
			versRange: "vers:cargo/>=1.0.0-alpha.1|<2.0.0",
			version:   "1.0.0-alpha.5",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo complex prerelease - not contained",
			versRange: "vers:cargo/>=1.0.0-alpha.5|<2.0.0",
			version:   "1.0.0-alpha.1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "cargo with build metadata - equal",
			versRange: "vers:cargo/>=1.0.0+build123",
			version:   "1.0.0+build456",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo with build metadata - contained",
			versRange: "vers:cargo/>=1.0.0+build123|<2.0.0",
			version:   "1.5.0+other",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo prerelease with build metadata",
			versRange: "vers:cargo/>=1.0.0-alpha.1+build|<2.0.0",
			version:   "1.0.0-beta.1+another",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo star constraint - matches all",
			versRange: "vers:cargo/*",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo star constraint - matches prerelease",
			versRange: "vers:cargo/*",
			version:   "1.0.0-alpha.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "cargo star constraint - matches build metadata",
			versRange: "vers:cargo/*",
			version:   "1.0.0+build123",
			want:      true,
			wantErr:   false,
		},
		// Error cases
		{
			name:      "cargo invalid version",
			versRange: "vers:cargo/>=1.0.0",
			version:   "invalid-version",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "cargo invalid constraint version",
			versRange: "vers:cargo/>=invalid",
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
