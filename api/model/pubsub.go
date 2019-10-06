package model

const PubsubUnSubscribre string = "pubsub-unsubscribe"
const PubsubSubscribe string = "pubsub-subscribe"
const PubsubSubscribeFeed string = "pubsub-subscribe-feed"

// PubsubMsg is the message recevied by the subscriber
type PubsubMsg struct {
	Subject string      `json:"subject"`
	Data    interface{} `json:"data"`
}

type PubsubCallback func(msg *PubsubMsg)

type PubsubUnsubscribe func() error

type PubsubSubscribeRequest struct {
	Subject string `json:"subject"`
	Queue   string `json:"queue"`
	Type    string `json:"type"`
	Token   string `json:"token"`
	Project string `json:"project"`
	Id      string `json:"id"`
}

type PubsubMsgResponse struct {
	Status int32  `json:"status"`
	Error  string `json:"error"`
	Msg    []byte `json:"msg"`
}

type PubsubPublishRequest struct {
	Subject string      `json:"subject"`
	Data    interface{} `json:"data"`
}

type TypeStore map[string]map[string]map[string]interface{}
type OnReceiveFunc func(subject string, data interface{})
type OnErrorFunc func(Error string)

type PubsubSubscribeEvents struct {
	OnReceive OnReceiveFunc
	OnError   OnErrorFunc
}
