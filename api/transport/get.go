package transport

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Read triggers the gRPC read function on space cloud
func (t *Transport) Read(ctx context.Context, meta *model.Meta, r *model.ReadRequest) (*model.Response, error) {
	url := t.generateDatabaseURL(meta, utils.Read)

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
