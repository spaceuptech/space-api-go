package model

import (
	"encoding/json"
)

// ChangedData is used to send the changed data response to the user via onSnapshot
type ChangedData struct {
	Data []byte
}

// Unmarshal parses the response data and stores the result in the value pointed to by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
func (changedData *ChangedData) Unmarshal(v interface{}) error {
	if len(changedData.Data) == 0 {
		return nil
	}
	return json.Unmarshal(changedData.Data, v)
}
