package cmd

import (
	"fmt"
	"knit/pkg/util"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "knit",
	Version: fmt.Sprintf("%s-%s #%s", util.GetVersion(), util.GetArchitecture(), util.GetShortHash()),
	Short:   "Knits together kubernetes manifests using KCL",
	Long: `knit helps you to create kubernetes manifests from various sources and produce them as rendered artifacts to easily enable the rendered manifests pattern using KCL.
'knit' is also an alias to 'knit render' when no valid commands are given.

Examples:
	# Initialises a project versioned as 1.0.0
	knit init --version 1.0.0
	
	# Adds a vendored helm chart to your KCL project
	knit add helm https://stefanprodan.github.io/podinfo podinfo --version 6.7.1
	
	# Renders 'main.k' as YAML to stdout
	knit render
	knit # When no valid commands are given, render is executed as an alias`,
	Args:    renderCmd.Args,
	Aliases: []string{"render"},
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

}

func argsGet(a []string, n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}
