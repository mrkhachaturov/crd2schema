// Command crd2schema turns CustomResourceDefinitions into a JSON Schema draft-4
// catalog laid out as <group>/<kind>_<version>.json.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"

	"github.com/mrkhachaturov/crd2schema/internal/convert"
)

func main() {
	out := flag.String("out", "schemas", "directory to write the catalog into")
	flag.Parse()

	written, err := run(*out, flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "crd2schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("wrote %d schemas to %s\n", written, *out)
}

func run(out string, paths []string) (int, error) {
	if len(paths) == 0 {
		return 0, errors.New("no input files given")
	}
	if err := os.MkdirAll(out, 0o750); err != nil {
		return 0, fmt.Errorf("create output directory: %w", err)
	}

	written := 0
	for _, path := range paths {
		crds, err := read(path)
		if err != nil {
			return written, err
		}
		for _, crd := range crds {
			n, err := write(out, crd)
			if err != nil {
				return written, fmt.Errorf("%s: %w", path, err)
			}
			written += n
		}
	}
	return written, nil
}

func read(path string) ([]*apiextensionsv1.CustomResourceDefinition, error) {
	var raw []byte
	var err error
	if path == "-" {
		raw, err = io.ReadAll(os.Stdin)
	} else {
		raw, err = os.ReadFile(path) //nolint:gosec // the caller chooses the path
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var crds []*apiextensionsv1.CustomResourceDefinition
	for _, doc := range strings.Split(string(raw), "\n---") {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		var list apiextensionsv1.CustomResourceDefinitionList
		if err := yaml.Unmarshal([]byte(doc), &list); err == nil && len(list.Items) > 0 {
			for i := range list.Items {
				crds = append(crds, &list.Items[i])
			}
			continue
		}

		var crd apiextensionsv1.CustomResourceDefinition
		if err := yaml.Unmarshal([]byte(doc), &crd); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		if crd.Kind == "CustomResourceDefinition" {
			crds = append(crds, &crd)
		}
	}
	return crds, nil
}

func write(out string, crd *apiextensionsv1.CustomResourceDefinition) (int, error) {
	written := 0
	for i := range crd.Spec.Versions {
		version := &crd.Spec.Versions[i]
		if version.Schema == nil || version.Schema.OpenAPIV3Schema == nil {
			continue
		}

		schema := version.Schema.OpenAPIV3Schema.DeepCopy()
		convert.Schema(schema)

		encoded, err := convert.Encode(schema)
		if err != nil {
			return written, fmt.Errorf("%s/%s: %w", crd.Spec.Names.Kind, version.Name, err)
		}

		dir := filepath.Join(out, crd.Spec.Group)
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return written, fmt.Errorf("create %s: %w", dir, err)
		}
		name := fmt.Sprintf("%s_%s.json", strings.ToLower(crd.Spec.Names.Kind), version.Name)
		if err := os.WriteFile(filepath.Join(dir, name), encoded, 0o600); err != nil {
			return written, fmt.Errorf("write %s: %w", name, err)
		}
		written++
	}
	return written, nil
}
