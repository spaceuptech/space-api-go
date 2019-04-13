package transport

import (
	"context"
	"encoding/json"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Read triggers the gRPC read function on space cloud
func (t *Transport) Read(ctx context.Context, meta *proto.Meta, find utils.M, op string, options *proto.ReadOptions) (*model.Response, error) {
	findJSON, err := json.Marshal(find)
	if err != nil {
		return nil, err
	}

	req := proto.ReadRequest{Find: findJSON, Meta: meta, Operation: op, Options: options}
	res, err := t.stub.Read(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 || res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}
