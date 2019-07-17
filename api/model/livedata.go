package model

import (
	"encoding/json"
)

// LiveData is used to send the live data response to the user via onSnapshot
type LiveData struct {
	DataList []*Storage
}

// Unmarshal parses the response data and stores the result in the value pointed to by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
func (liveData *LiveData) Unmarshal(v interface{}) error {
	b := make([]byte, 0)
	if len(liveData.DataList) > 0 {
		b = append(b, []byte("[")...)
		for _, l := range liveData.DataList {
			if !l.IsDeleted {
				b = append(b, l.Payload...)
				b = append(b, []byte(",")...)
			}
		}
		b = b[:len(b)-1]
		b = append(b, []byte("]")...)
	} else {
		b = []byte("[]")
	}
	return json.Unmarshal(b, v)
}
