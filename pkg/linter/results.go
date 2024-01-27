package linter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"gopkg.in/yaml.v3"
)

type JunitTestCaseFailure struct {
	Message string `xml:",chardata"`
}

type JunitTestCase struct {
	Name    string                 `xml:"name,attr"`
	Failure []JunitTestCaseFailure `xml:"failure,omitempty"`
}

type JunitTestSuite struct {
	Testcases []JunitTestCase `xml:"testcase"`
	Name      string          `xml:"name,attr"`
	Id        int             `xml:"id,attr"`
	Disabled  int             `xml:"disabled,attr"`
	Skipped   int             `xml:"skipped,attr"`
	Errors    int             `xml:"errors,attr"`
	Failures  int             `xml:"failures,attr"`
	Tests     int             `xml:"tests,attr"`
}

type JunitTestSuites struct {
	XMLName    xml.Name         `xml:"testsuites"`
	TestSuites []JunitTestSuite `xml:"testsuite"`
}

type LinterOutput string

const (
	OutputJson        LinterOutput = "json"
	OutputYaml        LinterOutput = "yaml"
	OutputJunitReport LinterOutput = "junit"
)

type RuleSeverity string

const (
	SeverityError   RuleSeverity = "error"
	SeverityWarning RuleSeverity = "warning"
)

type LinterResult struct {
	Name     string                 `json:"name"`
	Success  bool                   `json:"success"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata"`
	Severity RuleSeverity           `json:"severity"`
}

type LinterResults []LinterResult

func (l *LinterResults) Output(output LinterOutput) (string, error) {
	switch output {
	case OutputJson:
		return l.ToJson()
	case OutputYaml:
		return l.ToYaml()
	case OutputJunitReport:
		return l.ToJunitReport()
	default:
		return "", fmt.Errorf("unknown output format: %s", output)
	}
}

func (l *LinterResults) ToJson() (string, error) {
	output, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (l *LinterResults) ToYaml() (string, error) {
	output, err := yaml.Marshal(l)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (l *LinterResults) ToJunitReport() (string, error) {
	// Counts errors and failures
	errors := 0
	failures := 0
	tests := len(*l)

	for _, result := range *l {
		if !result.Success {
			switch result.Severity {
			case SeverityError:
				errors++
			case SeverityWarning:
				failures++
			}
		}
	}

	// Build test cases
	testCases := []JunitTestCase{}
	for _, result := range *l {
		testCase := JunitTestCase{
			Name:    result.Name,
			Failure: []JunitTestCaseFailure{},
		}
		if !result.Success {
			testCase.Failure = append(testCase.Failure, JunitTestCaseFailure{
				Message: result.Message,
			})
		}
		testCases = append(testCases, testCase)
	}

	// Build test suites
	testSuites := JunitTestSuites{
		TestSuites: []JunitTestSuite{
			{
				Name:      "Gitlab CI Component Linter",
				Id:        0,
				Testcases: testCases,
				Disabled:  0, // TODO: Implement later
				Skipped:   0, // TODO: Implement later
				Errors:    errors,
				Failures:  failures,
				Tests:     tests,
			},
		},
	}

	output, err := xml.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(output), nil
}
