/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/env"
	"kcl-lang.io/kpm/pkg/opt"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/reporter"
)

var initOpts = opt.InitOptions{}
// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit(args)
	},
}

func runInit(args []string) error {
	pwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	var pkgName string
	var pkgRootPath string
	// 1. If no package name is given, the current directory name is used as the package name.
	if len(args) == 0 {
		pkgName = filepath.Base(pwd)
		pkgRootPath = pwd
	} else {
		// 2. If the package name is given, create a new directory for the package.
		pkgName = argsGet(args, 0)
		pkgRootPath = filepath.Join(pwd, pkgName)
		err = os.MkdirAll(pkgRootPath, 0755)
		if err != nil {
			return err
		}
	}

	initOpts.Name = pkgName
	initOpts.InitPath = pkgRootPath

	err = initOpts.Validate()
	if err != nil {
		return err
	}

	kclPkg := pkg.NewKclPkg(&initOpts)

	globalPkgPath, err := env.GetAbsPkgPath()

	if err != nil {
		return err
	}

	err = kclPkg.ValidateKpmHome(globalPkgPath)

	if err != (*reporter.KpmEvent)(nil) {
		return err
	}

	cli, err := client.NewKpmClient()
	if err != nil {
		return err
	}
	err = cli.InitEmptyPkg(&kclPkg)
	if err != nil {
		return err
	}

	fmt.Printf("package '%s' init finished\n", pkgName)
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.
	initCmd.Flags().StringVar(&initOpts.Version, "version", "", "kcl mod init module version")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
