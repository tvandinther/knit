package cmd

import (
	"knit/pkg/helm"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	version   string
	directory string
)

var helmCmd = &cobra.Command{
	Use:   "helm [repository] [name]",
	Short: "Adds a helm chart to your project",
	Long: `Add a helm chart to your project ready to be imported in a KCL file.

Example:
	# Adds a vendored helm chart to your KCL project
	knit add helm https://stefanprodan.github.io/podinfo podinfo --version 6.7.1

You can then import the podinfo chart from vendored/helm/podinfo.`,
	Args: cobra.MatchAll(cobra.ExactArgs(2)),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := helm.Import(&helm.ChartRef{
			Repository: argsGet(args, 0),
			Name:       argsGet(args, 1),
			Version:    version,
		}, directory)
		return err
	},
}

func init() {
	addCmd.AddCommand(helmCmd)

	helmCmd.Flags().StringVarP(&version, "version", "v", "", "version of the helm chart")
	helmCmd.MarkFlagRequired("version")
	helmCmd.Flags().StringVar(&directory, "dir", filepath.Join("vendored", "helm"), "directory to add helm chart configuration")
}
