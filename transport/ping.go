package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-api-go/types"
	"github.com/spaceuptech/space-api-go/utils"
)

// Call triggers the gRPC call function on space cloud
func (t *Transport) Ping(ctx context.Context, token string) error {
	scheme := "http"
	if t.sslEnabled {
		scheme = "https"
	}

	// Make a http request
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s://%s/v1/config/env", scheme, t.addr), nil)
	if err != nil {
		return err
	}

	// Add appropriate headers
	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("Content-Type", contentTypeJSON)

	// Fire the request
	res, err := t.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer utils.CloseTheCloser(res.Body)

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}

	// Unmarshal the response
	result := types.M{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}

	return fmt.Errorf("Service responde with status code (%v) with error message - (%v) ", res.StatusCode, result["error"].(string))
}
