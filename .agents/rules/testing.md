# Testing

Refines `~/.claude/rule-details/go/writing-tests.md`.

- **TST1** Tests are black-box, in package `cerebras_test`, and use only the exported API.
- **TST2** Every test runs the client against an `httptest` server and checks method, path, the `Authorization: Bearer` header and the decoded response. No test talks to the real API.
- **TST3** Union request and response bodies are built and read through the generated `From...` and `As...` methods, the way callers use them.
- **TST4** `example_test.go` holds the README usage example and runs against a local server with an `// Output:` line. Change both together.
- **TST5** `make test` runs `go test -race`.
