package render

import (
	"fmt"

	_ "knit/http"

	"kcl-lang.io/kcl-go/pkg/kcl"                   // Import the native API
	_ "kcl-lang.io/kcl-go/pkg/plugin/hello_plugin" // Import the hello plugin
)

func Render() {
	yaml := kcl.MustRun("main.k", kcl.WithCode(code)).GetRawYamlResult()
	fmt.Println(yaml)
}

const code = `
import kcl_plugin.hello
import kcl_plugin.http

name = "kcl"
three = hello.add(1,2)  # hello.add is written by Go
http_status = http.get("https://google.com").status
`
