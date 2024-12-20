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

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a project in the current directory",
	Long: `Initialises a project in the current directory. Running this command will initialise a KCL module just like running 'kcl mod init' would.

Example:
	# Initialises a project versioned as 1.0.0
	knit init --version 1.0.0`,
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

	initCmd.Flags().StringVar(&initOpts.Version, "version", "", "kcl mod init module version")
}
