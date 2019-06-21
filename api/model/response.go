package model

import (
	"encoding/json"
	"errors"
)

// Response is the object recieved from the server. Status is the status code received from the server.
// Data is either a map[string]interface{} or
type Response struct {
	Status int
	Data   []byte
	Error  string
}

// Unmarshal parses the response data and stores the result in the value pointed to by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
func (res *Response) Unmarshal(v interface{}) error {
	if res.Status < 200 || res.Status >= 300 {
		return errors.New("Result not present")
	}
	return json.Unmarshal(res.Data, v)
}
