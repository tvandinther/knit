package render

import (
	"fmt"

	_ "knit/pkg/plugin"
	"knit/pkg/util"

	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/opt"

	"path/filepath"
)

func Render(file string) error {
	moduleRoot, err := util.FindModuleRoot()
	if err != nil {
		return err
	}

	client, err := client.NewKpmClient()
	if err != nil {
		return err
	}
	err = client.AcquirePackageCacheLock()
	if err != nil {
		return err
	}
	defer func() {
		releaseErr := client.ReleasePackageCacheLock()
		if releaseErr != nil && err == nil {
			err = releaseErr
		}
	}()

	entries := []string{filepath.Join(moduleRoot, file)}
	result, err := client.RunWithOpts(opt.WithEntries(entries))
	if err != nil {
		return err
	}
	yaml := result.GetRawYamlResult()
	fmt.Println(yaml)

	return nil
}
