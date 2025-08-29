package vers

import (
	"testing"
)

// TestContains_Debian tests VERS functionality specifically for the Debian ecosystem
func TestContains_Debian(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		{
			name:      "debian simple range - contained",
			versRange: "vers:debian/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian simple range - not contained",
			versRange: "vers:debian/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian exact match",
			versRange: "vers:debian/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian exact match - not equal",
			versRange: "vers:debian/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian lower bound only",
			versRange: "vers:debian/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian upper bound only",
			versRange: "vers:debian/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian greater than",
			versRange: "vers:debian/>1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian greater than - not satisfied",
			versRange: "vers:debian/>1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian less than",
			versRange: "vers:debian/<2.0.0",
			version:   "1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian less than - not satisfied",
			versRange: "vers:debian/<2.0.0",
			version:   "2.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian not equal - satisfied",
			versRange: "vers:debian/!=1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian not equal - not satisfied",
			versRange: "vers:debian/!=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian multiple constraints - AND logic",
			versRange: "vers:debian/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian multiple constraints - excluded",
			versRange: "vers:debian/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian epoch version - contained",
			versRange: "vers:debian/>=1:1.0.0",
			version:   "1:2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian epoch version - not contained",
			versRange: "vers:debian/>=2:1.0.0",
			version:   "1:2.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian epoch version range",
			versRange: "vers:debian/>=1:2.4.1|<1:3.0.0",
			version:   "1:2.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian epoch version range - not contained",
			versRange: "vers:debian/>=2:1.0.0|<=2:2.0.0",
			version:   "1:5.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian revision version - contained",
			versRange: "vers:debian/>=2.4.1-1",
			version:   "2.4.1-2",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian revision version - not contained",
			versRange: "vers:debian/>=2.4.2-1",
			version:   "2.4.1-9",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian ubuntu revision version",
			versRange: "vers:debian/>=2.4.1-1ubuntu1|<2.4.2",
			version:   "2.4.1-1ubuntu2",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian ubuntu revision version - not contained",
			versRange: "vers:debian/>=2.4.1-1ubuntu1|<2.4.2",
			version:   "2.4.2-1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian epoch and revision version",
			versRange: "vers:debian/>=1:2.4.1-1|<1:3.0.0-1",
			version:   "1:2.5.0-2",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian epoch and revision version - not contained",
			versRange: "vers:debian/>=1:2.4.1-1|<1:3.0.0-1",
			version:   "1:4.0.0-1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian complex epoch range",
			versRange: "vers:debian/>=1:2.4.1-1|<=1:3.0.0-1|!=1:2.5.0-1",
			version:   "1:2.6.0-1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian complex epoch range - excluded",
			versRange: "vers:debian/>=1:2.4.1-1|<=1:3.0.0-1|!=1:2.5.0-1",
			version:   "1:2.5.0-1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "debian star constraint - matches all",
			versRange: "vers:debian/*",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian star constraint - matches epoch",
			versRange: "vers:debian/*",
			version:   "1:2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian star constraint - matches revision",
			versRange: "vers:debian/*",
			version:   "2.0.0-1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "debian star constraint - matches epoch and revision",
			versRange: "vers:debian/*",
			version:   "1:2.0.0-1ubuntu1",
			want:      true,
			wantErr:   false,
		},
		// Error cases
		{
			name:      "debian invalid version",
			versRange: "vers:debian/>=1.0.0",
			version:   "1.0.0@invalid",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "debian invalid constraint version",
			versRange: "vers:debian/>=1.0.0@invalid",
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
