package utils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"renebizelli/go/observability/SystemA/dtos"
)

func ExecRequestWithContext[T any](ctx context.Context, URL string, headers map[string]string) (*T, *dtos.RequestError) {

	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return nil, &dtos.RequestError{Message: err.Error(), StatusCode: 400}
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err := http.DefaultClient.Do(req)

	select {
	case <-ctx.Done():
		return nil, &dtos.RequestError{Message: err.Error(), StatusCode: 499}
	default:

		if err != nil {
			return nil, &dtos.RequestError{Message: err.Error(), StatusCode: 500}
		}

		if res.StatusCode != 200 {
			return nil, &dtos.RequestError{Message: "Invalid request", StatusCode: res.StatusCode}
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, &dtos.RequestError{Message: err.Error(), StatusCode: 500}
		}

		var o *T
		err = json.Unmarshal(body, &o)

		if err != nil {
			return nil, &dtos.RequestError{Message: err.Error(), StatusCode: 500}
		}

		return o, nil
	}

}
