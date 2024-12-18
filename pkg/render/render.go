package render

import (
	"fmt"

	_ "knit/http"

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

const code = `
import kcl_plugin.hello
import kcl_plugin.http

name = "kcl"
three = hello.add(1,2)  # hello.add is written by Go
http_status = http.get("https://google.com").status
`
