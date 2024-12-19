package helm

import (
	"encoding/json"
	"fmt"
	"knit/pkg/util"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"kcl-lang.io/kcl-go"
	"kcl-lang.io/kcl-go/pkg/tools/gen"
)

type JSONType int

const (
	TypeNull JSONType = iota
	TypeString
	TypeInteger
	TypeNumber
	TypeBool
	TypeArray
	TypeObject
)

type ValuesNode struct {
	Type     JSONType
	Name     string
	Value    interface{}
	Comments string
	SubNodes []*ValuesNode
}

func RunValues(chartRef *ChartRef) error {
	return getValues(chartRef)
}

func getValues(chartRef *ChartRef) error {
	settings := cli.New()

	chart, err := getChart(chartRef, settings)
	if err != nil {
		return fmt.Errorf("could not get helm chart: %w", err)
	}

	valuesFile, err := getValuesFile(chart)
	if err != nil {
		return err
	}

	var root yaml.Node
	err = yaml.Unmarshal(valuesFile.Data, &root)
	if err != nil {
		return fmt.Errorf("could not unmarshal values yaml: %w", err)
	}

	if len(root.Content) != 1 {
		return fmt.Errorf("expecting a single document in values.yaml")
	}

	tree, err := parseYAML(root.Content[0], "root")
	if err != nil {
		return fmt.Errorf("could not parse the values yaml document: %w", err)
	}

	schema, err := ValuesNodeToJsonSchema(tree)
	if err != nil {
		return err
	}
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
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

	directory := filepath.Join("vendored", "helm", "podinfo")
	err = os.MkdirAll(directory, 0744)
	if err != nil {
		return err
	}
	valuesFilepath := filepath.Join(directory, "values.k")
	f, err := os.Create(valuesFilepath)
	if err != nil {
		return err
	}
	gen.GenKcl(f, valuesSchemaFile.Name(), nil, &gen.GenKclOptions{Mode: gen.ModeJsonSchema})

	valuesFilepath = filepath.Join(directory, "chart.k")
	f, err = os.Create(valuesFilepath)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, chartFileContent, chartRef.Repository, chartRef.Repository, chartRef.Name, chartRef.Name, chartRef.Version, chartRef.Version)
	_, err = kcl.FormatPath(directory)
	if err != nil {
		return err
	}

	return nil
}

func getValuesFile(chart *chart.Chart) (*chart.File, error) {
	for _, file := range chart.Raw {
		if file.Name == "values.yaml" {
			return file, nil
		}
	}

	return nil, fmt.Errorf("could not find default values.yaml file")
}

func collectComments(nodes ...*yaml.Node) string {
	comments := make([]string, 0)

	addComment := func(c string) {
		if c != "" {
			s := strings.TrimPrefix(c, "#")
			s = strings.TrimSpace(s)
			comments = append(comments, s)
		}
	}

	for _, node := range nodes {
		addComment(node.HeadComment)
		addComment(node.LineComment)
		addComment(node.FootComment)
	}

	return strings.Join(comments, "\n")
}

func parseYAML(node *yaml.Node, name string) (*ValuesNode, error) {
	vNode := &ValuesNode{
		Name: name,
	}

	switch node.Kind {
	case yaml.ScalarNode:
		err := node.Decode(&vNode.Value)
		vNode.Comments = collectComments(node)
		if err != nil {
			return nil, fmt.Errorf("could not decode yaml value: %w", err)
		}
		switch node.Tag {
		case "!!str":
			vNode.Type = TypeString
		case "!!int":
			vNode.Type = TypeInteger
		case "!!float":
			vNode.Type = TypeNumber
		case "!!bool":
			vNode.Type = TypeBool
		case "!!null":
			vNode.Type = TypeNull
		default:
			return nil, fmt.Errorf("unknown scalar type: %s", node.Tag)
		}
	case yaml.SequenceNode:
		vNode.Type = TypeArray
		for _, subNode := range node.Content {
			child, err := parseYAML(subNode, "")
			if err != nil {
				return nil, err
			}
			child.Comments = collectComments(subNode)
			vNode.SubNodes = append(vNode.SubNodes, child)
		}
	case yaml.MappingNode:
		vNode.Type = TypeObject
		for i := 0; i < len(node.Content); i += 2 {
			keyNode, valueNode := node.Content[i], node.Content[i+1]
			if keyNode.Kind != yaml.ScalarNode {
				return nil, fmt.Errorf("non-scalar key in mapping")
			}
			child, err := parseYAML(valueNode, keyNode.Value)
			child.Comments = collectComments(keyNode, valueNode)
			if err != nil {
				return nil, err
			}
			vNode.SubNodes = append(vNode.SubNodes, child)
		}
	case yaml.AliasNode:
		return parseYAML(node.Alias, name)
	default:
		return nil, fmt.Errorf("unsupported node kind: %d", node.Kind)
	}

	return vNode, nil
}
