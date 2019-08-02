package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// UnsubscribeFunction is the function sent to the user, to help him/her unsubscribe
type UnsubscribeFunction func()

// SnapshotFunction is the function that will be called when new realtime data is received
type SnapshotFunction func(*model.LiveData, string, *model.ChangedData)

// ErrorFunction is the function that will be called when an error occurred. Unsubscribe is automatically called
type ErrorFunction func(error)

// LiveQuery contains the methods for the liveQuery instance
type LiveQuery struct {
	config       *config.Config
	db           string
	col          string
	id           string
	find         utils.M
	store        []*model.Storage
	options      *liveQueryOptions
	subscription *LiveQuerySubscription
}

type liveQueryOptions struct {
	changesOnly bool
	skipInitial bool
}

func (l *LiveQuery) getOptions() []byte {
	// opts, err := json.Marshal("{\"skipInitial\":"+l.options.skipInitial+"}")
	opts, err := json.Marshal(l.options)
	if err != nil {
		log.Println("Could not marshal the options clause")
	}
	return opts
}

// Init returns a LiveQuery object
func Init(config *config.Config, db, col string) *LiveQuery {
	id := uuid.NewV1().String()
	return &LiveQuery{config, db, col, id, make(utils.M), make([]*model.Storage, 0), &liveQueryOptions{false, false}, nil}
}

// Where sets the where clause for the request
func (l *LiveQuery) Where(conds ...utils.M) *LiveQuery {
	if len(conds) == 1 {
		l.find = utils.GenerateFind(conds[0])
	} else {
		l.find = utils.GenerateFind(utils.And(conds...))
	}
	return l
}

func (l *LiveQuery) Options(options *model.LiveQueryOptions) *LiveQuery {
	l.options = &liveQueryOptions{options.ChangesOnly, options.ChangesOnly}
	return l
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

func (l *LiveQuery) snapshotCallback(feedData []*proto.FeedData, onSnapshot SnapshotFunction) {
	if l.options.changesOnly {
		for _, data := range feedData {
			if !(l.options.skipInitial && data.Type == utils.RealtimeInitial) {
				if data.Type != utils.RealtimeDelete {
					onSnapshot(&model.LiveData{nil}, data.Type, &model.ChangedData{data.Payload})
				} else {
					if l.db == utils.Mongo {
						onSnapshot(&model.LiveData{nil}, data.Type, &model.ChangedData{[]byte("{\"_id\":" + data.DocId + "}")})
					} else {
						onSnapshot(&model.LiveData{nil}, data.Type, &model.ChangedData{[]byte("{\"id\":" + data.DocId + "}")})
					}
				}
			}
		}
	} else {
		for _, data := range feedData {
			switch data.Type {
			case utils.RealtimeInitial:
				l.store = append(l.store, &model.Storage{data.DocId, data.TimeStamp, data.Payload, false})
			case utils.RealtimeInsert, utils.RealtimeUpdate:
				exists := false
				for _, s := range l.store {
					if s.Id == data.DocId {
						exists = true
						if s.Time <= data.TimeStamp {
							s.Time = data.TimeStamp
							s.Payload = data.Payload
							s.IsDeleted = false
						}
					}
				}
				if !exists {
					l.store = append(l.store, &model.Storage{data.DocId, data.TimeStamp, data.Payload, false})
				}
			case utils.RealtimeDelete:
				for _, s := range l.store {
					if s.Id == data.DocId && s.Time <= data.TimeStamp {
						s.Time = data.TimeStamp
						s.Payload = data.Payload
						s.IsDeleted = true
					}
				}
			}
		}
		if len(feedData) == 0 {
			liveData := &model.LiveData{l.store}
			l.subscription.snapshot = liveData
			onSnapshot(liveData, "initial", &model.ChangedData{make([]byte, 0)})
			return
		}
		changeType := feedData[0].Type
		if changeType == utils.RealtimeInitial {
			if !l.options.skipInitial {
				liveData := &model.LiveData{l.store}
				l.subscription.snapshot = liveData
				onSnapshot(liveData, changeType, &model.ChangedData{make([]byte, 0)})
			}
		} else {
			if !(changeType == utils.RealtimeDelete) {
				liveData := &model.LiveData{l.store}
				l.subscription.snapshot = liveData
				onSnapshot(liveData, changeType, &model.ChangedData{feedData[0].Payload})
			} else {
				if l.db == utils.Mongo {
					liveData := &model.LiveData{l.store}
					l.subscription.snapshot = liveData
					onSnapshot(liveData, changeType, &model.ChangedData{[]byte("{\"_id\":" + feedData[0].DocId + "}")})
				} else {
					liveData := &model.LiveData{l.store}
					l.subscription.snapshot = liveData
					onSnapshot(liveData, changeType, &model.ChangedData{[]byte("{\"id\":" + feedData[0].DocId + "}")})
				}
			}
		}
	}
}

// Subscribe is used to subscribe to a particular live query
func (l *LiveQuery) Subscribe(onSnapshot SnapshotFunction, onError ErrorFunction) *LiveQuerySubscription {
	conn := l.config.Transport.GetConn()
	var stream proto.SpaceCloud_RealTimeClient
	go func() {
		for {
			state := conn.GetState()
			if state.String() == "READY" {
				var err error
				stream, err = l.config.Transport.GetStub().RealTime(context.TODO())
				if err != nil {
					continue
				}
				findJSON, err := json.Marshal(l.find)
				if err != nil {
					log.Println("Could not marshal the where clause")
				}
				subscribeRequest := &proto.RealTimeRequest{Token: l.config.Token, DbType: l.db, Project: l.config.Project, Group: l.col, Type: utils.TypeRealtimeSubscribe, Id: l.id, Where: findJSON, Options: l.getOptions()}
				if err := stream.Send(subscribeRequest); err != nil {
					log.Println("Failed to send the subscribe request")
				}
				for {
					in, err := stream.Recv()
					if err == io.EOF {
						return
					}
					if err != nil {
						log.Println("Failed to receive a request:", err)
						break
					}

					if in.Id == l.id {
						if in.Ack {
							l.snapshotCallback(in.FeedData, onSnapshot)
						} else {
							onError(errors.New(in.Error))
							l.unsubscribe(stream)()
						}
					}
				}
			} else {
				// log.Println("Not connected. Attempting to Reconnect...")
				conn.WaitForStateChange(context.TODO(), state)
			}
		}
	}()
	l.subscription = &LiveQuerySubscription{l.unsubscribe(stream), &model.LiveData{nil}}
	return l.subscription
}
