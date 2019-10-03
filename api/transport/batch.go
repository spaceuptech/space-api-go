package transport

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/utils"

	"github.com/spaceuptech/space-api-go/api/model"
)

// Batch triggers the gRPC batch function on space cloud
func (t *Transport) Batch(ctx context.Context, meta *model.Meta, r *model.BatchRequest) (*model.Response, error) {
	url := t.generateDatabaseURL(meta, utils.Batch)

	// Fire the http request
	status, result, err := t.makeHTTPRequest(meta.Token, url, r)
	if err != nil {
		return nil, err
	}

	if status >= 200 && status < 300 {
		return &model.Response{Status: status, Data: result}, nil
	}

	return &model.Response{Status: status, Error: result["error"].(string)}, nil
}
