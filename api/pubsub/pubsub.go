package pubsub

import (
	"bytes"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/transport/websocket"
	"log"
	"net/http"
)

type Pubsub struct {
	url    string
	appId  string
	client *websocket.Socket
	store  model.TypeStore
}

func Init(url, appId string, client *websocket.Socket) *Pubsub {
	return &Pubsub{url: url, appId: appId, client: client, store: model.TypeStore{}}
}

func (p *Pubsub) queueSubscribe(subject, queue string, onReceive model.OnReceiveFunc, onError model.OnErrorFunc) *pubsubSubscription {
	s := PubsubSubscribeInit(p.appId, subject, queue, p.client, p.store)
	return s.subscribe(onReceive, onError)
}

func (p *Pubsub) subscribe(subject string, onReceive model.OnReceiveFunc, onError model.OnErrorFunc) *pubsubSubscription {
	return p.queueSubscribe(subject, "", onReceive, onError)
}

func (p *Pubsub) publish(subject string, data interface{}) (string, map[string]interface{}) {
	p.client.RegisterOnReconnectCallback(func() {
		for subject, subjectValue := range p.store {
			for queue, queueValue := range subjectValue {
				for range queueValue {
					p.queueSubscribe(subject, queue, nil, nil)
				}
			}
		}
	})

	p.client.RegisterCallback(model.PubsubSubscribeFeed, func(data interface{}) {
		m := &model.WebsocketMessage{}
		if err := mapstructure.Decode(data, m); err != nil {
			log.Println("pubsub publishing structure decoding error", err)
			return
		}
		for _, subjectValue := range p.store {
			for _, queueValue := range subjectValue {
				for id, value := range queueValue {
					if id == m.ID {
						onReceive, ok := value.(model.OnReceiveFunc)
						if ok {
							msg, ok := m.Data.(model.PubsubPublishRequest)
							if !ok {

							}
							onReceive(msg.Subject, msg.Data)
						}
					}
				}
			}
		}
	})

	requestBody, err := json.Marshal(map[string]interface{}{
		"subject": subject,
		"data":    data,
	})
	if err != nil {

	}
	res, err := http.Post(p.url+"v1/api"+p.appId+"/pubsub", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {

	}
	defer res.Body.Close()
	result := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(result); err != nil {

	}
	return res.Status, result

}
