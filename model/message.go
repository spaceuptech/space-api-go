package model

import (
	"encoding/json"
)

type Message struct {
	Data   []byte
}

// Unmarshal parses the response data and stores the result in the value pointed to by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
func (msg *Message) Unmarshal(v interface{}) error {
	return json.Unmarshal(msg.Data, v)
}
