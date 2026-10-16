---
name: go-cerebras
description: Unofficial Go client for the Cerebras Inference API, generated with oapi-codegen. `cerebras.NewClientWithResponses("https://api.cerebras.ai", cerebras.WithRequestEditorFn(...))` from `github.com/psyb0t/go-cerebras` gives typed methods for all 7 paths of the spec, chat completions (`ChatV1ChatCompletionsPost`), text completions (`CompletionsV1CompletionsPost`), model listing and lookup, and TCP warming. Union schemas come with `As...`/`From...` helpers. Use when Go code needs to call Cerebras inference, or already imports go-cerebras.
homepage: https://github.com/psyb0t/go-cerebras
user-invocable: true
metadata:
  openclaw:
    emoji: "🧠"
    requires:
      bins:
        - go
      env:
        - CEREBRAS_API_KEY
permissions:
  filesystem: "Guides edits to the Go source of the current project. The library reads and writes no files."
  shell: "The Go toolchain (`go get`, `go build`, `go test`) to add and compile the dependency."
  network: "`go get` fetches the module from your Go module proxy. At runtime the client sends HTTPS requests only to the base URL you give it, normally https://api.cerebras.ai."
---

# go-cerebras

Typed Go client for the Cerebras Inference API. Every operation in the spec is a method on `ClientWithResponses`, named after its operation ID.

For install, the method list, unions and error handling in detail, see [references/setup.md](references/setup.md).

## Security & safety

- Read the API key from an environment variable such as `CEREBRAS_API_KEY`. Never hardcode it, commit it or log it, and never print request headers.
- Requests go only to the base URL passed to `NewClientWithResponses`. Prompts sent through the client reach Cerebras.

## When to use

- Go code calls Cerebras for chat or text completions, or lists the models it serves.
- The project already imports `github.com/psyb0t/go-cerebras`.

## When NOT to use

- You need streaming. The generated client returns the whole body; for `stream: true` read `HTTPResponse.Body` from the plain `Client` methods yourself.
- The code is not Go. Use Cerebras's own SDKs.

## Usage

```go
import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/psyb0t/go-cerebras"
)

auth := cerebras.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+os.Getenv("CEREBRAS_API_KEY"))

	return nil
})

client, err := cerebras.NewClientWithResponses("https://api.cerebras.ai", auth)
if err != nil {
	return err
}

var content cerebras.UserMessageRequest_Content
if err := content.FromUserMessageRequestContent0("Say hi"); err != nil {
	return err
}

role := cerebras.UserMessageRequestRoleUser

var msg cerebras.ChatCompletionRequest_Messages_Item
if err := msg.FromUserMessageRequest(cerebras.UserMessageRequest{Content: content, Role: &role}); err != nil {
	return err
}

res, err := client.ChatV1ChatCompletionsPostWithResponse(ctx,
	&cerebras.ChatV1ChatCompletionsPostParams{},
	cerebras.ChatCompletionRequest{Model: "<model-id>", Messages: []cerebras.ChatCompletionRequest_Messages_Item{msg}},
)
if err != nil {
	return err // transport error
}

if res.JSON200 == nil {
	return fmt.Errorf("chat: status %d: %s", res.StatusCode(), res.Body)
}

completion, err := res.JSON200.AsChatCompletionResponse()
if err != nil {
	return err
}

fmt.Println(*completion.Choices[0].Message.Content)
```

Listing models:

```go
res, err := client.ListModelsV1ModelsGetWithResponse(ctx, &cerebras.ListModelsV1ModelsGetParams{})
if err != nil {
	return err
}

list, err := res.JSON200.AsModelMetadataList()
```

Check `JSON200` for nil before reading it. A non-2xx status is not a Go error; the status is in `StatusCode()` and the raw body in `Body`.
