package linter

import (
	"strings"
	"testing"
)

func TestLinterResultsOutput(t *testing.T) {
	// A sample result set with both successful and failed results for testing
	results := LinterResults{
		{Name: "TestRule1", Success: true, Severity: SeverityWarning},
		{Name: "TestRule2", Success: false, Severity: SeverityError, Message: "Error message"},
		{Name: "TestRule3", Success: false, Severity: SeverityWarning, Message: "Warning message"},
	}

	// Define test cases
	tests := []struct {
		name         string
		output       LinterOutput
		wantContains []string // strings that should be present in output
		wantErr      bool     // whether an error is expected
	}{
		{"Output JSON", OutputJson, []string{`"name": "TestRule1"`, `"severity": "warning"`}, false},
		{"Output YAML", OutputYaml, []string{"name: TestRule1", "severity: warning"}, false},
		{"Output JUnit", OutputJunitReport, []string{"<testsuite ", "<testcase ", "Error message", "Warning message"}, false},
		{"Output Unknown", LinterOutput("unknown"), nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := results.Output(tc.output)
			if (err != nil) != tc.wantErr {
				t.Fatalf("results.Output() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil {
				for _, wantStr := range tc.wantContains {
					if !strings.Contains(got, wantStr) {
						t.Errorf("results.Output() = %v, want it to contain %v", got, wantStr)
					}
				}
			}
		})
	}
}
