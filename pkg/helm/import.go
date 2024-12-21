package helm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"knit/pkg/util"
	"os"
	"path/filepath"
	"strings"

	"kcl-lang.io/kcl-go"
	"kcl-lang.io/kcl-go/pkg/tools/gen"
)

func Import(chartRef *ChartRef, directory string) error {
	valuesNode, err := getValues(chartRef)
	if err != nil {
		return err
	}

	schema, err := ValuesNodeToJsonSchema(valuesNode)
	if err != nil {
		return err
	}

	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	tmpDir, err := util.NewTempDir(fmt.Sprintf("values.%s-%s_", chartRef.Name, chartRef.Version))
	if err != nil {
		return err
	}
	defer tmpDir.Remove()

	valuesSchemaFile, err := tmpDir.CreateFile("values.json")
	if err != nil {
		return err
	}

	valuesSchemaFile.Write(schemaJSON)

	chartDirectory := filepath.Join(directory, strings.ReplaceAll(chartRef.Name, "-", "_"))
	err = os.MkdirAll(chartDirectory, 0744)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(chartDirectory, "values.k"))
	if err != nil {
		return err
	}

	var buf = bytes.NewBuffer([]byte("import knit.helm"))
	err = gen.GenKcl(buf, valuesSchemaFile.Name(), nil, &gen.GenKclOptions{Mode: gen.ModeJsonSchema})
	if err != nil {
		return err
	}
	enhancedValuesSchema := strings.Replace(buf.String(), "schema Values:", "schema Values(helm.Values):", 1)
	_, err = f.Write([]byte(enhancedValuesSchema))
	if err != nil {
		return err
	}

	f, err = os.Create(filepath.Join(chartDirectory, "chart.k"))
	if err != nil {
		return err
	}
	fmt.Fprintf(f, chartFileContent, chartRef.Repository, chartRef.Repository, chartRef.Name, chartRef.Name, chartRef.Version, chartRef.Version)
	_, err = kcl.FormatPath(chartDirectory)
	if err != nil {
		return err
	}

	fmt.Printf("Helm chart successfully imported to %s\n", chartDirectory)
	fmt.Printf("To import this chart use: `import %s`\n", strings.ReplaceAll(chartDirectory, string(filepath.Separator), "."))

	return nil
}
