# Security Policy

## Supported versions

Only the latest minor release line is supported. Older tags receive fixes on a
best-effort basis if the fix is trivial; otherwise upgrade.

| Version | Supported          |
| ------- | ------------------ |
| `0.1.x` | :white_check_mark: |
| `< 0.1` | :x:                |

## Reporting a vulnerability

Report privately through
[GitHub Security Advisories](https://github.com/mrkhachaturov/crd2schema/security/advisories/new).
Please do not open a public issue for a vulnerability.

Expect an acknowledgement within a few days. Once a fix is released the advisory
is published with credit, unless you ask otherwise.

## Scope

`crd2schema` reads CustomResourceDefinitions and writes JSON files. It makes no
network calls and needs no cluster credentials of its own — the caller supplies
the CRDs, usually by piping `kubectl get crds -o yaml` into it.

The output is consumed by validators and editors, so a schema that is wrong in
the permissive direction is a security-relevant defect: it lets an invalid
manifest through review. Reports of that kind are in scope.
