package transport

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-api-go/model"
)

// Call triggers the gRPC call function on space cloud
func (t *Transport) Call(ctx context.Context, token, project, service, endpoint string, params interface{}, timeout int) (*model.Response, error) {
	url := t.generateServiceURL(project, service, endpoint)
	// Fire the http request
	status, result, err := t.makeHTTPRequest(ctx, token, url, &model.ServiceRequest{Params: params, Timeout: timeout})
	if err != nil {
		return nil, err
	}

	if status >= 200 && status < 300 {
		return &model.Response{Status: status, Data: result}, nil
	}

	return &model.Response{Status: status, Error: result["error"].(string)}, nil
}

func (t *Transport) generateServiceURL(project, service, endpoint string) string {
	scheme := "http"
	if t.sslEnabled {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s/v1/api/%s/services/%s/%s", scheme, t.addr, project, service, endpoint)
}
