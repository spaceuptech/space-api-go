package transport

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
)

// Batch triggers the gRPC batch function on space cloud
func (t *Transport) Batch(ctx context.Context, meta *proto.Meta, requests []*proto.AllRequest) (*model.Response, error) {
	req := proto.BatchRequest{Meta: meta, Batchrequest: requests}
	res, err := t.stub.Batch(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}
