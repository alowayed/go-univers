package alpm

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
		// Basic versions
		{
			name:  "simple version",
			input: "1.0.0",
			want: &Version{
				original: "1.0.0",
				epoch:    0,
				pkgver:   "1.0.0",
				pkgrel:   0,
			},
		},
		{
			name:  "version with release",
			input: "1.0.0-1",
			want: &Version{
				original: "1.0.0-1",
				epoch:    0,
				pkgver:   "1.0.0",
				pkgrel:   1,
			},
		},
		{
			name:  "version with epoch",
			input: "2:1.0.0",
			want: &Version{
				original: "2:1.0.0",
				epoch:    2,
				pkgver:   "1.0.0",
				pkgrel:   0,
			},
		},
		{
			name:  "version with epoch and release",
			input: "3:1.2.3-5",
			want: &Version{
				original: "3:1.2.3-5",
				epoch:    3,
				pkgver:   "1.2.3",
				pkgrel:   5,
			},
		},
		// Real Arch Linux package examples
		{
			name:  "linux kernel",
			input: "6.1.1-1",
			want: &Version{
				original: "6.1.1-1",
				epoch:    0,
				pkgver:   "6.1.1",
				pkgrel:   1,
			},
		},
		{
			name:  "glibc",
			input: "2.36-6",
			want: &Version{
				original: "2.36-6",
				epoch:    0,
				pkgver:   "2.36",
				pkgrel:   6,
			},
		},
		{
			name:  "firefox",
			input: "108.0.2-1",
			want: &Version{
				original: "108.0.2-1",
				epoch:    0,
				pkgver:   "108.0.2",
				pkgrel:   1,
			},
		},
		{
			name:  "texlive with epoch",
			input: "1:2022.62885-17",
			want: &Version{
				original: "1:2022.62885-17",
				epoch:    1,
				pkgver:   "2022.62885",
				pkgrel:   17,
			},
		},
		// Alphanumeric versions (vercmp examples)
		{
			name:  "version with alpha",
			input: "1.0a-1",
			want: &Version{
				original: "1.0a-1",
				epoch:    0,
				pkgver:   "1.0a",
				pkgrel:   1,
			},
		},
		{
			name:  "version with beta",
			input: "2.0beta-2",
			want: &Version{
				original: "2.0beta-2",
				epoch:    0,
				pkgver:   "2.0beta",
				pkgrel:   2,
			},
		},
		{
			name:  "version with rc",
			input: "3.0rc1-1",
			want: &Version{
				original: "3.0rc1-1",
				epoch:    0,
				pkgver:   "3.0rc1",
				pkgrel:   1,
			},
		},
		{
			name:  "version with pre",
			input: "1.5pre2-3",
			want: &Version{
				original: "1.5pre2-3",
				epoch:    0,
				pkgver:   "1.5pre2",
				pkgrel:   3,
			},
		},
		// Complex versions
		{
			name:  "version with underscores",
			input: "2.4_rc1-1",
			want: &Version{
				original: "2.4_rc1-1",
				epoch:    0,
				pkgver:   "2.4_rc1",
				pkgrel:   1,
			},
		},
		{
			name:  "version with plus",
			input: "1.0+git20220101-1",
			want: &Version{
				original: "1.0+git20220101-1",
				epoch:    0,
				pkgver:   "1.0+git20220101",
				pkgrel:   1,
			},
		},
		{
			name:  "version with hyphens in pkgver",
			input: "20221201-git-1",
			want: &Version{
				original: "20221201-git-1",
				epoch:    0,
				pkgver:   "20221201-git",
				pkgrel:   1,
			},
		},
		// Edge cases
		{
			name:  "zero epoch",
			input: "0:1.0-1",
			want: &Version{
				original: "0:1.0-1",
				epoch:    0,
				pkgver:   "1.0",
				pkgrel:   1,
			},
		},
		{
			name:  "zero release",
			input: "1.0-0",
			want: &Version{
				original: "1.0-0",
				epoch:    0,
				pkgver:   "1.0",
				pkgrel:   0,
			},
		},
		{
			name:  "no release number",
			input: "1.2.3",
			want: &Version{
				original: "1.2.3",
				epoch:    0,
				pkgver:   "1.2.3",
				pkgrel:   0,
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
			name:    "invalid epoch",
			input:   "abc:1.0-1",
			wantErr: true,
		},
		{
			name:    "negative epoch",
			input:   "-1:1.0-1",
			wantErr: true,
		},
		{
			name:    "invalid character in pkgver",
			input:   "1.0@-1",
			wantErr: true,
		},
		{
			name:    "invalid epoch format",
			input:   "a:1.0-1",
			wantErr: true,
		},
		{
			name:    "empty pkgver",
			input:   "1:-1",
			wantErr: true,
		},
		{
			name:    "only colon",
			input:   ":",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "::",
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
		// Basic comparisons
		{
			name: "same version",
			v1:   "1.0.0-1",
			v2:   "1.0.0-1",
			want: 0,
		},
		{
			name: "different pkgver",
			v1:   "1.0.0-1",
			v2:   "1.0.1-1",
			want: -1,
		},
		{
			name: "different pkgrel",
			v1:   "1.0.0-1",
			v2:   "1.0.0-2",
			want: -1,
		},
		{
			name: "different epoch",
			v1:   "1:1.0.0-1",
			v2:   "2:1.0.0-1",
			want: -1,
		},
		// Epoch precedence (most important)
		{
			name: "epoch overrides version",
			v1:   "2:1.0-1",
			v2:   "1:3.6-1",
			want: 1,
		},
		{
			name: "epoch overrides release",
			v1:   "1:1.0-1",
			v2:   "0:1.0-100",
			want: 1,
		},
		// vercmp alphanumeric precedence examples
		{
			name: "alpha vs beta",
			v1:   "1.0a-1",
			v2:   "1.0b-1",
			want: -1,
		},
		{
			name: "beta vs rc",
			v1:   "1.0beta-1",
			v2:   "1.0rc-1",
			want: -1,
		},
		{
			name: "pre vs rc",
			v1:   "1.0pre-1",
			v2:   "1.0rc-1",
			want: -1,
		},
		{
			name: "rc vs release",
			v1:   "1.0rc-1",
			v2:   "1.0-1",
			want: -1,
		},
		{
			name: "release vs dot extension",
			v1:   "1.0-1",
			v2:   "1.0.a-1",
			want: 1, // "1.0" (no suffix) > "1.0.a" (with suffix)
		},
		{
			name: "dot extension vs dot number",
			v1:   "1.0.a-1",
			v2:   "1.0.1-1",
			want: 1, // "1.0.a" (letter) > "1.0.1" (numeric after dot)
		},
		// Numeric precedence
		{
			name: "numeric comparison",
			v1:   "1-1",
			v2:   "1.0-1",
			want: 1, // "1" (shorter) > "1.0" (longer)
		},
		{
			name: "more numeric comparison",
			v1:   "1.1-1",
			v2:   "1.1.1-1",
			want: 1, // "1.1" (shorter) > "1.1.1" (longer)
		},
		{
			name: "version precedence",
			v1:   "1.2-1",
			v2:   "2.0-1",
			want: -1,
		},
		// Release number comparisons
		{
			name: "release precedence",
			v1:   "1.0-1",
			v2:   "1.0-2",
			want: -1,
		},
		{
			name: "no release vs with release",
			v1:   "1.0",
			v2:   "1.0-1",
			want: -1,
		},
		// Real Arch package comparisons
		{
			name: "linux kernel versions",
			v1:   "6.1.0-1",
			v2:   "6.1.1-1",
			want: -1,
		},
		{
			name: "firefox releases",
			v1:   "108.0.1-1",
			v2:   "108.0.2-1",
			want: -1,
		},
		{
			name: "glibc releases",
			v1:   "2.36-5",
			v2:   "2.36-6",
			want: -1,
		},
		// Complex version strings
		{
			name: "git versions",
			v1:   "1.0+git20220101-1",
			v2:   "1.0+git20220102-1",
			want: -1,
		},
		{
			name: "versions with underscores",
			v1:   "2.4_rc1-1",
			v2:   "2.4_rc2-1",
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
			input: "1.0.0-1",
		},
		{
			name:  "version with epoch",
			input: "2:1.5.3-2",
		},
		{
			name:  "version without release",
			input: "3.2.1",
		},
		{
			name:  "complex version",
			input: "1:2.4_rc1+git20220101-5",
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
