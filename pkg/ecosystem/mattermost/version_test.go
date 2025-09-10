package mattermost

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
			name:  "Mattermost semantic v8.1.5",
			input: "v8.1.5",
			want: &Version{
				original: "v8.1.5",
				prefix:   "v",
				major:    8,
				minor:    1,
				patch:    5,
			},
		},
		{
			name:  "Mattermost semantic v10.12.0",
			input: "v10.12.0",
			want: &Version{
				original: "v10.12.0",
				prefix:   "v",
				major:    10,
				minor:    12,
				patch:    0,
			},
		},
		// Semantic versions without prefix
		{
			name:  "Mattermost semantic 8.1.5",
			input: "8.1.5",
			want: &Version{
				original: "8.1.5",
				prefix:   "",
				major:    8,
				minor:    1,
				patch:    5,
			},
		},
		{
			name:  "Mattermost semantic 10.12.0",
			input: "10.12.0",
			want: &Version{
				original: "10.12.0",
				prefix:   "",
				major:    10,
				minor:    12,
				patch:    0,
			},
		},
		// ESR versions
		{
			name:  "Mattermost ESR v8.1.5-esr",
			input: "v8.1.5-esr",
			want: &Version{
				original:  "v8.1.5-esr",
				prefix:    "v",
				major:     8,
				minor:     1,
				patch:     5,
				qualifier: "esr",
				number:    0,
			},
		},
		{
			name:  "Mattermost ESR 7.8.11-esr",
			input: "7.8.11-esr",
			want: &Version{
				original:  "7.8.11-esr",
				prefix:    "",
				major:     7,
				minor:     8,
				patch:     11,
				qualifier: "esr",
				number:    0,
			},
		},
		// Release candidates
		{
			name:  "Mattermost RC v8.1.0-rc1",
			input: "v8.1.0-rc1",
			want: &Version{
				original:  "v8.1.0-rc1",
				prefix:    "v",
				major:     8,
				minor:     1,
				patch:     0,
				qualifier: "rc",
				number:    1,
			},
		},
		{
			name:  "Mattermost RC v10.12.0-rc2",
			input: "v10.12.0-rc2",
			want: &Version{
				original:  "v10.12.0-rc2",
				prefix:    "v",
				major:     10,
				minor:     12,
				patch:     0,
				qualifier: "rc",
				number:    2,
			},
		},
		{
			name:  "Mattermost RC 8.1.0-rc",
			input: "8.1.0-rc",
			want: &Version{
				original:  "8.1.0-rc",
				prefix:    "",
				major:     8,
				minor:     1,
				patch:     0,
				qualifier: "rc",
				number:    0,
			},
		},
		// Real Mattermost examples
		{
			name:  "Mattermost current v10.11.2",
			input: "v10.11.2",
			want: &Version{
				original: "v10.11.2",
				prefix:   "v",
				major:    10,
				minor:    11,
				patch:    2,
			},
		},
		{
			name:  "Mattermost v9.5.10",
			input: "v9.5.10",
			want: &Version{
				original: "v9.5.10",
				prefix:   "v",
				major:    9,
				minor:    5,
				patch:    10,
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
			input:   "v8.1",
			wantErr: true,
		},
		{
			name:    "invalid format - no minor",
			input:   "v8",
			wantErr: true,
		},
		{
			name:    "invalid format - letters in version numbers",
			input:   "v8.a.3",
			wantErr: true,
		},
		{
			name:    "invalid format - unsupported qualifier",
			input:   "v8.1.0-alpha",
			wantErr: true,
		},
		{
			name:    "invalid format - multiple qualifiers",
			input:   "v8.1.0-rc1-esr",
			wantErr: true,
		},
		{
			name:    "invalid format - leading zeros",
			input:   "v08.01.05",
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
			v1:   "v8.1.5",
			v2:   "v8.1.5",
			want: 0,
		},
		{
			name: "major version difference",
			v1:   "v8.1.0",
			v2:   "v9.0.0",
			want: -1,
		},
		{
			name: "minor version difference",
			v1:   "v8.1.0",
			v2:   "v8.2.0",
			want: -1,
		},
		{
			name: "patch version difference",
			v1:   "v8.1.1",
			v2:   "v8.1.2",
			want: -1,
		},
		// Prefix handling
		{
			name: "v prefix vs no prefix same version",
			v1:   "v8.1.5",
			v2:   "8.1.5",
			want: 0,
		},
		// Qualifier comparisons
		{
			name: "release vs rc",
			v1:   "v8.1.0",
			v2:   "v8.1.0-rc1",
			want: 1,
		},
		{
			name: "release vs esr",
			v1:   "v8.1.5",
			v2:   "v8.1.5-esr",
			want: 1,
		},
		{
			name: "rc vs esr",
			v1:   "v8.1.0-rc1",
			v2:   "v8.1.5-esr",
			want: -1, // rc < esr in precedence
		},
		{
			name: "numbered rc versions",
			v1:   "v8.1.0-rc1",
			v2:   "v8.1.0-rc2",
			want: -1,
		},
		{
			name: "rc with and without number",
			v1:   "v8.1.0-rc",
			v2:   "v8.1.0-rc1",
			want: -1,
		},
		// Real Mattermost scenarios
		{
			name: "Mattermost 10.x series",
			v1:   "v10.11.0",
			v2:   "v10.11.2",
			want: -1,
		},
		{
			name: "Mattermost ESR comparison",
			v1:   "v8.1.5-esr",
			v2:   "v7.8.11-esr",
			want: 1,
		},
		{
			name: "Mattermost RC to release",
			v1:   "v10.12.0-rc1",
			v2:   "v10.12.0",
			want: -1,
		},
		{
			name: "Cross-major versions",
			v1:   "v9.5.10",
			v2:   "v10.0.0",
			want: -1,
		},
		{
			name: "ESR vs newer release",
			v1:   "v8.1.5-esr",
			v2:   "v9.0.0",
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
			input: "v8.1.5",
		},
		{
			name:  "semantic version without prefix",
			input: "8.1.5",
		},
		{
			name:  "ESR version",
			input: "v8.1.5-esr",
		},
		{
			name:  "release candidate",
			input: "v8.1.0-rc1",
		},
		{
			name:  "double digit versions",
			input: "v10.12.0",
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
