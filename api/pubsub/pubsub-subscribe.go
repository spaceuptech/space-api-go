package pubsub

import (
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/transport/websocket"
	"log"
)

type pubsubSubscribe struct {
	appId   string
	subject string
	queue   string
	client  *websocket.Socket
	store   model.TypeStore
}

func PubsubSubscribeInit(appId, subject string, queue string, client *websocket.Socket, store model.TypeStore) *pubsubSubscribe {
	return &pubsubSubscribe{appId: appId, subject: subject, queue: queue, client: client, store: store}
}

func (p *pubsubSubscribe) addSubscription(id string, onReceive model.OnReceiveFunc, onError model.OnErrorFunc) func() {
	p.store[p.subject][p.queue][id] = model.PubsubSubscribeEvents{OnReceive: onReceive, OnError: onError}
	return func() {
		_, err := p.client.Request(model.PubsubUnSubscribre, model.PubsubSubscribeRequest{Subject: p.subject, Id: id})
		if err != nil {
			log.Println("pubsub add subscription", err)
		}
		delete(p.store[p.subject][p.queue], id)
	}
}

func (p *pubsubSubscribe) subscribe(onReceive model.OnReceiveFunc, onError model.OnErrorFunc) *pubsubSubscription {
	id := uuid.NewV1().String()
	return p.subscribeRaw(id, onReceive, onError)
}

func (p *pubsubSubscribe) subscribeRaw(id string, onReceive model.OnReceiveFunc, onError model.OnErrorFunc) *pubsubSubscription {
	req := &model.PubsubSubscribeRequest{Id: id, Subject: p.subject, Queue: p.queue, Project: p.appId}

	_, ok := p.store[p.subject]
	if !ok {
		p.store[p.subject] = map[string]map[string]interface{}{}
	}

	_, ok = p.store[p.subject][p.queue]
	if !ok {
		p.store[p.subject][p.queue] = map[string]interface{}{}
	}

	p.store[p.subject][p.queue][id] = model.PubsubSubscribeEvents{} // check this

	unsubscribe := p.addSubscription(id, onReceive, onError)

	data, err := p.client.Request(model.PubsubSubscribe, req)
	if err != nil {
		onError(err.Error())
		unsubscribe()
	}

	res := model.PubsubMsgResponse{}
	if err := mapstructure.Decode(data, res); err != nil {
		log.Println("pubsub subscribe raw error", err)
		return nil
	}
	if res.Status != 200 {
		onError(res.Error)
		unsubscribe()
	}

	p.client.RegisterCallback(model.PubsubSubscribeFeed, func(data interface{}) {
		m := &model.PubsubPublishRequest{}
		if err := mapstructure.Decode(data, m); err != nil {
			log.Println("pubsub subscribe raw error", err)
		}
		onReceive(m.Subject, m.Data)
	})
	return pubsubSubscriptionInit(p.subject, unsubscribe)
}
