package render

import (
	"fmt"

	_ "knit/pkg/plugin"
	"knit/pkg/util"

	"kcl-lang.io/kcl-go/pkg/kcl"

	"path/filepath"
)

func Render(file string) error {
	moduleRoot, err := util.FindModuleRoot()
	if err != nil {
		return err
	}

	result, err := kcl.Run(filepath.Join(moduleRoot, file))
	if err != nil {
		return err
	}
	yaml := result.GetRawYamlResult()
	fmt.Println(yaml)

	return nil
}
