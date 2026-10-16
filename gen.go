// Package cerebras is a Go client for the Cerebras Inference API, generated
// with oapi-codegen from the OpenAPI document Cerebras publishes.
package cerebras

//go:generate bash scripts/fetch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.yaml
