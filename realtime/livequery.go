package realtime

import (
	"encoding/json"
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-api-go/transport/websocket"

	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/utils"
)

// UnsubscribeFunction is the function sent to the user, to help him/her unsubscribe
type UnsubscribeFunction func()

// SnapshotFunction is the function that will be called when new realtime data is received

// ErrorFunction is the function that will be called when an error occurred. Unsubscribe is automatically called
type ErrorFunction func(error)

// LiveQuery contains the methods for the liveQuery instance
type LiveQuery struct {
	appId   string
	db      string
	col     string
	client  *websocket.Socket
	store   model.DbStore
	options *model.LiveQueryOptions
	params  *model.RealtimeParams
}

// New returns a LiveQuery object
func New(appId, db, col string, client *websocket.Socket, store model.DbStore) *LiveQuery {
	return &LiveQuery{appId: appId, db: db, col: col, client: client, store: store, options: &model.LiveQueryOptions{}, params: &model.RealtimeParams{}}
}

// Where sets the where clause for the request
func (l *LiveQuery) Where(conds ...utils.M) *LiveQuery {
	if len(conds) == 1 {
		l.params.Find = utils.GenerateFind(conds[0])
	} else {
		l.params.Find = utils.GenerateFind(utils.And(conds...))
	}
	return l
}

func (l *LiveQuery) getOptions() []byte {
	// opts, err := json.Marshal("{\"skipInitial\":"+l.options.skipInitial+"}")
	opts, err := json.Marshal(l.options)
	if err != nil {
		log.Println("Could not marshal the options clause")
	}
	return opts
}

func (l *LiveQuery) Options(options *model.LiveQueryOptions) *LiveQuery {
	l.options = &model.LiveQueryOptions{ChangesOnly: options.ChangesOnly, SkipInitial: options.ChangesOnly}
	return l
}

// used internally
func (l *LiveQuery) addSubscription(id string) func() {
	return func() {
		_, err := l.client.Request(model.RealtimeUnsubscribe, model.RealtimeRequest{Group: l.col, ID: id, Options: l.options})
		if err != nil {
			log.Println("Failed to unsubscribe", err)
		}
		delete(l.store[l.db][l.col], id)
	}
}

// Subscribe is used to subscribe to a new document
func (l *LiveQuery) Subscribe() *model.StoreSubscriptionObject {
	id := uuid.NewV1().String()
	return l.subscribeRaw(id)
}

// subscribeRaw is used to subscribe to a particular live query
func (l *LiveQuery) subscribeRaw(id string) *model.StoreSubscriptionObject {
	req := model.RealtimeRequest{DBType: l.db, Project: l.appId, Group: l.col, ID: id, Where: l.params.Find, Options: l.options}

	_, ok := l.store[l.db]
	if !ok {
		l.store[l.db] = model.ColStore{}
	}
	_, ok = l.store[l.db][l.col]
	if !ok {
		l.store[l.db][l.col] = model.IdStore{}
	}
	c := make(chan *model.SubscriptionEvent, 5)
	v := &model.Store{Snapshot: []*model.SnapshotData{}, C: c, Find: l.params.Find, Options: l.options}
	l.store[l.db][l.col][id] = v

	unsubscribe := l.addSubscription(id)

	data, err := l.client.Request(utils.TypeRealtimeSubscribe, req)
	if err != nil {
		c <- model.NewSubscriptionEvent("", nil, nil, fmt.Errorf("error unable to subscribe to realtime feature %v", err))
		unsubscribe()
	}

	go func() {
		mapData, ok := data.(map[string]interface{})
		if ok {
			ack, ok := mapData["ack"]
			if ok && !ack.(bool) {
				err, ok := mapData["error"].(string)
				if ok {
					c <- model.NewSubscriptionEvent("", nil, nil, fmt.Errorf("error from server %v", err))
					unsubscribe()
				}
			}
			docs := mapData["docs"].([]interface{})
			for _, doc := range docs {
				c <- model.NewSubscriptionEvent("initial", doc.(map[string]interface{})["payload"], l.params.Find, nil)
			}
		}
	}()

	v.SubscriptionObject = model.LiveQuerySubscriptionInit(unsubscribe, []model.SnapshotData{}, c)

	return &v.SubscriptionObject
}
