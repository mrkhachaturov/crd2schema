// Package convert rewrites CustomResourceDefinition schemas into JSON Schema
// draft-4.
package convert

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Schema converts a CRD structural schema in place.
func Schema(schema *apiextensionsv1.JSONSchemaProps) {
	if schema == nil {
		return
	}

	denyUnknownFields(schema)
	expandIntOrString(schema)
	expandNullable(schema)
	dropExtensions(schema)

	schemaMap(schema.Properties)
	schemaMap(schema.PatternProperties)
	for _, list := range [][]apiextensionsv1.JSONSchemaProps{
		schema.AllOf, schema.AnyOf, schema.OneOf,
	} {
		for i := range list {
			Schema(&list[i])
		}
	}
	Schema(schema.Not)
	schemaOrArray(schema.Items)
	schemaOrBool(schema.AdditionalProperties)
	schemaOrBool(schema.AdditionalItems)
}

func schemaMap(schemas map[string]apiextensionsv1.JSONSchemaProps) {
	for name := range schemas {
		child := schemas[name]
		Schema(&child)
		schemas[name] = child
	}
}

func schemaOrArray(node *apiextensionsv1.JSONSchemaPropsOrArray) {
	if node == nil {
		return
	}
	Schema(node.Schema)
	for i := range node.JSONSchemas {
		Schema(&node.JSONSchemas[i])
	}
}

func schemaOrBool(node *apiextensionsv1.JSONSchemaPropsOrBool) {
	if node == nil {
		return
	}
	Schema(node.Schema)
}

// denyUnknownFields spells out the strictness the apiserver gets from pruning,
// which leaves no trace in the schema a validator downloads.
func denyUnknownFields(schema *apiextensionsv1.JSONSchemaProps) {
	if len(schema.Properties) == 0 || schema.AdditionalProperties != nil {
		return
	}
	if schema.XPreserveUnknownFields != nil && *schema.XPreserveUnknownFields {
		return
	}
	schema.AdditionalProperties = &apiextensionsv1.JSONSchemaPropsOrBool{Allows: false}
}

func expandIntOrString(schema *apiextensionsv1.JSONSchemaProps) {
	if !schema.XIntOrString {
		return
	}
	schema.Type = ""
	schema.XIntOrString = false
	schema.OneOf = []apiextensionsv1.JSONSchemaProps{
		{Type: "string"},
		{Type: "integer"},
	}
}

func expandNullable(schema *apiextensionsv1.JSONSchemaProps) {
	if !schema.Nullable {
		return
	}
	schema.Nullable = false
	if schema.Type == "" || schema.Type == "null" {
		return
	}
	schema.OneOf = append(schema.OneOf,
		apiextensionsv1.JSONSchemaProps{Type: schema.Type},
		apiextensionsv1.JSONSchemaProps{Type: "null"},
	)
	schema.Type = ""
}

func dropExtensions(schema *apiextensionsv1.JSONSchemaProps) {
	schema.XPreserveUnknownFields = nil
	schema.XEmbeddedResource = false
	schema.XListMapKeys = nil
	schema.XListType = nil
	schema.XMapType = nil
	schema.XValidations = nil
}
