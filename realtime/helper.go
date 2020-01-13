package realtime

import (
	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/utils"
)

func snapshotCallback(store model.DbStore, rows []FeedData) {
	if len(rows) == 0 {
		return
	}
	var obj = new(model.Store)
	var opts = model.LiveQueryOptions{}
	for _, data := range rows {
		obj = store[data.DBType][data.Group][data.QueryID]
		opts = obj.QueryOptions
		if opts.ChangesOnly {
			if !(opts.SkipInitial && data.Type == utils.RealtimeInitial) {
				if data.Type != utils.RealtimeDelete {
					obj.C <- model.NewSubscriptionEvent(data.Type, data.Payload, data.Find, nil)
				} else {
					obj.C <- model.NewSubscriptionEvent(data.Type, nil, data.Find, nil)
				}
			}
		} else {
			if data.Type == utils.RealtimeInitial {
				obj.Snapshot = append(obj.Snapshot, &model.SnapshotData{Find: data.Find, Time: data.TimeStamp, Payload: data.Payload, IsDeleted: false})
				obj.C <- model.NewSubscriptionEvent(data.Type, data.Payload, data.Find, nil)
			} else if data.Type == utils.RealtimeInsert || data.Type == utils.RealtimeUpdate {
				isExisting := false
				for _, row := range obj.Snapshot {
					if validate(data.Find, row.Payload.(map[string]interface{})) {
						isExisting = true
						if row.Time <= data.TimeStamp {
							row.Time = data.TimeStamp
							row.Payload = data.Payload
							row.IsDeleted = false
							obj.C <- model.NewSubscriptionEvent(data.Type, data.Payload, data.Find, nil)
						}
					}
				}
				if !isExisting {
					obj.Snapshot = append(obj.Snapshot, &model.SnapshotData{Find: data.Find, Time: data.TimeStamp, Payload: data.Payload, IsDeleted: false})
					obj.C <- model.NewSubscriptionEvent(data.Type, data.Payload, data.Find, nil)
				}
			} else if data.Type == utils.RealtimeDelete {
				for _, row := range obj.Snapshot {
					if validate(row.Find, row.Payload.(map[string]interface{})) && row.Time <= data.TimeStamp {
						row.Time = data.TimeStamp
						row.Payload = map[string]interface{}{}
						row.IsDeleted = true
						obj.C <- model.NewSubscriptionEvent(data.Type, nil, data.Find, nil)
					}
				}
			}
		}
	}
}

func validate(find map[string]interface{}, doc map[string]interface{}) bool {
	for k, v := range find {
		keyValue, p := doc[k]
		if !p {
			return false
		}

		if keyValue != v {
			return false
		}
	}
	return true
}

// FeedData is the format to send realtime data
type FeedData struct {
	QueryID   string                 `json:"id" mapstructure:"id" structs:"id"`
	Find      map[string]interface{} `json:"find" structs:"find"`
	Type      string                 `json:"type" structs:"type"`
	Payload   interface{}            `json:"payload" structs:"payload"`
	TimeStamp int64                  `json:"time" structs:"time"`
	Group     string                 `json:"group" structs:"group"`
	DBType    string                 `json:"dbType" structs:"dbType"`
	TypeName  string                 `json:"__typename,omitempty" structs:"__typename,omitempty"`
}

// RealtimeResponse is the object sent for realtime requests
type RealtimeResponse struct {
	Group string      `json:"group"` // Group is the collection name
	ID    string      `json:"id"`    // id is the query id
	Ack   bool        `json:"ack"`
	Error string      `json:"error"`
	Docs  []*FeedData `json:"docs"`
}

// Message is the request body of the message
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	ID   string      `json:"id"` // the request id
}
