package vers

import (
	"testing"
)

// TestContains_Gem tests VERS functionality specifically for the RubyGems ecosystem
func TestContains_Gem(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		{
			name:      "gem simple range - contained",
			versRange: "vers:gem/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem simple range - not contained",
			versRange: "vers:gem/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "gem exact match",
			versRange: "vers:gem/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem exact match - not equal",
			versRange: "vers:gem/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "gem lower bound only",
			versRange: "vers:gem/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem upper bound only",
			versRange: "vers:gem/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem greater than",
			versRange: "vers:gem/>1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem greater than - not satisfied",
			versRange: "vers:gem/>1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "gem less than",
			versRange: "vers:gem/<2.0.0",
			version:   "1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem less than - not satisfied",
			versRange: "vers:gem/<2.0.0",
			version:   "2.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "gem not equal - satisfied",
			versRange: "vers:gem/!=1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem not equal - not satisfied",
			versRange: "vers:gem/!=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "gem multiple constraints - AND logic",
			versRange: "vers:gem/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem multiple constraints - excluded",
			versRange: "vers:gem/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "gem prerelease version",
			versRange: "vers:gem/>=1.0.0-alpha",
			version:   "1.0.0-beta",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem prerelease vs release",
			versRange: "vers:gem/>=1.0.0",
			version:   "1.0.0-alpha",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "gem release vs prerelease",
			versRange: "vers:gem/>=1.0.0-alpha",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem complex prerelease",
			versRange: "vers:gem/>=1.0.0.pre1|<2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem complex prerelease - contained",
			versRange: "vers:gem/>=1.0.0.pre1|<2.0.0",
			version:   "1.0.0.pre2",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem complex prerelease - not contained",
			versRange: "vers:gem/>=1.0.0.pre1|<2.0.0",
			version:   "1.0.0.alpha",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "gem star constraint - matches all",
			versRange: "vers:gem/*",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "gem star constraint - matches prerelease",
			versRange: "vers:gem/*",
			version:   "1.0.0-alpha",
			want:      true,
			wantErr:   false,
		},
		// Error cases
		{
			name:      "gem invalid version",
			versRange: "vers:gem/>=1.0.0",
			version:   "invalid-version",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "gem invalid constraint version",
			versRange: "vers:gem/>=invalid",
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
