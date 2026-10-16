package cerebras_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-cerebras"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAPIKey = "test-key"

func newClient(
	t *testing.T,
	handler http.HandlerFunc,
) *cerebras.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := cerebras.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testAPIKey)

			return nil
		},
	)

	client, err := cerebras.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func assertBearer(t *testing.T, r *http.Request) {
	t.Helper()

	assert.Equal(t, "Bearer "+testAPIKey, r.Header.Get("Authorization"))
}

func TestListModels(t *testing.T) {
	t.Parallel()

	client := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/models", r.URL.Path)
		assertBearer(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"object":"list","data":[{"id":"model-x",` +
				`"object":"model","owned_by":"owner-x"}]}`,
		))
	})

	res, err := client.ListModelsV1ModelsGetWithResponse(
		context.Background(),
		&cerebras.ListModelsV1ModelsGetParams{},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)

	list, err := res.JSON200.AsModelMetadataList()
	require.NoError(t, err)
	require.Len(t, list.Data, 1)

	assert.Equal(t, "model-x", list.Data[0].Id)
	assert.Equal(t, "owner-x", *list.Data[0].OwnedBy)
}

func TestChatCompletion(t *testing.T) {
	t.Parallel()

	client := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assertBearer(t, r)

		var body map[string]any

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "model-x", body["model"])
		assert.Equal(t, []any{
			map[string]any{"role": "user", "content": "hello"},
		}, body["messages"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"id":"c1","object":"chat.completion","created":1,` +
				`"model":"model-x","system_fingerprint":"fp",` +
				`"choices":[{"index":0,"finish_reason":"stop",` +
				`"message":{"role":"assistant","content":"hi"}}]}`,
		))
	})

	var content cerebras.UserMessageRequest_Content

	require.NoError(t, content.FromUserMessageRequestContent0("hello"))

	role := cerebras.UserMessageRequestRoleUser

	var msg cerebras.ChatCompletionRequest_Messages_Item

	require.NoError(t, msg.FromUserMessageRequest(cerebras.UserMessageRequest{
		Content: content,
		Role:    &role,
	}))

	res, err := client.ChatV1ChatCompletionsPostWithResponse(
		context.Background(),
		&cerebras.ChatV1ChatCompletionsPostParams{},
		cerebras.ChatCompletionRequest{
			Model:    "model-x",
			Messages: []cerebras.ChatCompletionRequest_Messages_Item{msg},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)

	completion, err := res.JSON200.AsChatCompletionResponse()
	require.NoError(t, err)
	require.Len(t, completion.Choices, 1)

	assert.Equal(t, "c1", completion.Id)
	assert.Equal(t, "hi", *completion.Choices[0].Message.Content)
}
