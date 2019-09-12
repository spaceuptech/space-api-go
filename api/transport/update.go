package transport

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Update triggers the gRPC update function on space cloud
func (t *Transport) Update(ctx context.Context, meta *model.Meta, d *model.UpdateRequest) (*model.Response, error) {
	url := t.generateDatabaseURL(meta, utils.Update)

	// Fire the http request
	status, result, err := t.makeHTTPRequest(meta.Token, url, utils.M{"find": d.Find, "op": d.Operation, "update": d.Update})
	if err != nil {
		return nil, err
	}

	if status >= 200 && status < 300 {
		return &model.Response{Status: status, Data: result}, nil
	}

	return &model.Response{Status: status, Error: result["error"].(string)}, nil
}
