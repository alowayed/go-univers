package github

import (
	"reflect"
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Version
		wantErr bool
	}{
		// Semantic versions with v prefix
		{
			name:  "GitHub semantic v1.0.0",
			input: "v1.0.0",
			want: &Version{
				original:    "v1.0.0",
				prefix:      "v",
				major:       1,
				minor:       0,
				patch:       0,
				isDateBased: false,
			},
		},
		{
			name:  "GitHub semantic v2.3.1",
			input: "v2.3.1",
			want: &Version{
				original:    "v2.3.1",
				prefix:      "v",
				major:       2,
				minor:       3,
				patch:       1,
				isDateBased: false,
			},
		},
		// Semantic versions without prefix
		{
			name:  "GitHub semantic 1.0.0",
			input: "1.0.0",
			want: &Version{
				original:    "1.0.0",
				prefix:      "",
				major:       1,
				minor:       0,
				patch:       0,
				isDateBased: false,
			},
		},
		{
			name:  "GitHub semantic 2.3.1",
			input: "2.3.1",
			want: &Version{
				original:    "2.3.1",
				prefix:      "",
				major:       2,
				minor:       3,
				patch:       1,
				isDateBased: false,
			},
		},
		// Versions with qualifiers
		{
			name:  "GitHub v1.0.0-alpha",
			input: "v1.0.0-alpha",
			want: &Version{
				original:    "v1.0.0-alpha",
				prefix:      "v",
				major:       1,
				minor:       0,
				patch:       0,
				qualifier:   "alpha",
				number:      0,
				isDateBased: false,
			},
		},
		{
			name:  "GitHub v2.0.0-beta.1",
			input: "v2.0.0-beta.1",
			want: &Version{
				original:    "v2.0.0-beta.1",
				prefix:      "v",
				major:       2,
				minor:       0,
				patch:       0,
				qualifier:   "beta",
				number:      1,
				isDateBased: false,
			},
		},
		{
			name:  "GitHub 1.5.0-rc.2",
			input: "1.5.0-rc.2",
			want: &Version{
				original:    "1.5.0-rc.2",
				prefix:      "",
				major:       1,
				minor:       5,
				patch:       0,
				qualifier:   "rc",
				number:      2,
				isDateBased: false,
			},
		},
		// Release prefixed versions
		{
			name:  "GitHub release-1.0.0",
			input: "release-1.0.0",
			want: &Version{
				original:    "release-1.0.0",
				prefix:      "release-",
				major:       1,
				minor:       0,
				patch:       0,
				isDateBased: false,
			},
		},
		{
			name:  "GitHub rel-2.3.1",
			input: "rel-2.3.1",
			want: &Version{
				original:    "rel-2.3.1",
				prefix:      "rel-",
				major:       2,
				minor:       3,
				patch:       1,
				isDateBased: false,
			},
		},
		// Date-based versions
		{
			name:  "GitHub date-based 2024.01.15",
			input: "2024.01.15",
			want: &Version{
				original:    "2024.01.15",
				prefix:      "",
				major:       2024,
				minor:       1,
				patch:       15,
				isDateBased: true,
			},
		},
		{
			name:  "GitHub date-based v2024.12.31",
			input: "v2024.12.31",
			want: &Version{
				original:    "v2024.12.31",
				prefix:      "v",
				major:       2024,
				minor:       12,
				patch:       31,
				isDateBased: true,
			},
		},
		{
			name:  "GitHub date-based 2023.06.01",
			input: "2023.06.01",
			want: &Version{
				original:    "2023.06.01",
				prefix:      "",
				major:       2023,
				minor:       6,
				patch:       1,
				isDateBased: true,
			},
		},
		// Error cases
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "invalid format - no patch",
			input:   "v1.2",
			wantErr: true,
		},
		{
			name:    "invalid format - no minor",
			input:   "v1",
			wantErr: true,
		},
		{
			name:    "invalid format - letters in version numbers",
			input:   "v1.a.3",
			wantErr: true,
		},
		{
			name:    "invalid date - month too high",
			input:   "2024.13.01",
			wantErr: true,
		},
		{
			name:    "invalid date - day too high",
			input:   "2024.01.32",
			wantErr: true,
		},
		{
			name:    "invalid date - month zero",
			input:   "2024.00.15",
			wantErr: true,
		},
		{
			name:    "invalid format - multiple hyphens",
			input:   "v1.0.0-beta-1",
			wantErr: true,
		},
	}

	ecosystem := &Ecosystem{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ecosystem.NewVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecosystem.NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ecosystem.NewVersion() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		// Basic semantic comparisons
		{
			name: "same version",
			v1:   "v1.0.0",
			v2:   "v1.0.0",
			want: 0,
		},
		{
			name: "major version difference",
			v1:   "v1.0.0",
			v2:   "v2.0.0",
			want: -1,
		},
		{
			name: "minor version difference",
			v1:   "v1.2.0",
			v2:   "v1.3.0",
			want: -1,
		},
		{
			name: "patch version difference",
			v1:   "v1.0.1",
			v2:   "v1.0.2",
			want: -1,
		},
		// Prefix handling
		{
			name: "v prefix vs no prefix same version",
			v1:   "v1.0.0",
			v2:   "1.0.0",
			want: 0,
		},
		{
			name: "release prefix comparison",
			v1:   "release-1.0.0",
			v2:   "release-1.0.1",
			want: -1,
		},
		// Qualifier comparisons
		{
			name: "release vs alpha",
			v1:   "v1.0.0",
			v2:   "v1.0.0-alpha",
			want: 1,
		},
		{
			name: "release vs beta",
			v1:   "v1.0.0",
			v2:   "v1.0.0-beta",
			want: 1,
		},
		{
			name: "alpha vs beta",
			v1:   "v1.0.0-alpha",
			v2:   "v1.0.0-beta",
			want: -1,
		},
		{
			name: "beta vs rc",
			v1:   "v1.0.0-beta",
			v2:   "v1.0.0-rc",
			want: -1,
		},
		{
			name: "rc vs snapshot",
			v1:   "v1.0.0-rc",
			v2:   "v1.0.0-snapshot",
			want: -1,
		},
		{
			name: "numbered qualifiers",
			v1:   "v1.0.0-beta.1",
			v2:   "v1.0.0-beta.2",
			want: -1,
		},
		// Date-based comparisons
		{
			name: "date-based same date",
			v1:   "2024.01.15",
			v2:   "2024.01.15",
			want: 0,
		},
		{
			name: "date-based different year",
			v1:   "2023.01.15",
			v2:   "2024.01.15",
			want: -1,
		},
		{
			name: "date-based different month",
			v1:   "2024.01.15",
			v2:   "2024.02.15",
			want: -1,
		},
		{
			name: "date-based different day",
			v1:   "2024.01.15",
			v2:   "2024.01.16",
			want: -1,
		},
		// Mixed date-based and semantic
		{
			name: "date-based vs semantic",
			v1:   "2024.01.15",
			v2:   "v1.0.0",
			want: -1,
		},
		{
			name: "semantic vs date-based",
			v1:   "v2.0.0",
			v2:   "2024.01.15",
			want: 1,
		},
		// Real GitHub examples
		{
			name: "Go releases",
			v1:   "v1.20.0",
			v2:   "v1.21.0",
			want: -1,
		},
		{
			name: "Kubernetes style",
			v1:   "v1.28.0",
			v2:   "v1.28.1",
			want: -1,
		},
		{
			name: "Release candidates",
			v1:   "v2.0.0-rc.1",
			v2:   "v2.0.0",
			want: -1,
		},
	}

	ecosystem := &Ecosystem{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := ecosystem.NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("Failed to parse v1 %s: %v", tt.v1, err)
			}
			v2, err := ecosystem.NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("Failed to parse v2 %s: %v", tt.v2, err)
			}

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "semantic version with v prefix",
			input: "v1.0.0",
		},
		{
			name:  "semantic version without prefix",
			input: "1.2.3",
		},
		{
			name:  "version with qualifier",
			input: "v2.0.0-beta.1",
		},
		{
			name:  "release prefixed version",
			input: "release-1.5.0",
		},
		{
			name:  "date-based version",
			input: "2024.01.15",
		},
	}

	ecosystem := &Ecosystem{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ecosystem.NewVersion(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt.input, err)
			}

			if got := v.String(); got != tt.input {
				t.Errorf("Version.String() = %v, want %v", got, tt.input)
			}
		})
	}
}
