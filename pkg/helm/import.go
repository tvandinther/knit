package helm

import (
	"encoding/json"
	"fmt"
	"knit/pkg/util"
	"os"
	"path/filepath"

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

	chartDirectory := filepath.Join(directory, chartRef.Name)
	err = os.MkdirAll(chartDirectory, 0744)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(chartDirectory, "values.k"))
	if err != nil {
		return err
	}

	gen.GenKcl(f, valuesSchemaFile.Name(), nil, &gen.GenKclOptions{Mode: gen.ModeJsonSchema})

	f, err = os.Create(filepath.Join(chartDirectory, "chart.k"))
	if err != nil {
		return err
	}
	fmt.Fprintf(f, chartFileContent, chartRef.Repository, chartRef.Repository, chartRef.Name, chartRef.Name, chartRef.Version, chartRef.Version)
	_, err = kcl.FormatPath(chartDirectory)
	if err != nil {
		return err
	}

	return nil
}
