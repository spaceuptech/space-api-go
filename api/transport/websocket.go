package transport

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type WebsocketConnection struct {
	connected            bool
	isConnecting         bool
	connectedOnce        bool
	callBacks            map[string]CallBackFunction
	onReconnectCallbacks []CallBackFunction
	url                  string
	pendingRequest       []Payload
	sslEnable            bool
	conn                 *websocket.Conn
	options              Payload
}

type Payload struct {
	Auth     []byte      `json:"auth"`
	Token    string      `json:"token"`
	Project  string      `json:"project"`
	Type     string      `json:"type"`
	Id       string      `json:"id"`
	Service  string      `json:"service"`
	Params   []byte      `json:"params"`
	Function string      `json:"function"`
	Error    string      `json:"error"`
	Data     interface{} `json:"data"`
}

func (w *WebsocketConnection) RegisterCallBackFuntion(Type string, function CallBackFunction) {
	w.callBacks[Type] = function
}

func (w *WebsocketConnection) UnRegisterCallBackFuntion(Type string) {
	delete(w.callBacks, Type)
}

func (w *WebsocketConnection) RegisterOnReconnectCallback(function CallBackFunction) {
	w.onReconnectCallbacks = append(w.onReconnectCallbacks, function)
}

func (w *WebsocketConnection) Init(addr string, sslEnabled bool) {
	w.isConnecting = false
	w.connectedOnce = false
	w.connected = false
	w.sslEnable = sslEnabled

	var u url.URL
	if sslEnabled {
		u = url.URL{Scheme: "wss", Host: addr, Path: "/v1/api/socket/json"}
	} else {
		u = url.URL{Scheme: "ws", Host: addr, Path: "/v1/api/socket/json"}
	}
	w.url = u.String()
}

func (w *WebsocketConnection) Connect() error {
	w.isConnecting = true

	c, _, err := websocket.DefaultDialer.Dial(w.url, nil)
	if err != nil {
		log.Println("dial:", err)
		w.isConnecting = false
		return err
	}

	w.conn = c
	w.connected = true
	w.connectedOnce = true

	// clear the pending requests
	if w.connected {
		for _, v := range w.pendingRequest {
			err := c.WriteJSON(v)
			if err != nil {
				log.Println("Error in sending file ", err)
			}
		}
		w.pendingRequest = w.pendingRequest[:0] // clear all pending requests
	}
	return nil
}

func (w *WebsocketConnection) Send(Typ string, data Payload) (string, error) {
	id := uuid.NewV1().String() // genrate id using satori
	data.Token = w.options.Token
	payload := Payload{Id: id, Type: Typ, Data: data}
	if !w.connected {
		w.pendingRequest = append(w.pendingRequest, payload)
		if !w.isConnecting {
			err := w.Connect()
			if err != nil {
				log.Println("Error establishing websocket", err)
				w.connected = false
				return id, err
			}
			w.connected = true
		}
	}
	w.conn.WriteJSON(payload)
	return id, nil
}

func (w *WebsocketConnection) Request(Type string, data Payload) {
	_, err := w.Send(Type, data)
	if err != nil {
		log.Println("Error in sending : ", err)
	}

	//	w.RegisterCallBackFuntion(id)
}
