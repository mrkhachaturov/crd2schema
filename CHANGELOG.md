# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Convert CustomResourceDefinitions to a JSON Schema draft-4 catalog laid out as
  `<group>/<kind>_<version>.json`.
- Reject unknown fields by writing `additionalProperties: false` for every object
  that lists its properties, matching how the apiserver prunes against a
  structural schema. Objects carrying `x-kubernetes-preserve-unknown-fields`, and
  CRDs that declare their own `additionalProperties`, are left alone.
- Translate `x-kubernetes-int-or-string` and `nullable` into their draft-4
  spellings, and drop the remaining `x-kubernetes-*` vocabulary.
- Compile every schema as draft-4 before writing it, so a catalog no validator
  can load is never produced.
- Accept both a `CustomResourceDefinitionList` and a stream of documents, so
  `kubectl get crds -o yaml` can be piped in directly.
