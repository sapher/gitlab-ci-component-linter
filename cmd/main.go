package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sapher/gitlab-ci-component-linter/pkg/linter"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitlab-ci-component-linter [workdir]",
	Short: "Linter for Gitlab CI Component",
	Long: `Validate Gitlab CI Component files against set of rules
Workdir is the directory where the Gitlab CI Component project is located`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		cmdFlags := cmd.Flags()
		onlyFailures, _ := cmdFlags.GetBool("only-failures")
		softFail, _ := cmdFlags.GetBool("soft-fail")
		output, _ := cmdFlags.GetString("output")

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		// Get workdir
		workdir := args[0]
		if !filepath.IsAbs(workdir) {
			workdir = filepath.Join(cwd, args[0])
		}

		// Set workdir to current directory if not provided
		if workdir == "." {
			currentWorkDir, err := os.Getwd()
			if err != nil {
				os.Stderr.WriteString(err.Error())
				os.Exit(1)
			}
			workdir = currentWorkDir
		}

		// Check workdir validity
		if !linter.IsDirExist(workdir) {
			os.Stderr.WriteString(fmt.Sprintf("Workdir does not exist: %s", workdir))
			os.Exit(1)
		}

		// Create new linter
		newLinter := linter.New(workdir, linter.LinterOutput(output))

		// Execute linter
		ruleResults, err := newLinter.Execute()
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		// Filter results
		if onlyFailures {
			var filteredResults linter.LinterResults
			for _, result := range ruleResults {
				if !result.Success {
					filteredResults = append(filteredResults, result)
				}
			}
			ruleResults = filteredResults
		}

		// If no results, everything is fine, exit with 0
		if len(ruleResults) == 0 {
			os.Exit(0)
		}

		// Output results
		if newLinter.Output != linter.OutputNone {
			output, err = ruleResults.Output(newLinter.Output)
			if err != nil {
				os.Stderr.WriteString(err.Error())
				os.Exit(1)
			}
			os.Stdout.WriteString(output)
		}

		// Exit with error if has failulres and soft-fail is false
		if ruleResults.HasFailures() && !softFail {
			os.Exit(1)
		}
	},
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().Bool("only-failures", false, "Only display checks that failed")
	rootCmd.PersistentFlags().BoolP("soft-fail", "s", false, "Run checks and exit with 0 even if errors are found")
	rootCmd.PersistentFlags().StringP("output", "o", string(linter.OutputTable), "Output format, one of: json, yaml, junit, table, none")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
