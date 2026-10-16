# go-cerebras

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-cerebras.svg)](https://pkg.go.dev/github.com/psyb0t/go-cerebras)
[![CI](https://github.com/psyb0t/go-cerebras/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-cerebras/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-cerebras/badges/version.svg)](https://github.com/psyb0t/go-cerebras/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-cerebras/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-cerebras/badges/importers.svg)](https://github.com/psyb0t/go-cerebras/blob/badges/importers.md)

Unofficial Go client for the [Cerebras Inference API](https://inference-docs.cerebras.ai), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Cerebras publishes. I am not affiliated with Cerebras.

The package is `cerebras`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-cerebras
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/psyb0t/go-cerebras"
)

func main() {
	auth := cerebras.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+os.Getenv("CEREBRAS_API_KEY"))

		return nil
	})

	client, err := cerebras.NewClientWithResponses("https://api.cerebras.ai", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.ListModelsV1ModelsGetWithResponse(
		context.Background(),
		&cerebras.ListModelsV1ModelsGetParams{},
	)
	if err != nil {
		panic(err)
	}

	list, err := res.JSON200.AsModelMetadataList()
	if err != nil {
		panic(err)
	}

	for _, m := range list.Data {
		fmt.Println(m.Id)
	}
}
```

`example_test.go` runs the same code against a local test server, so `make test` checks it.

Method names come from the spec's operation IDs, so they are long: `ChatV1ChatCompletionsPost`, `CompletionsV1CompletionsPost`, `ListModelsV1ModelsGet`, `GetModelV1ModelsModelIdGet`.

Many schemas in the spec are `anyOf` unions, and oapi-codegen turns each into a struct with `As...` and `From...` methods. The model list response is one (`AsModelMetadataList`, or the OpenRouter and Hugging Face formats), the chat response is another (`AsChatCompletionResponse`), and so is every chat message. A user message is built like this:

```go
var content cerebras.UserMessageRequest_Content
_ = content.FromUserMessageRequestContent0("hello")

role := cerebras.UserMessageRequestRoleUser

var msg cerebras.ChatCompletionRequest_Messages_Item
_ = msg.FromUserMessageRequest(cerebras.UserMessageRequest{Content: content, Role: &role})

res, err := client.ChatV1ChatCompletionsPostWithResponse(ctx,
	&cerebras.ChatV1ChatCompletionsPostParams{},
	cerebras.ChatCompletionRequest{Model: "<model-id>", Messages: []cerebras.ChatCompletionRequest_Messages_Item{msg}},
)
```

## Auth

Cerebras uses a Bearer token: send the API key as `Authorization: Bearer <key>`. Create a key in the Cerebras Cloud console. The request editor in the usage example sets the header on every request. The spec doesn't list a server, so pass the base URL, `https://api.cerebras.ai`, yourself.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` reads `openapi_spec_url` from the `.stats.yml` of [Cerebras/cerebras-cloud-sdk-python](https://github.com/Cerebras/cerebras-cloud-sdk-python) and downloads that URL to `api/openapi.yaml`. Stainless hosts the spec under a URL with a content hash that changes on every update, so the script never hardcodes it.
- oapi-codegen runs with `config.yaml`. The config maps `number` to `float64`, because oapi-codegen would otherwise use `float32`.

## Agent integrations

The [skill](.agents/skills/go-cerebras) teaches a coding agent how to use go-cerebras when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-cerebras@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-cerebras@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-cerebras
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Cerebras's published API specification. The specification itself is not redistributed here.
