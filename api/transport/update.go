package transport

import (
	"context"
	"encoding/json"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Update triggers the gRPC update function on space cloud
func (t *Transport) Update(ctx context.Context, meta *proto.Meta, op string, find, update utils.M) (*model.Response, error) {
	updateJSON, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}

	findJSON, err := json.Marshal(find)
	if err != nil {
		return nil, err
	}

	req := proto.UpdateRequest{Find: findJSON, Update: updateJSON, Meta: meta, Operation: op}
	res, err := t.stub.Update(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}
