package initialise

import (
	"os"
	"path/filepath"

	"kcl-lang.io/kcl-go"
)

func WriteFiles(pkgRootPath string) error {
	directory := filepath.Join(pkgRootPath, "knit")
	err := os.MkdirAll(directory, 0744)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(directory, "kubernetes.k"), []byte(kubernetesBaseFileContent), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(directory, "helm.k"), []byte(helmBaseFileContent), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(directory, "kustomize.k"), []byte(kustomizeFileContent), 0644)
	if err != nil {
		return err
	}

	_, err = kcl.FormatPath(directory)
	if err != nil {
		return err
	}

	return nil
}
