package util

import (
	"fmt"
	"knit/pkg/helm"
	"slices"
	"strings"
)

type KCLWriter struct {
	filepath string
	sb       strings.Builder
}

type KCLSchema struct {
	Name        string
	Inheritance string
	Properties  []KCLSchemaProperty
	Description string
}

type KCLTypeDiscriminator int

const (
	TypeNone KCLTypeDiscriminator = iota
	TypeString
	TypeInt
	TypeFloat
	TypeBool
	TypeList
	TypeDict
)

type KCLType struct {
	Discriminator KCLTypeDiscriminator
	ChildrenTypes []KCLType
	KeyType       *KCLType
	ValueType     *KCLType
}

type KCLSchemaProperty struct {
	Name         string
	Types        []KCLType
	Optional     bool
	DefaultValue any
	Description  string
}

const docstringDelimiter = `"""`
const indentation = `    `

func New(filepath string) *KCLWriter {
	return &KCLWriter{
		filepath: filepath,
	}
}

func (k *KCLWriter) AddImport(path string) {
	k.sb.WriteString(fmt.Sprintf("import %s\n", path))
}

func (k *KCLWriter) AddDocstring(contents string) {
	k.sb.WriteString(fmt.Sprintf("%s\n%s\n%s\n", docstringDelimiter, contents, docstringDelimiter))
}

func (k *KCLWriter) AddRawDocstring(contents string) {
	k.sb.WriteString(fmt.Sprintf("r%s\n%s\n%s\n", docstringDelimiter, contents, docstringDelimiter))
}

func (k *KCLWriter) AddSchema(schema KCLSchema) {
	k.startSchema(schema)
	for _, property := range schema.Properties {
		k.addSchemaProperty(property)
	}
}

func (k *KCLWriter) startSchema(schema KCLSchema) {
	k.sb.WriteString(fmt.Sprintf("schema %s", schema.Name))
	if schema.Inheritance != "" {
		k.sb.WriteString(fmt.Sprintf("(%s)", schema.Inheritance))
	}
	k.sb.WriteString(":\n")
}

func (k *KCLWriter) addSchemaProperty(property KCLSchemaProperty) {
	k.sb.WriteString(indentation)
	k.sb.WriteString(property.Name)
	if property.Optional {
		k.sb.WriteString("?")
	}
	k.sb.WriteString(": ")
	// k.sb.WriteString(typeName)
	if !property.Optional {
		k.sb.WriteString(" = ")
		// k.sb.WriteString(defaultValue)
	}
	k.sb.WriteString("\n")
}

func JSONSchemaToKCLSchema(name string, jsonSchema helm.JsonSchema) (*KCLSchema, error) {
	kclSchema := KCLSchema{
		Name: fmt.Sprintf("%sValues", name),
	}

	if name == "" {
		kclSchema.Inheritance = "helm.Values"
	}

	kclSchema.Description = jsonSchema.Description

	kclSchema.Properties = make([]KCLSchemaProperty, 0, len(jsonSchema.Properties))
	for key, jsonProperty := range jsonSchema.Properties {
		kclProperty := KCLSchemaProperty{
			Name:         key,
			Description:  jsonProperty.Description,
			Optional:     !slices.Contains(jsonSchema.Required, key),
			DefaultValue: jsonSchema.Default,
		}

		var propertyTypes []string

		propertyType, ok := jsonProperty.Type.(string)
		if ok {
			propertyTypes = append(propertyTypes, propertyType)
		} else {
			propertyTypes, ok = jsonProperty.Type.([]string)
			if !ok {
				return nil, fmt.Errorf("internal error: invalid property type")
			}
		}

		for _, propType := range propertyTypes {
			discriminator := jsonSchemaTypeToKCLTypeDiscriminator(propType)

			kclType := KCLType{
				Discriminator: discriminator,
			}

			if discriminator == TypeList {
				jsonProperty.Items.Type
			}
		}
	}

	return &kclSchema
}

func makeKclType(propertyType string, parentProperty helm.JsonSchema) *KCLType {
	discriminator := jsonSchemaTypeToKCLTypeDiscriminator(propType)

	kclType := KCLType{
		Discriminator: discriminator,
	}

	if discriminator == TypeList {
		kclType.ChildrenTypes = makeKclType(parentProperty.Items.Type)
	}
}

func jsonSchemaTypeToKCLTypeDiscriminator(t string) KCLTypeDiscriminator {
	switch t {
	case "array":
		return TypeList
	case "object":
		return TypeDict
	case "string":
		return TypeString
	case "number":
		return TypeFloat
		// return TypeInt
	case "null":
		return TypeNone
	case "boolean":
		return TypeBool
	default:
		panic(fmt.Sprintf("non-exhastive pattern match on %T", t))
	}
}
