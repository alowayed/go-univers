package npm

import (
	"testing"
)

func TestNewVersion(t *testing.T) {
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
			wantErr: false,
		},
		{
			name:  "version with v prefix",
			input: "v1.2.3",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "1.2.3",
			},
			wantErr: false,
		},
		{
			name:  "version with equals prefix",
			input: "=1.2.3",
			want: &Version{
				major:    1,
				minor:    2,
				patch:    3,
				original: "1.2.3",
			},
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
		},
		{
			name:    "invalid version",
			input:   "1.2",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty version",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "non-numeric version",
			input:   "a.b.c",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVersion(tt.input)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("Failed to parse v1: %v", err)
			}
			v2, err := NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("Failed to parse v2: %v", err)
			}

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_Normalize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "basic version",
			input: "1.2.3",
			want:  "1.2.3",
		},
		{
			name:  "version with prerelease",
			input: "1.2.3-alpha.1",
			want:  "1.2.3-alpha.1",
		},
		{
			name:  "version with build",
			input: "1.2.3+build.1",
			want:  "1.2.3+build.1",
		},
		{
			name:  "version with prerelease and build",
			input: "1.2.3-alpha.1+build.1",
			want:  "1.2.3-alpha.1+build.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewVersion(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse version: %v", err)
			}

			got := v.Normalize()
			if got != tt.want {
				t.Errorf("Version.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewVersionRange(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantErr      bool
		wantOriginal string
	}{
		{
			name:         "exact version",
			input:        "1.2.3",
			wantErr:      false,
			wantOriginal: "1.2.3",
		},
		{
			name:         "caret range",
			input:        "^1.2.3",
			wantErr:      false,
			wantOriginal: "^1.2.3",
		},
		{
			name:         "tilde range",
			input:        "~1.2.3",
			wantErr:      false,
			wantOriginal: "~1.2.3",
		},
		{
			name:         "x-range",
			input:        "1.x",
			wantErr:      false,
			wantOriginal: "1.x",
		},
		{
			name:         "hyphen range",
			input:        "1.2.3 - 2.3.4",
			wantErr:      false,
			wantOriginal: "1.2.3 - 2.3.4",
		},
		{
			name:         "multiple constraints",
			input:        ">=1.0.0 <2.0.0",
			wantErr:      false,
			wantOriginal: ">=1.0.0 <2.0.0",
		},
		{
			name:         "OR range",
			input:        "1.x || 2.x",
			wantErr:      false,
			wantOriginal: "1.x || 2.x",
		},
		{
			name:    "empty range",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersionRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.wantOriginal {
				t.Errorf("NewVersionRange().String() = %v, want %v", got.String(), tt.wantOriginal)
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	tests := []struct {
		name        string
		rangeStr    string
		version     string
		want        bool
		wantRangeErr bool
		wantVersionErr bool
	}{
		{
			name:     "exact match",
			rangeStr: "1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "exact no match",
			rangeStr: "1.2.3",
			version:  "1.2.4",
			want:     false,
		},
		{
			name:     "caret range match",
			rangeStr: "^1.2.3",
			version:  "1.2.5",
			want:     true,
		},
		{
			name:     "caret range no match major",
			rangeStr: "^1.2.3",
			version:  "2.0.0",
			want:     false,
		},
		{
			name:     "tilde range match",
			rangeStr: "~1.2.3",
			version:  "1.2.5",
			want:     true,
		},
		{
			name:     "tilde range no match minor",
			rangeStr: "~1.2.3",
			version:  "1.3.0",
			want:     false,
		},
		{
			name:     "x-range match",
			rangeStr: "1.x",
			version:  "1.5.0",
			want:     true,
		},
		{
			name:     "x-range no match",
			rangeStr: "1.x",
			version:  "2.0.0",
			want:     false,
		},
		{
			name:     "hyphen range match",
			rangeStr: "1.2.3 - 2.3.4",
			version:  "1.5.0",
			want:     true,
		},
		{
			name:     "hyphen range no match",
			rangeStr: "1.2.3 - 2.3.4",
			version:  "3.0.0",
			want:     false,
		},
		{
			name:     "multiple constraints match",
			rangeStr: ">=1.0.0 <2.0.0",
			version:  "1.5.0",
			want:     true,
		},
		{
			name:     "multiple constraints no match",
			rangeStr: ">=1.0.0 <2.0.0",
			version:  "2.0.0",
			want:     false,
		},
		{
			name:     "wildcard match",
			rangeStr: "*",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "prerelease constraint",
			rangeStr: ">=1.0.0-alpha",
			version:  "1.0.0-beta",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr, err := NewVersionRange(tt.rangeStr)
			if (err != nil) != tt.wantRangeErr {
				t.Errorf("NewVersionRange() error = %v, wantRangeErr %v", err, tt.wantRangeErr)
				return
			}
			if tt.wantRangeErr {
				return
			}

			v, err := NewVersion(tt.version)
			if (err != nil) != tt.wantVersionErr {
				t.Errorf("NewVersion() error = %v, wantVersionErr %v", err, tt.wantVersionErr)
				return
			}
			if tt.wantVersionErr {
				return
			}

			got := vr.Contains(v)
			if got != tt.want {
				t.Errorf("VersionRange.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstraint_Matches(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		constVer string
		testVer  string
		want     bool
	}{
		{
			name:     "equals match",
			operator: "=",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     true,
		},
		{
			name:     "equals no match",
			operator: "=",
			constVer: "1.2.3",
			testVer:  "1.2.4",
			want:     false,
		},
		{
			name:     "not equals match",
			operator: "!=",
			constVer: "1.2.3",
			testVer:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than match",
			operator: ">",
			constVer: "1.2.3",
			testVer:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than no match",
			operator: ">",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     false,
		},
		{
			name:     "greater than or equal match",
			operator: ">=",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     true,
		},
		{
			name:     "less than match",
			operator: "<",
			constVer: "1.2.3",
			testVer:  "1.2.2",
			want:     true,
		},
		{
			name:     "less than or equal match",
			operator: "<=",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     true,
		},
		{
			name:     "wildcard match",
			operator: "*",
			constVer: "*",
			testVer:  "1.2.3",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Constraint{
				operator: tt.operator,
				version:  tt.constVer,
			}

			v, err := NewVersion(tt.testVer)
			if err != nil {
				t.Fatalf("Failed to parse test version: %v", err)
			}

			got := c.Matches(v)
			if got != tt.want {
				t.Errorf("Constraint.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}