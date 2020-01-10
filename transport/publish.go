package transport

import (
	"context"
	"encoding/json"

	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/proto"
)

// PubsubPublish triggers the gRPC PubsubPublish function on space cloud
func (t *Transport) PubsubPublish(ctx context.Context, meta *proto.Meta, subject string, msg interface{}) (*model.Response, error) {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req := proto.PubsubPublishRequest{Subject: subject, Msg: msgJSON, Meta: meta}
	res, err := t.stub.PubsubPublish(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: nil}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}
