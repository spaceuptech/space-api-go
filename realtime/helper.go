package realtime

import (
	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/utils"
)

func (l *LiveQuery) snapshotCallback(store model.TypeStore, rows []FeedData) {
	if len(rows) == 0 {
		return
	}
	var obj = model.Store{}
	var opts = model.LiveQueryOptions{}
	for _, data := range rows {
		obj = store[data.DBType][data.Group][data.QueryID]
		opts = obj.QueryOptions
		if opts.ChangesOnly {
			if !(opts.SkipInitial && data.Type == utils.RealtimeInitial) {
				if data.Type != utils.RealtimeDelete {
					obj.Subscription.OnSnapShot(nil, data.Type, data.Payload)
				} else {
					if l.db == utils.Mongo {
						obj.Subscription.OnSnapShot(nil, data.Type, map[string]string{"_id": data.DocID})
					} else {
						obj.Subscription.OnSnapShot(nil, data.Type, map[string]string{"id": data.DocID})
					}
				}
			}
		} else {
			if data.Type == utils.RealtimeInitial {
				obj.Snapshot = append(obj.Snapshot, model.SnapshotData{Id: data.DocID, Time: data.TimeStamp, Payload: data.Payload, IsDeleted: false})
			} else if data.Type == utils.RealtimeInsert || data.Type == utils.RealtimeUpdate {
				isExisting := false
				temp := []interface{}{}
				for _, row := range obj.Snapshot {
					if row.Id == data.DocID {
						isExisting = true
						if row.Time <= data.TimeStamp {
							q := model.SnapshotData{Id: row.Id, Time: data.TimeStamp, Payload: data.Payload, IsDeleted: false}
							temp = append(temp, q)
						}
						temp = append(temp, row)
					}
				}
				if !isExisting {
					obj.Snapshot = append(obj.Snapshot, model.SnapshotData{Id: data.DocID, Time: data.TimeStamp, Payload: data.Payload, IsDeleted: false,})
				}
			} else if data.Type == utils.RealtimeDelete {
				temp := []interface{}{}
				for _, row := range obj.Snapshot {
					if row.Id == data.DocID && row.Time <= data.TimeStamp {
						q := model.SnapshotData{Id: row.Id, Time: data.TimeStamp, Payload: data.Payload, IsDeleted: false}
						temp = append(temp, q)
					}
					temp = append(temp, row)
				}
			}
		}
	}
	if !opts.ChangesOnly {
		changeType := rows[0].Type
		if changeType == utils.RealtimeInitial {
			if !opts.SkipInitial {
				temp := []model.SnapshotData{}
				for _, row := range obj.Snapshot {
					if !row.IsDeleted {
						a := model.SnapshotData{Payload: row.Payload}
						temp = append(temp, a)
					}
				}
				obj.SubscriptionObject.SnapShot = temp
				obj.Subscription.OnSnapShot(obj.SubscriptionObject.SnapShot, changeType, nil)
			}
		} else {
			if changeType != utils.RealtimeDelete {
				temp := []model.SnapshotData{}
				for _, row := range obj.Snapshot {
					if !row.IsDeleted {
						a := model.SnapshotData{Payload: row.Payload}
						temp = append(temp, a)
					}
				}
				obj.SubscriptionObject.SnapShot = temp
				obj.Subscription.OnSnapShot(obj.SubscriptionObject.SnapShot, changeType, nil)
			}
		}
	}

}

// FeedData is the format to send realtime data
type FeedData struct {
	QueryID   string      `json:"id" structs:"id"`
	DocID     string      `json:"docId" structs:"docId"`
	Type      string      `json:"type" structs:"type"`
	Payload   interface{} `json:"payload" structs:"payload"`
	TimeStamp int64       `json:"time" structs:"time"`
	Group     string      `json:"group" structs:"group"`
	DBType    string      `json:"dbType" structs:"dbType"`
	TypeName  string      `json:"__typename,omitempty" structs:"__typename,omitempty"`
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
