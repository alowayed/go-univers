package npm

import (
	"testing"
)

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
			wantOriginal: "1.2.3",
		},
		{
			name:         "caret range",
			input:        "^1.2.3",
			wantOriginal: "^1.2.3",
		},
		{
			name:         "tilde range",
			input:        "~1.2.3",
			wantOriginal: "~1.2.3",
		},
		{
			name:         "x-range",
			input:        "1.x",
			wantOriginal: "1.x",
		},
		{
			name:         "hyphen range",
			input:        "1.2.3 - 2.3.4",
			wantOriginal: "1.2.3 - 2.3.4",
		},
		{
			name:         "multiple constraints",
			input:        ">=1.0.0 <2.0.0",
			wantOriginal: ">=1.0.0 <2.0.0",
		},
		{
			name:         "OR range",
			input:        "1.x || 2.x",
			wantOriginal: "1.x || 2.x",
		},
		{
			name:    "empty range",
			input:   "",
			wantErr: true,
		},
		{
			name:         "range with leading whitespace",
			input:        " ^1.2.3",
			wantOriginal: "^1.2.3",
		},
		{
			name:         "range with trailing whitespace",
			input:        "~1.2.3 ",
			wantOriginal: "~1.2.3",
		},
		{
			name:         "OR with extra spaces",
			input:        "1.x  ||  2.x",
			wantOriginal: "1.x  ||  2.x",
		},
		{
			name:    "invalid double tilde",
			input:   "~~1.2.3",
			wantErr: true,
		},
		{
			name:    "invalid double caret",
			input:   "^^1.2.3",
			wantErr: true,
		},
		{
			name:    "malformed hyphen range - missing end",
			input:   "1.2.3 -",
			wantErr: true,
		},
		{
			name:    "malformed hyphen range - missing start",
			input:   "- 1.2.3",
			wantErr: true,
		},
		{
			name:    "malformed hyphen range - double hyphen",
			input:   "1.2.3 - - 2.0.0",
			wantErr: true,
		},
		{
			name:         "parentheses around range",
			input:        "(>=1.0.0 <2.0.0)",
			wantOriginal: "(>=1.0.0 <2.0.0)",
		},
		{
			name:    "invalid characters",
			input:   "1.2.3@invalid",
			wantErr: true,
		},
		{
			name:    "only whitespace",
			input:   "   ",
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
		{
			name:     "caret range with prerelease base",
			rangeStr: "^1.2.3-alpha",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "caret range with prerelease - includes prerelease",
			rangeStr: "^1.2.3-alpha",
			version:  "1.3.0-alpha",
			want:     true,
		},
		{
			name:     "caret range with zero major",
			rangeStr: "^0.0.1",
			version:  "0.0.1",
			want:     true,
		},
		{
			name:     "caret range with zero major - no match",
			rangeStr: "^0.0.1",
			version:  "0.0.2",
		},
		{
			name:     "tilde range with zero minor",
			rangeStr: "~0.1.0",
			version:  "0.1.5",
			want:     true,
		},
		{
			name:     "caret range with zero minor",
			rangeStr: "^0.1.0",
			version:  "0.1.5",
			want:     true,
		},
		{
			name:     "caret range with zero minor - no match",
			rangeStr: "^0.1.0",
			version:  "0.2.0",
		},
		{
			name:     "upper boundary exactness - caret excludes next major",
			rangeStr: "^1.2.3",
			version:  "2.0.0",
		},
		{
			name:     "caret range excludes prerelease from next major",
			rangeStr: "^1.2.3",
			version:  "2.0.0-alpha",
		},
		{
			name:     "complex OR logic with prerelease",
			rangeStr: "1.x || >=2.0.0-alpha <3.0.0",
			version:  "2.5.0-beta",
			want:     true,
		},
		{
			name:     "range with build metadata (ignored)",
			rangeStr: "1.2.3+build",
			version:  "1.2.3+different",
			want:     true,
		},
		{
			name:     "prerelease version against non-prerelease range",
			rangeStr: ">=1.0.0",
			version:  "1.0.0-alpha",
		},
		{
			name:     "zero patch version in range",
			rangeStr: "^1.2.0",
			version:  "1.2.5",
			want:     true,
		},
		{
			name:     "hyphen range with prerelease",
			rangeStr: "1.0.0-alpha - 2.0.0",
			version:  "1.5.0-beta",
			want:     true,
		},
		{
			name:     "X-range includes prereleases - major",
			rangeStr: "1.x",
			version:  "1.0.0-alpha",
			want:     true,
		},
		{
			name:     "X-range includes prereleases - minor",
			rangeStr: "1.2.x",
			version:  "1.2.0-beta",
			want:     true,
		},
		{
			name:     "X-range includes prereleases - complex",
			rangeStr: "2.x",
			version:  "2.5.0-rc.1",
			want:     true,
		},
		{
			name:     "X-range excludes different major with prerelease",
			rangeStr: "1.x",
			version:  "2.0.0-alpha",
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