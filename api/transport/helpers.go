package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/utils"
)

const contentTypeJSON string = "application/json"
const post string = "POST"

func (t *Transport) generateDatabaseURL(meta *model.Meta, op string) string {
	scheme := "http"
	if t.sslEnabled {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s/v1/api/%s/crud/%s/%s/%s", scheme, t.addr, meta.Project, meta.DB, meta.Col, op)
}

func (t *Transport) makeHTTPRequest(token, url string, payload interface{}) (int, utils.M, error) {
	// Marshal the payoad
	data, err := json.Marshal(payload)
	if err != nil {
		return -1, nil, err
	}

	// Make a http request
	r, err := http.NewRequest(post, url, bytes.NewBuffer(data))
	if err != nil {
		return -1, nil, err
	}

	// Add appropriate headers
	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("Content-Type", "application/json")

	// Fire the request
	res, err := t.httpClient.Do(r)
	if err != nil {
		return -1, nil, err
	}
	defer res.Body.Close()

	// Unmarshal the response
	result := utils.M{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return -1, nil, err
	}

	// Return the final response
	return res.StatusCode, result, nil
}
