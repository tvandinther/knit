package helm

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
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
	Type       JSONType
	Name       string
	Value      interface{}
	Comments   string
	SubNodes   []*ValuesNode
}

func RunValues(chartRef *ChartRef) error {
	return getValues(chartRef)
}

func getValues(chartRef *ChartRef) error {
	settings := cli.New()

	chart, err := getChart(chartRef, settings)
	if err != nil {
		return fmt.Errorf("Could not get helm chart", err)
	}

	valuesFile, err := getValuesFile(chart)
	if err != nil {
		return err
	}

	// TODO: Convert yaml data into a yaml AST with comments included and relate lines of config with comments by building up a JSON schema for KCL to import as a KCL schema (https://pkg.go.dev/kcl-lang.io/kcl-go@v0.10.8/pkg/tools/gen#GenKcl)

	var root yaml.Node
	err = yaml.Unmarshal(valuesFile.Data, &root)
	if err != nil {
		return fmt.Errorf("Could not unmarshal values yaml", err)
	}

	if len(root.Content) != 1 {
		return fmt.Errorf("Expecting a single document in values.yaml")
	}

	tree, err := parseYAML(root.Content[0], "root")
	if err != nil {
		return fmt.Errorf("Could not parse the values yaml document", err)
	}

	schema, err := ValuesNodeToJsonSchema(tree)
	if err != nil {
		return err
	}
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}
	os.WriteFile("values.json", schemaJSON, 0644)

	directory := path.Join("vendored", "helm", "podinfo") 
	err = os.MkdirAll(directory, 0744)
	if err != nil {
		return err
	}
	filepath := path.Join(directory, "values.k")
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	gen.GenKcl(f, "values.json", nil, &gen.GenKclOptions{Mode: gen.ModeJsonSchema})

	return nil
}

func getValuesFile(chart *chart.Chart) (*chart.File, error) {
	for _, file := range chart.Raw {
		if file.Name == "values.yaml" {
			return file, nil
		}
	}

	return nil, fmt.Errorf("Could not find default values.yaml file")
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
		Name:     name,
	}

	switch node.Kind {
	case yaml.ScalarNode:
		err := node.Decode(&vNode.Value)
		vNode.Comments = collectComments(node)
		if err != nil {
			return nil, fmt.Errorf("Could not decode yaml value", err)
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

