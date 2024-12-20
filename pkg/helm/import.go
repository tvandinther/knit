package helm

import (
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

	fmt.Printf("Helm chart successfully imported to %s\n", chartDirectory)
	fmt.Printf(`Example usage:
	
import %s
import kcl_plugin.helm

_chart = %s.Chart {
    values = %s.Values {
        # Typed Helm values map
    }
}

manifests.yaml_stream(helm.template(_chart)) 
`, strings.ReplaceAll(chartDirectory, string(filepath.Separator), "."), filepath.Base(chartDirectory), filepath.Base(chartDirectory))

	return nil
}
