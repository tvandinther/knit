/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"knit/pkg/util"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "knit",
	Version: fmt.Sprintf("%s #%s", util.GetVersion(), util.GetShortHash()),
	Short:   "Knits together kubernetes manifests using KCL",
	Long: `knit helps you to create kubernetes manifests from various sources and produce them as rendered artifacts to easily enable the rendered manifests pattern using KCL.

Examples:
	# Initialises a project versioned as 1.0.0
	knit init --version 1.0.0
	
	# Adds a vendored helm chart to your KCL project
	knit add helm https://stefanprodan.github.io/podinfo podinfo --version 6.7.1
	
	# Renders 'main.k' as YAML to stdout
	knit render`,
	Args: renderCmd.Args,
	RunE: func(cmd *cobra.Command, args []string) error {
		return renderCmd.RunE(cmd, args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func argsGet(a []string, n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}
