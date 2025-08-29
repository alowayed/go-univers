package vers

import (
	"testing"
)

// TestContains_RPM tests VERS functionality specifically for the RPM ecosystem
func TestContains_RPM(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		{
			name:      "rpm simple range - contained",
			versRange: "vers:rpm/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm simple range - not contained",
			versRange: "vers:rpm/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm exact match",
			versRange: "vers:rpm/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm exact match - not equal",
			versRange: "vers:rpm/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm lower bound only",
			versRange: "vers:rpm/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm upper bound only",
			versRange: "vers:rpm/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm greater than",
			versRange: "vers:rpm/>1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm greater than - not satisfied",
			versRange: "vers:rpm/>1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm less than",
			versRange: "vers:rpm/<2.0.0",
			version:   "1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm less than - not satisfied",
			versRange: "vers:rpm/<2.0.0",
			version:   "2.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm not equal - satisfied",
			versRange: "vers:rpm/!=1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm not equal - not satisfied",
			versRange: "vers:rpm/!=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm multiple constraints - AND logic",
			versRange: "vers:rpm/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm multiple constraints - excluded",
			versRange: "vers:rpm/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm epoch version - contained",
			versRange: "vers:rpm/>=1:1.0.0",
			version:   "1:2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm epoch version - not contained",
			versRange: "vers:rpm/>=2:1.0.0",
			version:   "1:2.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm epoch version range",
			versRange: "vers:rpm/>=1:2.4.1|<1:3.0.0",
			version:   "1:2.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm epoch version range - not contained",
			versRange: "vers:rpm/>=2:1.0.0|<=2:2.0.0",
			version:   "1:5.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm release version - contained",
			versRange: "vers:rpm/>=2.4.1-1.el8",
			version:   "2.4.1-1.el9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm release version - not contained",
			versRange: "vers:rpm/>=2.4.2-1.el8",
			version:   "2.4.1-1.el9",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm epoch and release version",
			versRange: "vers:rpm/>=1:2.4.1-1.el8|<1:3.0.0-1.el9",
			version:   "1:2.5.0-1.el8",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm epoch and release version - not contained",
			versRange: "vers:rpm/>=1:2.4.1-1.el8|<1:3.0.0-1.el9",
			version:   "1:4.0.0-1.el9",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm complex epoch range",
			versRange: "vers:rpm/>=1:2.4.1-1.el8|<=1:3.0.0-1.el9|!=1:2.5.0-1.el8",
			version:   "1:2.6.0-1.el8",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm complex epoch range - excluded",
			versRange: "vers:rpm/>=1:2.4.1-1.el8|<=1:3.0.0-1.el9|!=1:2.5.0-1.el8",
			version:   "1:2.5.0-1.el8",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "rpm star constraint - matches all",
			versRange: "vers:rpm/*",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm star constraint - matches epoch",
			versRange: "vers:rpm/*",
			version:   "1:2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm star constraint - matches release",
			versRange: "vers:rpm/*",
			version:   "2.0.0-1.el8",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "rpm star constraint - matches epoch and release",
			versRange: "vers:rpm/*",
			version:   "1:2.0.0-1.el8",
			want:      true,
			wantErr:   false,
		},
		// Error cases
		{
			name:      "rpm invalid version",
			versRange: "vers:rpm/>=1.0.0",
			version:   "1.0.0@invalid",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "rpm invalid constraint version",
			versRange: "vers:rpm/>=1.0.0@invalid",
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
