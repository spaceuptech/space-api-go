package websocket

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/model"
)

type websocketOptions struct {
	projectId string
	token     string
}

type Socket struct {
	url                  string
	isConnect            bool
	isConnecting         bool
	connectedOnce        bool
	options              websocketOptions
	pendingMsg           []model.WebsocketMessage
	socket               *websocket.Conn
	sendMessage          chan model.WebsocketMessage
	registerCallbackMap  map[string]func(data interface{})
	onReconnectCallbacks []func()
	mux                  sync.RWMutex
}

func Init(url string, config *config.Config) *Socket {
	url = "ws://" + url + "/v1/api/" + config.Project + "/socket/json"
	if config.IsSecure {
		url = "wss://" + url + "/v1/api/" + config.Project + "/socket/json"
	}

	s := &Socket{
		url:                 url,
		options:             websocketOptions{projectId: config.Project, token: config.Token},
		registerCallbackMap: map[string]func(data interface{}){},
		pendingMsg:          []model.WebsocketMessage{},
		mux:                 sync.RWMutex{},
	}

	writeMessage := make(chan model.WebsocketMessage)
	s.setWriterChannel(writeMessage)

	// create a websocket reader & writer
	go s.read()
	go s.writerRoutine()

	return s
}

func (s *Socket) connect() error {
	if !s.checkIsConnecting() {
		return nil
	}
	conn, _, err := websocket.DefaultDialer.Dial(s.url, nil)
	if err != nil {
		log.Println("websocket dialer error", err)
		s.resetIsConnecting()
		return err
	}

	s.resetIsConnecting()
	s.setSocket(conn)
	s.setConnected(true)

	if s.isConnectedOnce() {
		for _, fn := range s.onReconnectCallbacks {
			fn()
		}
	}

	s.setConnectedOnce(true)

	s.mux.Lock()
	if len(s.pendingMsg) > 0 {
		for _, payload := range s.pendingMsg {
			if err := s.socket.WriteJSON(payload); err != nil {
				log.Println("error writing pending messages into websocket", err)
			}
		}
		s.pendingMsg = []model.WebsocketMessage{}
	}
	s.mux.Unlock()

	return nil
}

func (s *Socket) Request(msgType string, data interface{}) (interface{}, error) {
	if !s.getConnected() {
		// connect to server
		if err := s.connect(); err != nil {
			return false, err
		}
	}

	id := s.Send(msgType, data)

	timer1 := time.NewTimer(10 * time.Second)
	defer timer1.Stop()

	// channel for receiving service register acknowledgement
	ch := make(chan interface{})
	defer close(ch)

	s.RegisterCallback(id, func(data interface{}) {
		ch <- data
	})

	select {
	case <-timer1.C:
		return false, errors.New("response time elapsed")
	case msg := <-ch:
		return msg, nil
	}
}

func (s *Socket) writerRoutine() {
	for msg := range s.sendMessage {
		if !s.getConnected() {
			s.mux.Lock()
			s.pendingMsg = append(s.pendingMsg, msg)
			s.mux.Unlock()
			continue
		}

		if err := s.socket.WriteJSON(msg); err != nil {
			log.Println("error writing into websocket", err)
		}
	}
}

// Send sends a message to server over websocket protocol
func (s *Socket) Send(Type string, data interface{}) string {
	id := uuid.NewV1().String()
	s.sendMessage <- model.WebsocketMessage{ID: id, Type: Type, Data: data}
	return id
}

func (s *Socket) read() {
	for {
		msg := &model.WebsocketMessage{}
		if s.getConnected() {
			if err := s.socket.ReadJSON(msg); err != nil {
				log.Println("error reading from websocket", err)
				s.setConnected(false)
				time.Sleep(5 * time.Second)
				continue
			}
		} else {
			if err := s.connect(); err != nil {
				log.Println(err)
				time.Sleep(5 * time.Second)
				continue
			}
		}

		if msg != nil {
			if msg.ID != "" {
				cb, ok := s.getRegisteredCallBack(msg.ID)
				if ok {
					log.Println(msg.ID)
					go cb(msg.Data)
					s.unregisterCallback(msg.ID)
				}
			}

			cb, ok := s.getRegisteredCallBack(msg.Type)
			if ok {
				go cb(msg.Data)
			}
		}
	}
}
