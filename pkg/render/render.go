package render

import (
	"fmt"

	_ "knit/pkg/plugin"

	"kcl-lang.io/kcl-go/pkg/kcl"                   // Import the native API
	_ "kcl-lang.io/kcl-go/pkg/plugin/hello_plugin" // Import the hello plugin
)

func Render() error {
	result, err := kcl.Run("main.k")
	if err != nil {
		return err
	}
	yaml := result.GetRawYamlResult()
	fmt.Println(yaml)

	return nil
}
