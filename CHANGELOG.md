# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-16

- Added `cerebras`, a REST client generated with oapi-codegen from Cerebras's OpenAPI document (Cerebras Cloud API, spec version 0.2.0, OpenAPI 3.1.0). It covers all 7 paths: chat completions, completions, the authenticated and public model listings, and TCP warming.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` reads the current spec URL from the `.stats.yml` of Cerebras's Python SDK repo, downloads the spec from Stainless, and regenerates the client. The spec is not committed.
- Added the `go-cerebras` agent skill: `claude plugin install go-cerebras@psyb0t`, `codex plugin add go-cerebras@psyb0t`, `openclaw skills install @psyb0t/go-cerebras`.
