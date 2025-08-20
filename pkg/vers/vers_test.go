package vers

import (
	"testing"
)

func TestParseVersString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantEco     string
		wantConstr  string
		wantErr     bool
	}{
		{
			name:       "valid maven vers",
			input:      "vers:maven/>=1.0.0|<=2.0.0",
			wantEco:    "maven",
			wantConstr: ">=1.0.0|<=2.0.0",
			wantErr:    false,
		},
		{
			name:       "valid npm vers",
			input:      "vers:npm/^1.2.3",
			wantEco:    "npm",
			wantConstr: "^1.2.3",
			wantErr:    false,
		},
		{
			name:       "valid pypi vers",
			input:      "vers:pypi/~=1.2.3|!=1.2.5",
			wantEco:    "pypi",
			wantConstr: "~=1.2.3|!=1.2.5",
			wantErr:    false,
		},
		{
			name:    "missing vers prefix",
			input:   "maven/>=1.0.0",
			wantErr: true,
		},
		{
			name:    "missing slash separator",
			input:   "vers:maven>=1.0.0",
			wantErr: true,
		},
		{
			name:    "empty ecosystem",
			input:   "vers:/>=1.0.0",
			wantErr: true,
		},
		{
			name:    "empty constraints",
			input:   "vers:maven/",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEco, gotConstr, err := ParseVersString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotEco != tt.wantEco {
					t.Errorf("ParseVersString() ecosystem = %v, want %v", gotEco, tt.wantEco)
				}
				if gotConstr != tt.wantConstr {
					t.Errorf("ParseVersString() constraints = %v, want %v", gotConstr, tt.wantConstr)
				}
			}
		})
	}
}

func TestParseVersConstraints(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Constraint
		wantErr bool
	}{
		{
			name:  "single constraint",
			input: ">=1.0.0",
			want: []Constraint{
				{Operator: ">=", Version: "1.0.0"},
			},
			wantErr: false,
		},
		{
			name:  "multiple constraints",
			input: ">=1.0.0|<=2.0.0|!=1.5.0",
			want: []Constraint{
				{Operator: ">=", Version: "1.0.0"},
				{Operator: "<=", Version: "2.0.0"},
				{Operator: "!=", Version: "1.5.0"},
			},
			wantErr: false,
		},
		{
			name:  "all operators",
			input: ">=1.0.0|<=2.0.0|>0.9.0|<3.0.0|=1.2.3|!=1.5.0",
			want: []Constraint{
				{Operator: ">=", Version: "1.0.0"},
				{Operator: "<=", Version: "2.0.0"},
				{Operator: ">", Version: "0.9.0"},
				{Operator: "<", Version: "3.0.0"},
				{Operator: "=", Version: "1.2.3"},
				{Operator: "!=", Version: "1.5.0"},
			},
			wantErr: false,
		},
		{
			name:  "with spaces",
			input: ">= 1.0.0 | <= 2.0.0",
			want: []Constraint{
				{Operator: ">=", Version: "1.0.0"},
				{Operator: "<=", Version: "2.0.0"},
			},
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid operator",
			input:   "~1.0.0",
			wantErr: true,
		},
		{
			name:    "missing version",
			input:   ">=",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersConstraints(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersConstraints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("ParseVersConstraints() length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i, constraint := range got {
					if constraint != tt.want[i] {
						t.Errorf("ParseVersConstraints()[%d] = %v, want %v", i, constraint, tt.want[i])
					}
				}
			}
		})
	}
}