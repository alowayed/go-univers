package cli

import (
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCode int
	}{
		{
			name:     "no arguments",
			args:     []string{},
			wantCode: 1,
		},
		{
			name:     "unknown ecosystem",
			args:     []string{"unknown"},
			wantCode: 1,
		},
		{
			name:     "npm ecosystem no command",
			args:     []string{"npm"},
			wantCode: 1,
		},
		{
			name:     "npm ecosystem unknown command",
			args:     []string{"npm", "unknown"},
			wantCode: 1,
		},
		{
			name:     "npm compare success",
			args:     []string{"npm", "compare", "1.0.0", "2.0.0"},
			wantCode: 0,
		},
		{
			name:     "npm compare invalid version",
			args:     []string{"npm", "compare", "invalid", "2.0.0"},
			wantCode: 1,
		},
		{
			name:     "npm sort success",
			args:     []string{"npm", "sort", "2.0.0", "1.0.0"},
			wantCode: 0,
		},
		{
			name:     "npm contains success",
			args:     []string{"npm", "contains", "^1.0.0", "1.5.0"},
			wantCode: 0,
		},
		{
			name:     "alpine ecosystem success",
			args:     []string{"alpine", "compare", "1.0.0", "2.0.0"},
			wantCode: 0,
		},
		{
			name:     "pypi ecosystem success",
			args:     []string{"pypi", "compare", "1.0.0", "2.0.0"},
			wantCode: 0,
		},
		{
			name:     "go ecosystem success",
			args:     []string{"go", "compare", "v1.0.0", "v2.0.0"},
			wantCode: 0,
		},
		{
			name:     "maven ecosystem success",
			args:     []string{"maven", "compare", "1.0.0", "2.0.0"},
			wantCode: 0,
		},
		{
			name:     "gem ecosystem success",
			args:     []string{"gem", "compare", "1.0.0", "2.0.0"},
			wantCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode := Run(tt.args)
			if gotCode != tt.wantCode {
				t.Errorf("Run(%+v) = %v, want %v", tt.args, gotCode, tt.wantCode)
			}
		})
	}
}

func TestRun_Private(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantOut  string
		wantCode int
	}{
		{
			name:     "no arguments",
			args:     []string{},
			wantOut:  "Usage: univers <ecosystem> <command> [args]",
			wantCode: 1,
		},
		{
			name:     "unknown ecosystem",
			args:     []string{"unknown"},
			wantOut:  "Unknown ecosystem: unknown",
			wantCode: 1,
		},
		{
			name:     "npm ecosystem no command",
			args:     []string{"npm"},
			wantOut:  "No command specified for npm",
			wantCode: 1,
		},
		{
			name:     "npm ecosystem unknown command",
			args:     []string{"npm", "unknown"},
			wantOut:  "Unknown npm command: unknown",
			wantCode: 1,
		},
		{
			name:     "npm compare success less than",
			args:     []string{"npm", "compare", "1.0.0", "2.0.0"},
			wantOut:  "-1",
			wantCode: 0,
		},
		{
			name:     "npm compare success greater than",
			args:     []string{"npm", "compare", "2.0.0", "1.0.0"},
			wantOut:  "1",
			wantCode: 0,
		},
		{
			name:     "npm compare success equal",
			args:     []string{"npm", "compare", "1.0.0", "1.0.0"},
			wantOut:  "0",
			wantCode: 0,
		},
		{
			name:     "npm compare invalid first version",
			args:     []string{"npm", "compare", "invalid", "2.0.0"},
			wantOut:  "Error running command 'compare': invalid version 'invalid': invalid NPM version: invalid",
			wantCode: 1,
		},
		{
			name:     "npm compare invalid second version",
			args:     []string{"npm", "compare", "1.0.0", "invalid"},
			wantOut:  "Error running command 'compare': invalid version 'invalid': invalid NPM version: invalid",
			wantCode: 1,
		},
		{
			name:     "npm compare wrong number of args",
			args:     []string{"npm", "compare", "1.0.0"},
			wantOut:  "Error running command 'compare': compare requires exactly 2 version arguments",
			wantCode: 1,
		},
		{
			name:     "npm sort success",
			args:     []string{"npm", "sort", "2.0.0", "1.0.0", "1.5.0"},
			wantOut:  "\"1.0.0\" \"1.5.0\" \"2.0.0\"",
			wantCode: 0,
		},
		{
			name:     "npm sort single version",
			args:     []string{"npm", "sort", "1.0.0"},
			wantOut:  "\"1.0.0\"",
			wantCode: 0,
		},
		{
			name:     "npm sort no args",
			args:     []string{"npm", "sort"},
			wantOut:  "Error running command 'sort': sort requires at least 1 version argument",
			wantCode: 1,
		},
		{
			name:     "npm sort invalid version",
			args:     []string{"npm", "sort", "invalid"},
			wantOut:  "Error running command 'sort': invalid version 'invalid': invalid NPM version: invalid",
			wantCode: 1,
		},
		{
			name:     "npm contains success true",
			args:     []string{"npm", "contains", "^1.0.0", "1.5.0"},
			wantOut:  "true",
			wantCode: 0,
		},
		{
			name:     "npm contains success false",
			args:     []string{"npm", "contains", "^1.0.0", "2.0.0"},
			wantOut:  "false",
			wantCode: 0,
		},
		{
			name:     "npm contains invalid range",
			args:     []string{"npm", "contains", "invalid", "1.0.0"},
			wantOut:  "false",
			wantCode: 0,
		},
		{
			name:     "npm contains invalid version",
			args:     []string{"npm", "contains", "^1.0.0", "invalid"},
			wantOut:  "Error running command 'contains': invalid version 'invalid': invalid NPM version: invalid",
			wantCode: 1,
		},
		{
			name:     "npm contains wrong number of args",
			args:     []string{"npm", "contains", "^1.0.0"},
			wantOut:  "Error running command 'contains': contains requires exactly 2 arguments: <version> <range>",
			wantCode: 1,
		},
		{
			name:     "alpine compare success less than",
			args:     []string{"alpine", "compare", "1.0.0", "2.0.0"},
			wantOut:  "-1",
			wantCode: 0,
		},
		{
			name:     "alpine compare success greater than",
			args:     []string{"alpine", "compare", "2.0.0", "1.0.0"},
			wantOut:  "1",
			wantCode: 0,
		},
		{
			name:     "alpine compare success equal",
			args:     []string{"alpine", "compare", "1.0.0", "1.0.0"},
			wantOut:  "0",
			wantCode: 0,
		},
		{
			name:     "alpine compare suffix vs release",
			args:     []string{"alpine", "compare", "1.0.0_alpha", "1.0.0"},
			wantOut:  "-1",
			wantCode: 0,
		},
		{
			name:     "alpine sort success",
			args:     []string{"alpine", "sort", "2.0.0", "1.0.0_alpha", "1.0.0"},
			wantOut:  "\"1.0.0_alpha\" \"1.0.0\" \"2.0.0\"",
			wantCode: 0,
		},
		{
			name:     "alpine contains true",
			args:     []string{"alpine", "contains", ">=1.0.0", "1.5.0"},
			wantOut:  "true",
			wantCode: 0,
		},
		{
			name:     "alpine contains false",
			args:     []string{"alpine", "contains", ">=2.0.0", "1.5.0"},
			wantOut:  "false",
			wantCode: 0,
		},
		{
			name:     "alpine contains invalid version",
			args:     []string{"alpine", "contains", ">=1.0.0", "invalid"},
			wantOut:  "Error running command 'contains': invalid version 'invalid': invalid Alpine version: invalid",
			wantCode: 1,
		},
		{
			name:     "pypi ecosystem success",
			args:     []string{"pypi", "compare", "1.0.0", "2.0.0"},
			wantOut:  "-1",
			wantCode: 0,
		},
		{
			name:     "go ecosystem success",
			args:     []string{"go", "compare", "v1.0.0", "v2.0.0"},
			wantOut:  "-1",
			wantCode: 0,
		},
		{
			name:     "maven compare success less than",
			args:     []string{"maven", "compare", "1.0.0-alpha", "1.0.0"},
			wantOut:  "-1",
			wantCode: 0,
		},
		{
			name:     "maven compare success greater than",
			args:     []string{"maven", "compare", "1.0.0-sp", "1.0.0"},
			wantOut:  "1",
			wantCode: 0,
		},
		{
			name:     "maven compare success equal",
			args:     []string{"maven", "compare", "1.0.0-ga", "1.0.0"},
			wantOut:  "0",
			wantCode: 0,
		},
		{
			name:     "maven sort success",
			args:     []string{"maven", "sort", "1.0.0", "1.0.0-alpha", "1.0.0-beta"},
			wantOut:  "\"1.0.0-alpha\" \"1.0.0-beta\" \"1.0.0\"",
			wantCode: 0,
		},
		{
			name:     "maven contains success true",
			args:     []string{"maven", "contains", "[1.0.0,2.0.0]", "1.5.0"},
			wantOut:  "true",
			wantCode: 0,
		},
		{
			name:     "maven contains success false",
			args:     []string{"maven", "contains", "[1.0.0]", "1.0.1"},
			wantOut:  "false",
			wantCode: 0,
		},
		{
			name:     "maven contains invalid version",
			args:     []string{"maven", "contains", "[1.0.0,2.0.0]", "invalid"},
			wantOut:  "Error running command 'contains': invalid version 'invalid': invalid Maven version format: invalid",
			wantCode: 1,
		},
		{
			name:     "gem compare success less than",
			args:     []string{"gem", "compare", "1.0.0", "2.0.0"},
			wantOut:  "-1",
			wantCode: 0,
		},
		{
			name:     "gem compare success greater than",
			args:     []string{"gem", "compare", "2.0.0", "1.0.0"},
			wantOut:  "1",
			wantCode: 0,
		},
		{
			name:     "gem compare success equal",
			args:     []string{"gem", "compare", "1.0.0", "1.0.0"},
			wantOut:  "0",
			wantCode: 0,
		},
		{
			name:     "gem compare prerelease vs release",
			args:     []string{"gem", "compare", "1.0.0-alpha", "1.0.0"},
			wantOut:  "-1",
			wantCode: 0,
		},
		{
			name:     "gem sort success",
			args:     []string{"gem", "sort", "2.0.0", "1.0.0-alpha", "1.0.0"},
			wantOut:  "\"1.0.0-alpha\" \"1.0.0\" \"2.0.0\"",
			wantCode: 0,
		},
		{
			name:     "gem contains pessimistic true",
			args:     []string{"gem", "contains", "~> 1.2.0", "1.2.5"},
			wantOut:  "true",
			wantCode: 0,
		},
		{
			name:     "gem contains pessimistic false",
			args:     []string{"gem", "contains", "~> 1.2.0", "1.3.0"},
			wantOut:  "false",
			wantCode: 0,
		},
		{
			name:     "gem contains invalid version",
			args:     []string{"gem", "contains", "~> 1.0.0", "invalid"},
			wantOut:  "Error running command 'contains': invalid version 'invalid': invalid Ruby Gem version: invalid",
			wantCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, gotCode := run(tt.args)
			if gotOut != tt.wantOut {
				t.Errorf("run(%+v) out = %q, want %q", tt.args, gotOut, tt.wantOut)
			}
			if gotCode != tt.wantCode {
				t.Errorf("run(%+v) code = %v, want %v", tt.args, gotCode, tt.wantCode)
			}
		})
	}
}

