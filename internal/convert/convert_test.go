package convert

import (
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func object(props map[string]apiextensionsv1.JSONSchemaProps) *apiextensionsv1.JSONSchemaProps {
	return &apiextensionsv1.JSONSchemaProps{Type: "object", Properties: props}
}

// Converters that walk untyped maps mistake a CRD field named "properties" for
// the keyword and write the strictness flag into the field list, where draft-4
// demands a schema.
func TestFieldNamedPropertiesStaysAField(t *testing.T) {
	schema := object(map[string]apiextensionsv1.JSONSchemaProps{
		"properties": {Type: "array", Items: &apiextensionsv1.JSONSchemaPropsOrArray{
			Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
		}},
	})

	Schema(schema)

	field, ok := schema.Properties["properties"]
	if !ok {
		t.Fatal("the field named properties was dropped")
	}
	if field.Type != "array" {
		t.Errorf("field type = %q, want array", field.Type)
	}
	if schema.AdditionalProperties == nil || schema.AdditionalProperties.Allows {
		t.Error("strictness was not applied to the enclosing object")
	}
	if _, err := Encode(schema); err != nil {
		t.Errorf("Encode: %v", err)
	}
}

func TestUnknownFieldsDenied(t *testing.T) {
	schema := object(map[string]apiextensionsv1.JSONSchemaProps{
		"name": {Type: "string"},
	})

	Schema(schema)

	if schema.AdditionalProperties == nil {
		t.Fatal("additionalProperties was not set")
	}
	if schema.AdditionalProperties.Allows {
		t.Error("additionalProperties allows unknown fields")
	}
}

func TestPreserveUnknownFieldsIsHonoured(t *testing.T) {
	preserve := true
	schema := object(map[string]apiextensionsv1.JSONSchemaProps{
		"name": {Type: "string"},
	})
	schema.XPreserveUnknownFields = &preserve

	Schema(schema)

	if schema.AdditionalProperties != nil {
		t.Error("strictness was applied to a schema that opts out of pruning")
	}
	if schema.XPreserveUnknownFields != nil {
		t.Error("the x-kubernetes extension survived into the output")
	}
}

func TestDeclaredAdditionalPropertiesIsKept(t *testing.T) {
	schema := object(map[string]apiextensionsv1.JSONSchemaProps{
		"name": {Type: "string"},
	})
	schema.AdditionalProperties = &apiextensionsv1.JSONSchemaPropsOrBool{
		Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
	}

	Schema(schema)

	if schema.AdditionalProperties.Schema == nil {
		t.Fatal("the declared additionalProperties schema was replaced")
	}
	if schema.AdditionalProperties.Schema.Type != "string" {
		t.Errorf("type = %q, want string", schema.AdditionalProperties.Schema.Type)
	}
}

func TestIntOrStringBecomesOneOf(t *testing.T) {
	schema := object(map[string]apiextensionsv1.JSONSchemaProps{
		"port": {XIntOrString: true},
	})

	Schema(schema)

	port := schema.Properties["port"]
	if len(port.OneOf) != 2 {
		t.Fatalf("oneOf has %d entries, want 2", len(port.OneOf))
	}
	if port.XIntOrString {
		t.Error("the x-kubernetes marker survived into the output")
	}
}

func TestNullableBecomesTypeUnion(t *testing.T) {
	schema := object(map[string]apiextensionsv1.JSONSchemaProps{
		"note": {Type: "string", Nullable: true},
	})

	Schema(schema)

	note := schema.Properties["note"]
	if note.Nullable {
		t.Error("nullable survived, and draft-4 does not know it")
	}
	if len(note.OneOf) != 2 {
		t.Fatalf("oneOf has %d entries, want 2", len(note.OneOf))
	}
}

func TestNestedObjectsAreReached(t *testing.T) {
	schema := object(map[string]apiextensionsv1.JSONSchemaProps{
		"spec": *object(map[string]apiextensionsv1.JSONSchemaProps{
			"replicas": {Type: "integer"},
		}),
	})

	Schema(schema)

	spec := schema.Properties["spec"]
	if spec.AdditionalProperties == nil || spec.AdditionalProperties.Allows {
		t.Error("a nested object was left permissive")
	}
}

func TestArrayElementsAreReached(t *testing.T) {
	schema := object(map[string]apiextensionsv1.JSONSchemaProps{
		"rules": {
			Type: "array",
			Items: &apiextensionsv1.JSONSchemaPropsOrArray{
				Schema: object(map[string]apiextensionsv1.JSONSchemaProps{
					"name": {Type: "string"},
				}),
			},
		},
	})

	Schema(schema)

	item := schema.Properties["rules"].Items.Schema
	if item.AdditionalProperties == nil || item.AdditionalProperties.Allows {
		t.Error("an array element schema was left permissive")
	}
}
