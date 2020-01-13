package realtime

import (
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/transport/websocket"
	"github.com/spaceuptech/space-api-go/utils"
)

type Realtime struct {
	appId  string
	store  model.DbStore
	client *websocket.Socket
}

// Init initialize the realtime module
func Init(appID string, client *websocket.Socket) *Realtime {
	r := &Realtime{appId: appID, client: client, store: make(model.DbStore, 0)}

	// on reconnect register again according the value in store
	r.client.RegisterOnReconnectCallback(func() {
		for db, dbValue := range r.store {
			for col, colValue := range dbValue {
				for id := range colValue {
					obj := r.store[db][col][id]
					q := r.LiveQuery(db, col)
					q.options = obj.Options.(*model.LiveQueryOptions)
					q.params = &model.RealtimeParams{Find: obj.Find.(utils.M)}
					q.subscribeRaw(id)
				}
			}
		}
	})

	// initialize the realtime on sc
	r.client.RegisterCallback(utils.TypeRealtimeFeed, func(data interface{}) {
		var feedData FeedData
		if err := mapstructure.Decode(data, &feedData); err != nil {
			log.Fatal("error while decoding map structure in realtime:", err)
		}
		snapshotCallback(r.store, []FeedData{feedData})
	})
	return r
}

// LiveQuery initialize the live query module
func (r *Realtime) LiveQuery(db, collection string) *LiveQuery {
	return New(r.appId, db, collection, r.client, r.store)
}
