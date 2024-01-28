package linter

import (
	_ "embed"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type JsonRule struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type JsonRuleMap struct {
	Rules map[string]JsonRule `json:"rules"`
}

//go:embed rules.json
var rawJsonRules string
var jsonRules JsonRuleMap

func init() {
	err := json.Unmarshal([]byte(rawJsonRules), &jsonRules)
	if err != nil {
		panic(err)
	}
}

type LinterRuleFunc func(linter *Linter) (LinterResult, error)

var ruleFuncs = []LinterRuleFunc{
	MissingRootReadmeRule,
	MissingRootTemplatesDirRule,
	WrongYamlFileExtensionRule,
}

func MissingRootReadmeRule(linter *Linter) (LinterResult, error) {
	ruleName := "missing-root-readme"
	yamlRule := jsonRules.Rules[ruleName]
	result := LinterResult{
		Name:     ruleName,
		Success:  false,
		Message:  yamlRule.Message,
		Metadata: map[string]interface{}{},
		Severity: RuleSeverity(yamlRule.Severity),
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
	ruleName := "missing-root-templates-dir"
	yamlRule := jsonRules.Rules[ruleName]
	result := LinterResult{
		Name:     ruleName,
		Success:  false,
		Message:  yamlRule.Message,
		Metadata: map[string]interface{}{},
		Severity: RuleSeverity(yamlRule.Severity),
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
	ruleName := "wrong-yaml-file-extension"
	yamlRule := jsonRules.Rules[ruleName]
	result := LinterResult{
		Name:     ruleName,
		Success:  false,
		Message:  yamlRule.Message,
		Metadata: map[string]interface{}{},
		Severity: RuleSeverity(yamlRule.Severity),
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
