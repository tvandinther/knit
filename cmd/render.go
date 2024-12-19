/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"knit/pkg/render"

	// Import the native API
	_ "kcl-lang.io/kcl-go/pkg/plugin/hello_plugin" // Import the hello plugin
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render [file]",
	Short: "Render a KCL file as YAML",
	Long: `Renders a KCL file as YAML to stdout. 'main.k' in the KCL module root will be rendered unless an argument is given.
Filepaths are resolved relative to the project root defined by the location of the 'kcl.mod' file if found or relative otherwise.

Examples:
	# Renders 'main.k' as YAML to stdout
	knit render
	
	# Renders 'dev.k' as YAML to stdout
	knit render dev.k`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := argsGet(args, 0)
		if filepath == "" {
			filepath = "main.k"
		}

		return render.Render(filepath)
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
