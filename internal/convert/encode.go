package convert

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Encode serialises a converted schema and rejects one no validator could use.
// kubeconform reports a schema that fails to compile as a missing schema, so an
// unusable file would otherwise reach the catalog and read as an absent entry.
func Encode(schema *apiextensionsv1.JSONSchemaProps) ([]byte, error) {
	encoded, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}
	encoded = append(encoded, '\n')

	if err := compiles(encoded); err != nil {
		return nil, err
	}
	return encoded, nil
}

func compiles(encoded []byte) error {
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(encoded))
	if err != nil {
		return fmt.Errorf("reparse: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft4)
	if err := compiler.AddResource("schema.json", doc); err != nil {
		return fmt.Errorf("add resource: %w", err)
	}
	if _, err := compiler.Compile("schema.json"); err != nil {
		return fmt.Errorf("does not compile as draft-4: %w", err)
	}
	return nil
}
