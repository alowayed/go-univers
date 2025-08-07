package npm

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Version
		wantErr bool
	}{
		{
			name:  "basic semantic version",
			input: "1.2.3",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "1.2.3",
			},
		},
		{
			name:  "version with v prefix",
			input: "v1.2.3",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "v1.2.3",
			},
		},
		{
			name:  "version with equals prefix",
			input: "=1.2.3",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "=1.2.3",
			},
		},
		{
			name:  "version with prerelease",
			input: "1.2.3-alpha.1",
			want: &Version{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "alpha.1",
				original:   "1.2.3-alpha.1",
			},
		},
		{
			name:  "version with build metadata",
			input: "1.2.3+build.1",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				build:    "build.1",
				original: "1.2.3+build.1",
			},
		},
		{
			name:  "version with prerelease and build",
			input: "1.2.3-alpha.1+build.1",
			want: &Version{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "alpha.1",
				build:      "build.1",
				original:   "1.2.3-alpha.1+build.1",
			},
		},
		{
			name:    "invalid version",
			input:   "1.2",
			wantErr: true,
		},
		{
			name:    "empty version",
			input:   "",
			wantErr: true,
		},
		{
			name:    "non-numeric version",
			input:   "a.b.c",
			wantErr: true,
		},
		{
			name:  "multiple v prefixes",
			input: "vv1.2.3",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "vv1.2.3",
			},
		},
		{
			name:    "triple v prefixes",
			input:   "vvv1.2.3",
			wantErr: true,
		},
		{
			name:  "zero-padded version",
			input: "01.02.03",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "01.02.03",
			},
		},
		{
			name:  "numeric prerelease with zero padding",
			input: "1.2.3-01",
			want: &Version{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "01",
				original:   "1.2.3-01",
			},
		},
		{
			name:  "complex prerelease",
			input: "1.2.3-alpha.1.beta",
			want: &Version{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "alpha.1.beta",
				original:   "1.2.3-alpha.1.beta",
			},
		},
		{
			name:  "build metadata with hyphens",
			input: "1.2.3+build-01.02",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				build:    "build-01.02",
				original: "1.2.3+build-01.02",
			},
		},
		{
			name:  "version with leading whitespace",
			input: " 1.2.3",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "1.2.3",
			},
		},
		{
			name:  "version with trailing whitespace",
			input: "1.2.3 ",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "1.2.3",
			},
		},
		{
			name:  "version with tabs",
			input: "1.2.3\t",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "1.2.3",
			},
		},
		{
			name:    "two-part version",
			input:   "1.2",
			wantErr: true,
		},
		{
			name:    "four-part version",
			input:   "1.2.3.4",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ecosystem{}
			got, err := e.NewVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.major != tt.want.major || got.minor != tt.want.minor || got.patch != tt.want.patch ||
					got.prerelease != tt.want.prerelease || got.build != tt.want.build || got.original != tt.want.original {
					t.Errorf("NewVersion() = %+v, want %+v", got, tt.want)
				}
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
		{
			name: "equal versions",
			v1:   "1.2.3",
			v2:   "1.2.3",
			want: 0,
		},
		{
			name: "major version difference",
			v1:   "1.2.3",
			v2:   "2.2.3",
			want: -1,
		},
		{
			name: "minor version difference",
			v1:   "1.3.3",
			v2:   "1.2.3",
			want: 1,
		},
		{
			name: "patch version difference",
			v1:   "1.2.2",
			v2:   "1.2.3",
			want: -1,
		},
		{
			name: "prerelease vs release",
			v1:   "1.2.3-alpha",
			v2:   "1.2.3",
			want: -1,
		},
		{
			name: "prerelease comparison",
			v1:   "1.2.3-alpha.1",
			v2:   "1.2.3-alpha.2",
			want: -1,
		},
		{
			name: "prerelease numeric vs string",
			v1:   "1.2.3-alpha.1",
			v2:   "1.2.3-alpha.beta",
			want: -1,
		},
		{
			name: "build metadata ignored in comparison",
			v1:   "1.2.3+build1",
			v2:   "1.2.3+build2",
			want: 0,
		},
		{
			name: "complex prerelease ordering alpha < beta",
			v1:   "1.2.3-alpha",
			v2:   "1.2.3-beta",
			want: -1,
		},
		{
			name: "complex prerelease ordering with numbers",
			v1:   "1.2.3-alpha.1.2",
			v2:   "1.2.3-alpha.1.10",
			want: -1,
		},
		{
			name: "zero-padded prerelease comparison",
			v1:   "1.2.3-01",
			v2:   "1.2.3-1",
			want: 0,
		},
		{
			name: "numeric vs alphabetic prerelease",
			v1:   "1.2.3-1",
			v2:   "1.2.3-a",
			want: -1,
		},
		{
			name: "prerelease with multiple components",
			v1:   "1.2.3-alpha.1",
			v2:   "1.2.3-alpha.1.0",
			want: -1,
		},
		{
			name: "prerelease rc vs release",
			v1:   "1.2.3-rc.1",
			v2:   "1.2.3",
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := mustNewVersion(t, tt.v1)
			v2 := mustNewVersion(t, tt.v2)

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

// mustNewVersion is a helper function to create a new Version.
func mustNewVersion(t *testing.T, version string) *Version {
	t.Helper()
	e := &Ecosystem{}
	v, err := e.NewVersion(version)
	if err != nil {
		t.Fatalf("Failed to create version %s: %v", version, err)
	}
	return v
}
