package apache

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
			name:  "Apache HTTP Server 2.4.41",
			input: "2.4.41",
			want: &Version{
				major:     2,
				minor:     4,
				patch:     41,
				qualifier: "",
				number:    0,
				original:  "2.4.41",
			},
		},
		{
			name:  "Apache Tomcat 9.0.45",
			input: "9.0.45",
			want: &Version{
				major:     9,
				minor:     0,
				patch:     45,
				qualifier: "",
				number:    0,
				original:  "9.0.45",
			},
		},
		{
			name:  "Apache HTTP Server 2.2.34",
			input: "2.2.34",
			want: &Version{
				major:     2,
				minor:     2,
				patch:     34,
				qualifier: "",
				number:    0,
				original:  "2.2.34",
			},
		},
		{
			name:  "Apache with RC qualifier",
			input: "2.4.41-RC1",
			want: &Version{
				major:     2,
				minor:     4,
				patch:     41,
				qualifier: "rc",
				number:    1,
				original:  "2.4.41-RC1",
			},
		},
		{
			name:  "Apache with beta qualifier",
			input: "3.0.0-beta",
			want: &Version{
				major:     3,
				minor:     0,
				patch:     0,
				qualifier: "beta",
				number:    0,
				original:  "3.0.0-beta",
			},
		},
		{
			name:  "Apache with alpha qualifier",
			input: "2.5.0-alpha",
			want: &Version{
				major:     2,
				minor:     5,
				patch:     0,
				qualifier: "alpha",
				number:    0,
				original:  "2.5.0-alpha",
			},
		},
		{
			name:  "Apache with dev qualifier",
			input: "2.5.0-dev",
			want: &Version{
				major:     2,
				minor:     5,
				patch:     0,
				qualifier: "dev",
				number:    0,
				original:  "2.5.0-dev",
			},
		},
		{
			name:  "Apache with RC2 qualifier",
			input: "2.4.0-RC2",
			want: &Version{
				major:     2,
				minor:     4,
				patch:     0,
				qualifier: "rc",
				number:    2,
				original:  "2.4.0-RC2",
			},
		},
		{
			name:  "Apache with milestone M4",
			input: "2.5.0-M4",
			want: &Version{
				major:     2,
				minor:     5,
				patch:     0,
				qualifier: "m",
				number:    4,
				original:  "2.5.0-M4",
			},
		},
		{
			name:  "Apache with full milestone qualifier",
			input: "3.0.0-milestone2",
			want: &Version{
				major:     3,
				minor:     0,
				patch:     0,
				qualifier: "m",
				number:    2,
				original:  "3.0.0-milestone2",
			},
		},
		{
			name:  "Apache with SNAPSHOT qualifier",
			input: "2.4.0-SNAPSHOT",
			want: &Version{
				major:     2,
				minor:     4,
				patch:     0,
				qualifier: "snapshot",
				number:    0,
				original:  "2.4.0-SNAPSHOT",
			},
		},
		{
			name:  "Apache Directory Project date version",
			input: "2.4.41-v20230415",
			want: &Version{
				major:     2,
				minor:     4,
				patch:     41,
				qualifier: "v",
				number:    20230415,
				original:  "2.4.41-v20230415",
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
			input:   "2.4",
			wantErr: true,
		},
		{
			name:    "invalid format - no minor",
			input:   "2",
			wantErr: true,
		},
		{
			name:    "invalid format - letters in version numbers",
			input:   "2.a.3",
			wantErr: true,
		},
		{
			name:    "invalid format - invalid qualifier",
			input:   "2.4.0-",
			wantErr: true,
		},
		{
			name:    "invalid format - multiple hyphens",
			input:   "2.4.0-RC-1",
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
			if got.major != tt.want.major {
				t.Errorf("Ecosystem.NewVersion() major = %v, want %v", got.major, tt.want.major)
			}
			if got.minor != tt.want.minor {
				t.Errorf("Ecosystem.NewVersion() minor = %v, want %v", got.minor, tt.want.minor)
			}
			if got.patch != tt.want.patch {
				t.Errorf("Ecosystem.NewVersion() patch = %v, want %v", got.patch, tt.want.patch)
			}
			if got.qualifier != tt.want.qualifier {
				t.Errorf("Ecosystem.NewVersion() qualifier = %v, want %v", got.qualifier, tt.want.qualifier)
			}
			if got.number != tt.want.number {
				t.Errorf("Ecosystem.NewVersion() number = %v, want %v", got.number, tt.want.number)
			}
			if got.original != tt.want.original {
				t.Errorf("Ecosystem.NewVersion() original = %v, want %v", got.original, tt.want.original)
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
			v1:   "2.4.41",
			v2:   "2.4.41",
			want: 0,
		},
		{
			name: "major version difference",
			v1:   "2.4.41",
			v2:   "3.0.0",
			want: -1,
		},
		{
			name: "minor version difference",
			v1:   "2.4.41",
			v2:   "2.5.0",
			want: -1,
		},
		{
			name: "patch version difference",
			v1:   "2.4.41",
			v2:   "2.4.42",
			want: -1,
		},
		// Qualifier comparisons
		{
			name: "release vs alpha",
			v1:   "2.4.41",
			v2:   "2.4.41-alpha",
			want: 1,
		},
		{
			name: "release vs beta",
			v1:   "2.4.41",
			v2:   "2.4.41-beta",
			want: 1,
		},
		{
			name: "release vs RC",
			v1:   "2.4.41",
			v2:   "2.4.41-RC1",
			want: 1,
		},
		{
			name: "alpha vs beta",
			v1:   "2.4.41-alpha",
			v2:   "2.4.41-beta",
			want: -1,
		},
		{
			name: "beta vs RC",
			v1:   "2.4.41-beta",
			v2:   "2.4.41-RC1",
			want: -1,
		},
		{
			name: "RC vs dev",
			v1:   "2.4.41-RC1",
			v2:   "2.4.41-dev",
			want: -1,
		},
		{
			name: "beta vs milestone",
			v1:   "2.4.41-beta",
			v2:   "2.4.41-M1",
			want: -1,
		},
		{
			name: "milestone vs RC",
			v1:   "2.4.41-M4",
			v2:   "2.4.41-RC1",
			want: -1,
		},
		{
			name: "RC vs snapshot",
			v1:   "2.4.41-RC1",
			v2:   "2.4.41-SNAPSHOT",
			want: -1,
		},
		{
			name: "snapshot vs dev",
			v1:   "2.4.41-SNAPSHOT",
			v2:   "2.4.41-dev",
			want: -1,
		},
		{
			name: "milestone comparison",
			v1:   "2.4.41-M1",
			v2:   "2.4.41-M4",
			want: -1,
		},
		{
			name: "date version comparison",
			v1:   "2.4.41-v20230415",
			v2:   "2.4.41-v20230416",
			want: -1,
		},
		{
			name: "date version vs release",
			v1:   "2.4.41-v20230415",
			v2:   "2.4.41",
			want: -1,
		},
		// Numbered qualifiers
		{
			name: "RC1 vs RC2",
			v1:   "2.4.41-RC1",
			v2:   "2.4.41-RC2",
			want: -1,
		},
		{
			name: "same RC",
			v1:   "2.4.41-RC1",
			v2:   "2.4.41-RC1",
			want: 0,
		},
		// Real Apache version comparisons
		{
			name: "Apache HTTP Server versions",
			v1:   "2.2.34",
			v2:   "2.4.41",
			want: -1,
		},
		{
			name: "Apache Tomcat versions",
			v1:   "8.5.75",
			v2:   "9.0.45",
			want: -1,
		},
		{
			name: "Apache HTTP Server patch versions",
			v1:   "2.4.40",
			v2:   "2.4.41",
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
			input: "2.4.41",
		},
		{
			name:  "version with RC",
			input: "2.4.41-RC1",
		},
		{
			name:  "version with beta",
			input: "3.0.0-beta",
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
