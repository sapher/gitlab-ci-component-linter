package main

import (
	"fmt"
	"os"

	"github.com/sapher/gitlab-ci-component-linter/pkg/linter"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitlab-ci-component-linter",
	Short: "Linter for Gitlab CI Component",
	Long:  `Validate Gitlab CI Component files against set of rules`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create new linter
		cmdFlags := cmd.Flags()
		workdir, _ := cmdFlags.GetString("workdir")
		softFail, _ := cmdFlags.GetBool("soft-fail")
		output, _ := cmdFlags.GetString("output")
		newLinter := linter.New(workdir, linter.LinterOutput(output))

		// Check workdir validity
		if linter.IsDirExist(workdir) {
			os.Stderr.WriteString(fmt.Sprintf("Workdir does not exist: %s", workdir))
			os.Exit(1)
		}

		// Execute linter
		ruleResults, err := newLinter.Execute()
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		// Has rules in failure
		hasErrors := false
		for _, ruleResult := range ruleResults {
			if !ruleResult.Success && ruleResult.Severity == linter.SeverityError {
				hasErrors = true
			}
		}

		// Output results
		output, err = ruleResults.Output(newLinter.Output)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		os.Stdout.WriteString(output)

		// Exit with error if has errors and soft-fail is false
		if hasErrors && !softFail {
			os.Exit(1)
		}
	},
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().String("workdir", ".", "Working directory")
	rootCmd.PersistentFlags().Bool("soft-fail", false, "Wether to fail or not on error")
	rootCmd.PersistentFlags().String("output", string(linter.OutputJson), "Output format, one of: json, yaml, junit")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
