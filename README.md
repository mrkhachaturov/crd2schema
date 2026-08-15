# crd2schema

Turns Kubernetes CustomResourceDefinitions into a JSON Schema draft-4 catalog,
laid out as `<group>/<kind>_<version>.json` — the layout
[kubeconform](https://github.com/yannh/kubeconform) and the editor's
yaml-language-server both expect.

```sh
kubectl get crds -o yaml | crd2schema -out schemas -
crd2schema -out schemas ./crds/*.yaml
```

## Why another converter

The usual path is datreeio's `crd-extractor.sh`, which downloads
`openapi2jsonschema.py` from kubeconform's default branch at run time. That
script walks the schema as untyped maps, so it cannot tell the JSON Schema
keyword `properties` from a CRD field that happens to be named `properties`.
When a CRD declares one — external-secrets does — it writes
`additionalProperties: false` inside the field list, where draft-4 requires a
schema. The file is then rejected by every validator, and kubeconform reports
that as *"could not find schema"*, which reads like the catalog is missing an
entry rather than serving a broken one.

`crd2schema` reads CRDs through the Kubernetes types themselves:

```go
Properties           map[string]JSONSchemaProps
AdditionalProperties *JSONSchemaPropsOrBool
```

The keyword is a struct field and the CRD's field is a map key, so the confusion
is not expressible. Every schema is compiled as draft-4 before it is written, so
a catalog that no validator can load is never produced.

## What it does to a schema

`JSONSchemaProps` is already draft-4 shaped. The conversion is limited to what
Kubernetes layers on top:

| Input | Output |
| --- | --- |
| `x-kubernetes-int-or-string` | `oneOf: [string, integer]` |
| `nullable: true` | `oneOf: [<type>, null]` |
| `x-kubernetes-*` | dropped |
| an object with `properties` | `additionalProperties: false` |

The last row is a deliberate addition. The apiserver rejects unknown fields by
pruning them against the structural schema, which leaves no trace in the schema
a validator downloads — so an object that lists its fields has to say the list
is exhaustive, or a typo passes review. Objects carrying
`x-kubernetes-preserve-unknown-fields` keep accepting anything, and a CRD that
declares its own `additionalProperties` keeps it.

## Install

```sh
mise use ubi:mrkhachaturov/crd2schema
```

or download a binary from the
[releases](https://github.com/mrkhachaturov/crd2schema/releases).

## Development

```sh
mise install
mise run gate     # fmt, vet, lint, test — what CI runs
```

`hk` gates commits and pushes; bypass a single command with `HK=0 git commit`.
