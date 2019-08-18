package pubsub

import (
	"context"
	"encoding/json"
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

// OnReceive is the function that will be called when new data is published to a subscribed subject
type OnReceive func(string, interface{})

// Pubsub contains the methods for the pubsub instance
type Pubsub struct {
	config       *config.Config
	id           string
	subscription *PubsubSubscription
	onReceive    OnReceive
}

// Init returns a Pubsub object
func Init(config *config.Config) *Pubsub {
	id := uuid.NewV1().String()
	return &Pubsub{config, id, nil, nil}
}

func (p *Pubsub) Unsubscribe(stream proto.SpaceCloud_PubsubSubscribeClient, subject string) UnsubscribeFunction {
	return func() {
		if stream == nil {
			return
		}
		stream.Send(&proto.PubsubSubscribeRequest{Subject: subject, Type: utils.TypePubsubUnsubscribe, Token: p.config.Token, Project: p.config.Project, Id: p.id})
		stream.CloseSend()
	}
}

// Subscribe is used to subscribe to a particular subject and its children
func (p *Pubsub) Subscribe(subject string, onReceive OnReceive) *PubsubSubscription {
	return p.QueueSubscribe(subject, "", onReceive)
}

// QueueSubscribe is used to subscribe to a particular subject and its children, using a queue
func (p *Pubsub) QueueSubscribe(subject, queue string, onReceive OnReceive) *PubsubSubscription {
	p.onReceive = onReceive
	conn := p.config.Transport.GetConn()
	var stream proto.SpaceCloud_PubsubSubscribeClient
	go func() {
		for {
			state := conn.GetState()
			if state.String() == "READY" {
				var err error
				stream, err = p.config.Transport.GetStub().PubsubSubscribe(context.TODO())
				if err != nil {
					continue
				}
				subscribeRequest := &proto.PubsubSubscribeRequest{Token: p.config.Token, Project: p.config.Project, Type: utils.TypePubsubSubscribe, Id: p.id, Subject: subject, Queue: queue}
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

					if in.Id == p.id {
						if in.Type == utils.TypePubsubSubscribeFeed {
							var m map[string]interface{}
							err := json.Unmarshal(in.Msg, &m)
							if err != nil {
								log.Println("Error decoding received message")
							}
							p.onReceive(m["subject"].(string), m["data"])
						} else if in.Status != 200 {
							log.Println("Pubsub Error:", "OperationType=", in.Type, "Status=", in.Status, in.Error)
                        	p.Unsubscribe(stream, subject)
                        	return
						}
					}
				}
			} else {
				// log.Println("Not connected. Attempting to Reconnect...")
				conn.WaitForStateChange(context.TODO(), state)
			}
		}
	}()
	p.subscription = &PubsubSubscription{subject, p.Unsubscribe(stream, subject)}
	return p.subscription
}

// Publish publishes a message to a particular subject
func (p *Pubsub) Publish(subject string, msg interface{}) (*model.Response, error) {
	return p.config.Transport.PubsubPublish(context.TODO(), &proto.Meta{Project: p.config.Project, Token: p.config.Token}, subject, msg)
}