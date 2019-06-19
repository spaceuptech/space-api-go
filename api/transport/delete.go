package transport

import (
	"context"
	"encoding/json"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Delete triggers the gRPC delete function on space cloud
func (t *Transport) Delete(ctx context.Context, meta *proto.Meta, op string, find utils.M) (*model.Response, error) {
	findJSON, err := json.Marshal(find)
	if err != nil {
		return nil, err
	}

	req := proto.DeleteRequest{Find: findJSON, Meta: meta, Operation: op}
	res, err := t.stub.Delete(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}
