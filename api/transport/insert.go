package transport

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/model"

	"github.com/spaceuptech/space-api-go/api/utils"
)

// Insert triggers the gRPC create function on space cloud
func (t *Transport) Insert(ctx context.Context, meta *model.Meta, i *model.CreateRequest) (*model.Response, error) {
	// Make url for request
	url := t.generateDatabaseURL(meta, utils.Create)

	// Fire the http request
	status, result, err := t.makeHTTPRequest(meta.Token, url, i)
	if err != nil {
		return nil, err
	}

	if status >= 200 && status < 300 {
		return &model.Response{Status: status, Data: result}, nil
	}

	return &model.Response{Status: status, Error: result["error"].(string)}, nil
}
