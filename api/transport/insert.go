package transport

import (
	"context"
	"encoding/json"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
)

// Insert triggers the gRPC create function on space cloud
func Insert(ctx context.Context, stub proto.SpaceCloudClient, meta *proto.Meta, op string, obj interface{}) (*model.Response, error) {
	objJSON, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	req := proto.CreateRequest{Document: objJSON, Meta: meta, Operation: op}
	res, err := stub.Create(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 || res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}
