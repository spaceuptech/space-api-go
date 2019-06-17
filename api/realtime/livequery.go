package realtime

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log"
	"io"
	"errors"
	"encoding/json"

	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// UnsubscribeFunction is the function sent to the user, to help him/her unsubscribe
type UnsubscribeFunction func()()

// SnapshotFunction is the function that will be called when new realtime data is received
type SnapshotFunction func(*model.LiveData, string)()

// ErrorFunction is the function that will be called when an error occurred. Unsubscribe is automatically called
type ErrorFunction func(error)()

// LiveQuery contains the methods for the liveQuery instance
type LiveQuery struct {
	config  *config.Config
	db      string
	col     string
	id      string
	find    utils.M
	store   []*model.Storage
}

// Init returns a LiveQuery object
func Init(config *config.Config, db, col string) *LiveQuery {
	id := uuid.NewV1().String()
	return &LiveQuery{config, db, col, id, make(utils.M), make([]*model.Storage, 0)}
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

func (l *LiveQuery) unsubscribe(stream proto.SpaceCloud_RealTimeClient) UnsubscribeFunction {
	return func() {
		stream.Send(&proto.RealTimeRequest{Token: l.config.Token, DbType: l.db, Project: l.config.Project, Group: l.col, Type: utils.TypeRealtimeUnsubscribe, Id: l.id})
		stream.CloseSend()
	}
}

func (l *LiveQuery) snapshotCallback(feedData []*proto.FeedData, onSnapshot SnapshotFunction) {
	if len(feedData) > 0 {
		for _, data := range feedData {
			switch data.Type {
			case utils.RealtimeInsert, utils.RealtimeUpdate:
				exists := false
				for _, s:= range l.store {
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
				for _, s:= range l.store {
					if s.Id == data.DocId && s.Time <= data.TimeStamp {
						s.Time = data.TimeStamp
						s.Payload = data.Payload
						s.IsDeleted = true
					}
				}
			}
		}
		if len(feedData) == 1 {
			onSnapshot(&model.LiveData{l.store}, feedData[0].Type)
		} else {
			onSnapshot(&model.LiveData{l.store}, "initial")
		}
	}
}

// Subscribe is used to subscribe to a particular live query
func (l *LiveQuery) Subscribe(onSnapshot SnapshotFunction, onError ErrorFunction) UnsubscribeFunction {
	conn := l.config.Transport.GetConn()
	var stream proto.SpaceCloud_RealTimeClient
	go func() {
		for {
			state := conn.GetState()
			if state.String() == "READY" {
				// log.Println("Connected to Space Cloud")
				var err error
				stream, err = l.config.Transport.GetStub().RealTime(context.TODO())
				if err != nil {
					continue
				}
				findJSON, err := json.Marshal(l.find)
				if err != nil {
					log.Println("Could not marshal the where clause")
				}
				subscribeRequest := &proto.RealTimeRequest{Token: l.config.Token, DbType: l.db, Project: l.config.Project, Group: l.col, Type: utils.TypeRealtimeSubscribe, Id: l.id, Where: findJSON}
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
	return l.unsubscribe(stream)
}