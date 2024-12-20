package cmd

import (
	"github.com/spf13/cobra"

	"knit/pkg/render"
)

var renderCmd = &cobra.Command{
	Use:   "render [file]",
	Short: "Render a KCL file as YAML",
	Long: `Renders a KCL file as YAML to stdout. 'main.k' in the KCL module root will be rendered unless an argument is given.
Filepaths are resolved relative to the project root defined by the location of the 'kcl.mod' file if found or relative otherwise.

Examples:
	# Renders 'main.k' as YAML to stdout
	knit render
	
	# Renders 'dev.k' as YAML to stdout
	knit render dev.k
	
	# You can also render from the top level command
	knit
	knit dev.k`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := argsGet(args, 0)
		if filepath == "" {
			filepath = "main.k"
		}
		cmd.SilenceUsage = true
		return render.Render(filepath)
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
}
