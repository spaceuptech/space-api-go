package realtime

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/spaceuptech/space-api-go/api/transport/websocket"
	"log"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/utils"
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
	store   model.TypeStore
	options *model.LiveQueryOptions
	params  *model.RealtimeParams
}

// Init returns a LiveQuery object
func Init(appId, db, col string, client *websocket.Socket, store []*model.Store) *LiveQuery {
	return &LiveQuery{appId: appId, db: db, col: col, client: client, store: model.TypeStore{}, options: &model.LiveQueryOptions{}, params: &model.RealtimeParams{}}
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

func (l *LiveQuery) addSubscription(id string, onSnapShot model.SnapshotFunction, onError model.OnErrorFunction) func() {
	l.store[l.db][l.col][id] = model.Store{Subscription: model.StoreSubscription{OnSnapShot: onSnapShot, OnError: onError}}
	return func() {
		_, err := l.client.Request(model.RealtimeUnsubscribe, model.RealtimeRequest{Group: l.col, ID: id, Options: l.options})
		if err != nil {

		}
		delete(l.store[l.db][l.col], id)
	}
}

func (l *LiveQuery) subscribe(onSnapShot model.SnapshotFunction, onError model.OnErrorFunction) *model.StoreSubscriptionObject {
	id := uuid.NewV1().String()
	return l.subscribeRaw(id, onSnapShot, onError)
}

func (l *LiveQuery) unsubscribe(stream proto.SpaceCloud_RealTimeClient) UnsubscribeFunction {
	return func() {
		if stream == nil {
			return
		}
		stream.Send(&proto.RealTimeRequest{Token: l.config.Token, DbType: l.db, Project: l.config.Project, Group: l.col, Type: utils.TypeRealtimeUnsubscribe, Id: l.id, Options: l.getOptions()})
		stream.CloseSend()
	}
}

// Subscribe is used to subscribe to a particular live query
func (l *LiveQuery) subscribeRaw(id string, onSnapShot model.SnapshotFunction, onError model.OnErrorFunction) *model.StoreSubscriptionObject {
	req := model.RealtimeRequest{DBType: l.db, Project: l.appId, Group: l.col, ID: id, Where: l.params.Find, Options: l.options}

	_, ok := l.store[l.db]
	if !ok {

	}
	_, ok = l.store[l.db][l.col]
	if !ok {

	}
	l.store[l.db][l.col][id] = model.Store{Snapshot: []model.SnapshotData{}, Subscription: model.StoreSubscription{}, Find: l.params.Find, Options: l.options}

	unsubscribe := l.addSubscription(id, onSnapShot, onError)
}
