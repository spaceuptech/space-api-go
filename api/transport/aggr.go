package transport

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Aggr triggers the gRPC aggr function on space cloud
func (t *Transport) Aggr(ctx context.Context, meta *model.Meta, a *model.AggregateRequest) (*model.Response, error) {

	url := t.generateDatabaseURL(meta, utils.Aggregate)

	status, result, err := t.makeHTTPRequest(meta.Token, url, a)
	if err != nil {
		return nil, err
	}

	if status >= 200 && status < 300 {
		return &model.Response{Status: status, Data: result}, nil
	}

	return &model.Response{Status: status, Error: result["error"].(string)}, nil
}
