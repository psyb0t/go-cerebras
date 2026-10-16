package cerebras_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/psyb0t/go-cerebras"
)

// This is the README usage example. A local test server stands in for
// https://api.cerebras.ai so the example runs without credentials.
func Example() {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(
				`{"object":"list","data":[{"id":"model-x"}]}`,
			))
		},
	))
	defer srv.Close()

	baseURL := srv.URL // https://api.cerebras.ai in real use

	auth := cerebras.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"Bearer "+os.Getenv("CEREBRAS_API_KEY"),
			)

			return nil
		},
	)

	client, err := cerebras.NewClientWithResponses(baseURL, auth)
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

	// Output: model-x
}
