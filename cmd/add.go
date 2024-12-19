package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a component to your project",
	Long: `Add a component such as a helm chart to your KCL project.
	
Example:
	# Adds a vendored helm chart to your KCL project
	knit add helm https://stefanprodan.github.io/podinfo podinfo --version 6.7.1`,
}

func init() {
	rootCmd.AddCommand(addCmd)
}
