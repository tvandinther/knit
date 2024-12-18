package helm

import (
	"fmt"
)

type JsonSchema struct {
	Type        string                 `json:"type,omitempty"`
	Default     interface{}            `json:"default,omitempty"`
	Description string                 `json:"description,omitempty"`
	Properties  map[string]*JsonSchema `json:"properties,omitempty"`
	Items       *JsonSchema            `json:"items,omitempty"`
	Required    []string			   `json:"required,omitempty"`
}

func ValuesNodeToJsonSchema(root *ValuesNode) (*JsonSchema, error) {
	return toJSONSchema(root)
}

var emptyObject = map[string]any{}
var emptyArray = []any{}

func toJSONSchema(node *ValuesNode) (*JsonSchema, error) {
	nodeType, err := mapTypeToSchema(node.Type)
	if err != nil {
		return nil, fmt.Errorf("Could not convert values to json schema", err)
	}
	
	schema := &JsonSchema{
		Type:        nodeType,
		Default:     node.Value,
		Description: node.Comments,
	}

	switch node.Type {
	case TypeObject:
		schema.Properties = make(map[string]*JsonSchema)
		if node.SubNodes == nil {
			schema.Default = emptyObject
		}
		for _, child := range node.SubNodes {
			s, err := toJSONSchema(child)
			if err != nil {
				return nil, fmt.Errorf("Could not convert values to json schema", err)
			}
			schema.Properties[child.Name] = s
			if s.Default != nil {
				schema.Required = append(schema.Required, child.Name)
			}
		}
	case TypeArray:
		// Assuming all items in the array have the same schema
		if len(node.SubNodes) > 0 {
			s, err := toJSONSchema(node.SubNodes[0])
			if err != nil {
				return nil, fmt.Errorf("Could not convert values to json schema", err)
			}
			schema.Items = s
		} else {
			schema.Items = &JsonSchema{}
			schema.Default = emptyArray
		}
	}

	return schema, nil
}

func mapTypeToSchema(t JSONType) (string, error) {
	switch t {
	case TypeNull:
		return "null", nil
	case TypeString:
		return "string", nil
	case TypeInteger:
		return "integer", nil
	case TypeNumber:
		return "number", nil
	case TypeBool:
		return "boolean", nil
	case TypeArray:
		return "array", nil
	case TypeObject:
		return "object", nil
	default:
		return "", fmt.Errorf("Could not map type to schema. Possible partial mapping")
	}
}
