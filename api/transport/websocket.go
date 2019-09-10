package transport

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type websocketConnection struct {
	connected            bool
	isConnecting         bool
	connectedOnce        bool
	callBacks            map[string]CallBackFunction
	onReconnectCallbacks []CallBackFunction
	url                  string
	pendingRequest       []Payload
	address              string
	sslEnable            bool
	conn                 *websocket.Conn
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

func (w *websocketConnection) RegisterCallBackFuntion(Type string, function CallBackFunction) {
	w.callBacks[Type] = function
}

func (w *websocketConnection) UnRegisterCallBackFuntion(Type string) {
	delete(w.callBacks, Type)
}

func (w *websocketConnection) RegisterOnReconnectCallback(function CallBackFunction) {
	w.onReconnectCallbacks = append(w.onReconnectCallbacks, function)
}

func (w *websocketConnection) Init() {
	w.isConnecting = false
	w.connectedOnce = false
	w.connected = false
}

func (w *websocketConnection) Connect(addr string, sslEnabled bool) error {
	w.address = addr
	w.sslEnable = sslEnabled
	w.isConnecting = true

	var u url.URL
	if sslEnabled {
		u = url.URL{Scheme: "wss", Host: addr, Path: "/v1/api/socket/json"}
	} else {
		u = url.URL{Scheme: "ws", Host: addr, Path: "/v1/api/socket/json"}
	}
	w.url = u.String()
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial:", err)
		w.connected = false
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

func (w *websocketConnection) Send(Typ string, data interface{}) error {
	payload := Payload{Type: Typ, Data: data}
	if !w.connected {
		w.pendingRequest = append(w.pendingRequest, payload)
		if !w.isConnecting {
			err := w.Connect(w.address, w.sslEnable)
			if err != nil {
				log.Println("Error establishing websocket", err)
				w.connected = false
				return err
			}
			w.connected = true
		}
	}
	w.conn.WriteJSON(payload)
	return nil
}
