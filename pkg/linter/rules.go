package linter

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type JsonRule struct {
	Name     string `json:"name"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

//go:embed rules.json
var rawJsonRules string

type LinterRuleFunc func(linter *Linter) (LinterResult, error)

var ruleFuncs = []LinterRuleFunc{
	MissingRootReadmeRule,
	MissingRootTemplatesDirRule,
	WrongYamlFileExtensionRule,
}

func ParseJsonRules() (map[string]JsonRule, error) {
	var jsonRules map[string]JsonRule
	// Unmarshal JSON rules
	if err := json.Unmarshal([]byte(rawJsonRules), &jsonRules); err != nil {
		return nil, err
	}
	return jsonRules, nil
}

func GetJsonRule(ruleName string) (JsonRule, error) {
	// Unmarshal JSON rules
	jsonRules, err := ParseJsonRules()
	if err != nil {
		return JsonRule{}, err
	}
	// Check if rule exists
	if rule, ok := jsonRules[ruleName]; ok {
		return rule, nil
	}
	return JsonRule{}, fmt.Errorf("rule not found: %s", ruleName)
}

func MissingRootReadmeRule(linter *Linter) (LinterResult, error) {
	result := LinterResult{
		Name:     "missing-root-readme",
		Success:  false,
		Message:  "No README.md file found in root directory",
		Metadata: map[string]interface{}{},
		Severity: SeverityWarning,
	}

	// Check if README.md exists in root directory
	filepath := path.Join(linter.Workdir, "README.md")
	if !IsFileExist(filepath) {
		return result, nil
	}

	result.Success = true
	return result, nil
}

func MissingRootTemplatesDirRule(linter *Linter) (LinterResult, error) {
	result := LinterResult{
		Name:     "missing-root-templates-dir",
		Success:  false,
		Message:  "No templates directory found in root directory",
		Metadata: map[string]interface{}{},
		Severity: SeverityError,
	}

	// Check if templates directory exists in root directory
	dirpath := path.Join(linter.Workdir, "templates")
	if !IsDirExist(dirpath) {
		return result, nil
	}

	result.Success = true
	return result, nil
}

func WrongYamlFileExtensionRule(linter *Linter) (LinterResult, error) {
	result := LinterResult{
		Name:    "wrong-yaml-file-extension",
		Success: false,
		Message: "YAML files must have .yml extension, not .yaml",
		Metadata: map[string]interface{}{
			"files": []string{},
		},
		Severity: SeverityError,
	}

	// Check if any files with .yaml extension exists in root directory and subdirectories
	err := filepath.Walk(linter.Workdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if path is a file and has .yaml extension
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".yaml") {
			result.Metadata["files"] = append(result.Metadata["files"].([]string), strings.TrimPrefix(path, linter.Workdir+"/"))
			result.Success = false
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	result.Success = true
	return result, nil
}
