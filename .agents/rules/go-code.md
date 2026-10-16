# Go code

Refines `~/.claude/rule-details/working-with-generated-code.md` and `~/.claude/rule-details/coding-style.md`.

- **GO1** `client.gen.go` is generated. Never edit it by hand. Change the spec patch or `config.yaml`, then run `make generate`.
- **GO2** The spec is downloaded into the gitignored `api/` directory and never committed. A spec problem oapi-codegen cannot handle gets fixed by a patch step under `scripts/` (jq for JSON, a Go tool under `internal/cmd` for YAML) that `gen.go` runs before oapi-codegen.
- **GO3** Hand-written Go is limited to `gen.go` and the tests. The package never logs and never prints.
- **GO4** `make lint` passes with zero issues. Lines are at most 80 columns, the `.golangci.yml` limit; split calls one argument per line rather than adding `//nolint`.
