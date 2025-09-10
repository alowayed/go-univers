package hex

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
		// Basic semantic versions
		{
			name:  "Hex basic version 1.0.0",
			input: "1.0.0",
			want: &Version{
				original: "1.0.0",
				major:    1,
				minor:    0,
				patch:    0,
			},
		},
		{
			name:  "Hex version 2.3.1",
			input: "2.3.1",
			want: &Version{
				original: "2.3.1",
				major:    2,
				minor:    3,
				patch:    1,
			},
		},
		{
			name:  "Hex version 0.1.0",
			input: "0.1.0",
			want: &Version{
				original: "0.1.0",
				major:    0,
				minor:    1,
				patch:    0,
			},
		},
		// Pre-release versions
		{
			name:  "Hex pre-release 1.0.0-alpha",
			input: "1.0.0-alpha",
			want: &Version{
				original:   "1.0.0-alpha",
				major:      1,
				minor:      0,
				patch:      0,
				preRelease: []string{"alpha"},
			},
		},
		{
			name:  "Hex pre-release 1.0.0-alpha.1",
			input: "1.0.0-alpha.1",
			want: &Version{
				original:   "1.0.0-alpha.1",
				major:      1,
				minor:      0,
				patch:      0,
				preRelease: []string{"alpha", "1"},
			},
		},
		{
			name:  "Hex pre-release 2.0.0-rc.1",
			input: "2.0.0-rc.1",
			want: &Version{
				original:   "2.0.0-rc.1",
				major:      2,
				minor:      0,
				patch:      0,
				preRelease: []string{"rc", "1"},
			},
		},
		{
			name:  "Hex pre-release 1.0.0-beta.2.3",
			input: "1.0.0-beta.2.3",
			want: &Version{
				original:   "1.0.0-beta.2.3",
				major:      1,
				minor:      0,
				patch:      0,
				preRelease: []string{"beta", "2", "3"},
			},
		},
		// Build metadata versions
		{
			name:  "Hex with build metadata 1.0.0+build.1",
			input: "1.0.0+build.1",
			want: &Version{
				original:      "1.0.0+build.1",
				major:         1,
				minor:         0,
				patch:         0,
				buildMetadata: "build.1",
			},
		},
		{
			name:  "Hex with build metadata 1.0.0+20130417",
			input: "1.0.0+20130417",
			want: &Version{
				original:      "1.0.0+20130417",
				major:         1,
				minor:         0,
				patch:         0,
				buildMetadata: "20130417",
			},
		},
		{
			name:  "Hex pre-release with build metadata 1.0.0-alpha.1+build.1",
			input: "1.0.0-alpha.1+build.1",
			want: &Version{
				original:      "1.0.0-alpha.1+build.1",
				major:         1,
				minor:         0,
				patch:         0,
				preRelease:    []string{"alpha", "1"},
				buildMetadata: "build.1",
			},
		},
		// Real Hex package examples
		{
			name:  "Phoenix version 1.7.10",
			input: "1.7.10",
			want: &Version{
				original: "1.7.10",
				major:    1,
				minor:    7,
				patch:    10,
			},
		},
		{
			name:  "Ecto version 3.10.3",
			input: "3.10.3",
			want: &Version{
				original: "3.10.3",
				major:    3,
				minor:    10,
				patch:    3,
			},
		},
		{
			name:  "Poison pre-release 3.0.0-rc.1",
			input: "3.0.0-rc.1",
			want: &Version{
				original:   "3.0.0-rc.1",
				major:      3,
				minor:      0,
				patch:      0,
				preRelease: []string{"rc", "1"},
			},
		},
		// Edge cases
		{
			name:  "Version with long pre-release 1.0.0-x.7.z.92",
			input: "1.0.0-x.7.z.92",
			want: &Version{
				original:   "1.0.0-x.7.z.92",
				major:      1,
				minor:      0,
				patch:      0,
				preRelease: []string{"x", "7", "z", "92"},
			},
		},
		{
			name:  "Version with complex build metadata 1.0.0+exp.sha.5114f85",
			input: "1.0.0+exp.sha.5114f85",
			want: &Version{
				original:      "1.0.0+exp.sha.5114f85",
				major:         1,
				minor:         0,
				patch:         0,
				buildMetadata: "exp.sha.5114f85",
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
			input:   "1.2.a",
			wantErr: true,
		},
		{
			name:    "invalid format - no minor",
			input:   "1",
			wantErr: true,
		},
		{
			name:    "invalid format - letters in version numbers",
			input:   "1.a.3",
			wantErr: true,
		},
		{
			name:    "invalid format - invalid characters",
			input:   "1.2.3.4",
			wantErr: true,
		},
		{
			name:    "invalid pre-release - empty identifier",
			input:   "1.0.0-alpha..1",
			wantErr: true,
		},
		{
			name:    "invalid format - no major version",
			input:   ".1.2",
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
			v1:   "1.0.0",
			v2:   "1.0.0",
			want: 0,
		},
		{
			name: "major version difference",
			v1:   "1.0.0",
			v2:   "2.0.0",
			want: -1,
		},
		{
			name: "minor version difference",
			v1:   "1.2.0",
			v2:   "1.3.0",
			want: -1,
		},
		{
			name: "patch version difference",
			v1:   "1.0.1",
			v2:   "1.0.2",
			want: -1,
		},
		// Pre-release comparisons
		{
			name: "release vs pre-release",
			v1:   "1.0.0",
			v2:   "1.0.0-alpha",
			want: 1,
		},
		{
			name: "pre-release vs release",
			v1:   "1.0.0-alpha",
			v2:   "1.0.0",
			want: -1,
		},
		{
			name: "alpha vs beta",
			v1:   "1.0.0-alpha",
			v2:   "1.0.0-beta",
			want: -1,
		},
		{
			name: "alpha vs rc",
			v1:   "1.0.0-alpha",
			v2:   "1.0.0-rc",
			want: -1,
		},
		{
			name: "beta vs rc",
			v1:   "1.0.0-beta",
			v2:   "1.0.0-rc",
			want: -1,
		},
		{
			name: "numbered pre-releases",
			v1:   "1.0.0-alpha.1",
			v2:   "1.0.0-alpha.2",
			want: -1,
		},
		{
			name: "alpha vs alpha.1",
			v1:   "1.0.0-alpha",
			v2:   "1.0.0-alpha.1",
			want: -1,
		},
		{
			name: "numeric vs string pre-release",
			v1:   "1.0.0-1",
			v2:   "1.0.0-alpha",
			want: -1,
		},
		{
			name: "complex pre-release comparison",
			v1:   "1.0.0-alpha.1",
			v2:   "1.0.0-alpha.beta",
			want: -1,
		},
		// Build metadata (should be ignored)
		{
			name: "same version with different build metadata",
			v1:   "1.0.0+build.1",
			v2:   "1.0.0+build.2",
			want: 0,
		},
		{
			name: "version vs version with build metadata",
			v1:   "1.0.0",
			v2:   "1.0.0+build.1",
			want: 0,
		},
		{
			name: "pre-release with build metadata",
			v1:   "1.0.0-alpha+build.1",
			v2:   "1.0.0-alpha+build.2",
			want: 0,
		},
		// Real-world Hex comparisons
		{
			name: "Phoenix versions",
			v1:   "1.7.9",
			v2:   "1.7.10",
			want: -1,
		},
		{
			name: "Ecto versions",
			v1:   "3.10.2",
			v2:   "3.10.3",
			want: -1,
		},
		{
			name: "Major upgrade",
			v1:   "2.1.0",
			v2:   "3.0.0",
			want: -1,
		},
		{
			name: "RC to release",
			v1:   "3.0.0-rc.1",
			v2:   "3.0.0",
			want: -1,
		},
		// Edge cases
		{
			name: "zero versions",
			v1:   "0.1.0",
			v2:   "0.2.0",
			want: -1,
		},
		{
			name: "complex pre-release with numbers",
			v1:   "1.0.0-x.7.z.92",
			v2:   "1.0.0-x.8.z.92",
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
			name:  "basic version",
			input: "1.0.0",
		},
		{
			name:  "version with pre-release",
			input: "2.0.0-alpha.1",
		},
		{
			name:  "version with build metadata",
			input: "1.0.0+build.1",
		},
		{
			name:  "version with both",
			input: "1.0.0-rc.1+build.123",
		},
		{
			name:  "real hex package version",
			input: "3.10.3",
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
