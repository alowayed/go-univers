package gomod

import (
	"testing"
	"time"
)

func TestEcosystem_NewVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Version
		wantErr bool
	}{
		{
			name:  "basic semantic version with v prefix",
			input: "v1.2.3",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "v1.2.3",
			},
		},
		{
			name:  "basic semantic version without v prefix",
			input: "1.2.3",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "1.2.3",
			},
		},
		{
			name:  "version with prerelease",
			input: "v1.2.3-beta.1",
			want: &Version{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "beta.1",
				original:   "v1.2.3-beta.1",
			},
		},
		{
			name:  "version with build metadata",
			input: "v1.2.3+build.1",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				build:    "build.1",
				original: "v1.2.3+build.1",
			},
		},
		{
			name:  "version with prerelease and build",
			input: "v1.2.3-alpha+build.1",
			want: &Version{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "alpha",
				build:      "build.1",
				original:   "v1.2.3-alpha+build.1",
			},
		},
		{
			name:  "pseudo-version pattern 1",
			input: "v1.0.0-20170915032832-14c0d48ead0c",
			want: &Version{
				major:    1,
				minor:    0,
				patch:    0,
				original: "v1.0.0-20170915032832-14c0d48ead0c",
				pseudo: &pseudoVersion{
					baseVersion: "v1.0.0",
					timestamp:   mustParseTime(t, "20170915032832"),
					revision:    "14c0d48ead0c",
				},
			},
		},
		{
			name:  "pseudo-version pattern 2",
			input: "v1.2.3-beta.0.20170915032832-14c0d48ead0c",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "v1.2.3-beta.0.20170915032832-14c0d48ead0c",
				pseudo: &pseudoVersion{
					baseVersion: "v1.2.3-beta",
					timestamp:   mustParseTime(t, "20170915032832"),
					revision:    "14c0d48ead0c",
				},
			},
		},
		{
			name:  "pseudo-version pattern 3",
			input: "v1.2.4-0.20170915032832-14c0d48ead0c",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    4,
				original: "v1.2.4-0.20170915032832-14c0d48ead0c",
				pseudo: &pseudoVersion{
					baseVersion: "v1.2.3",
					timestamp:   mustParseTime(t, "20170915032832"),
					revision:    "14c0d48ead0c",
				},
			},
		},
		{
			name:    "invalid version",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty version",
			input:   "",
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
					return
				}

				// Check pseudo-version specifics
				if tt.want.pseudo != nil {
					if got.pseudo == nil {
						t.Errorf("NewVersion() pseudo = nil, want %+v", tt.want.pseudo)
					} else {
						if got.pseudo.baseVersion != tt.want.pseudo.baseVersion {
							t.Errorf("NewVersion() pseudo.baseVersion = %q, want %q", got.pseudo.baseVersion, tt.want.pseudo.baseVersion)
						}
						if got.pseudo.revision != tt.want.pseudo.revision {
							t.Errorf("NewVersion() pseudo.revision = %q, want %q", got.pseudo.revision, tt.want.pseudo.revision)
						}
						if !got.pseudo.timestamp.Equal(tt.want.pseudo.timestamp) {
							t.Errorf("NewVersion() pseudo.timestamp = %v, want %v", got.pseudo.timestamp, tt.want.pseudo.timestamp)
						}
					}
				} else if got.pseudo != nil {
					t.Errorf("NewVersion() pseudo = %+v, want nil", got.pseudo)
				}
			}
		})
	}
}

func TestVersionCompare(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		{
			name: "equal versions",
			v1:   "v1.2.3",
			v2:   "v1.2.3",
			want: 0,
		},
		{
			name: "major version difference",
			v1:   "v1.2.3",
			v2:   "v2.2.3",
			want: -1,
		},
		{
			name: "minor version difference",
			v1:   "v1.3.3",
			v2:   "v1.2.3",
			want: 1,
		},
		{
			name: "patch version difference",
			v1:   "v1.2.3",
			v2:   "v1.2.4",
			want: -1,
		},
		{
			name: "prerelease vs release",
			v1:   "v1.2.3-alpha",
			v2:   "v1.2.3",
			want: -1,
		},
		{
			name: "prerelease comparison",
			v1:   "v1.2.3-alpha",
			v2:   "v1.2.3-beta",
			want: -1,
		},
		{
			name: "pseudo-version vs release",
			v1:   "v1.2.3-0.20170915032832-14c0d48ead0c",
			v2:   "v1.2.3",
			want: -1,
		},
		{
			name: "pseudo-version comparison",
			v1:   "v1.2.3-0.20170915032832-14c0d48ead0c",
			v2:   "v1.2.3-0.20170916032832-14c0d48ead0c",
			want: -1,
		},
		{
			name: "build metadata ignored",
			v1:   "v1.2.3+build1",
			v2:   "v1.2.3+build2",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ver1 := mustNewVersion(t, tt.v1)
			ver2 := mustNewVersion(t, tt.v2)

			got := ver1.Compare(ver2)
			if got != tt.want {
				t.Errorf("Version.Compare() = %d, want %d", got, tt.want)
			}

			// Test symmetry
			if tt.want != 0 {
				reverseGot := ver2.Compare(ver1)
				if reverseGot != -tt.want {
					t.Errorf("Version.Compare() not symmetric: %q vs %q = %d, but %q vs %q = %d",
						tt.v1, tt.v2, got, tt.v2, tt.v1, reverseGot)
				}
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"v prefix", "v1.2.3", "v1.2.3"},
		{"no v prefix", "1.2.3", "1.2.3"},
		{"with prerelease", "v1.2.3-beta", "v1.2.3-beta"},
		{"pseudo-version", "v1.0.0-20170915032832-14c0d48ead0c", "v1.0.0-20170915032832-14c0d48ead0c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustNewVersion(t, tt.input)
			if got.String() != tt.want {
				t.Errorf("Version.String() = %q, want %q", got.String(), tt.want)
			}
		})
	}
}

// mustParseTime is a test helper to parse timestamps
func mustParseTime(t *testing.T, timestamp string) time.Time {
	t.Helper()
	ts, err := time.Parse("20060102150405", timestamp)
	if err != nil {
		t.Fatalf("failed to parse time %q: %v", timestamp, err)
	}
	return ts
}

func mustNewVersion(t *testing.T, input string) *Version {
	t.Helper()
	e := &Ecosystem{}
	v, err := e.NewVersion(input)
	if err != nil {
		t.Fatalf("NewVersion(%q) error = %v", input, err)
	}
	return v
}
