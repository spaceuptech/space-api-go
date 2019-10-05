package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"log"
	"sync"
	"time"
)

const ServiceRegister string = "service-register"
const ServiceRequest string = "service-request"

type websocketOptions struct {
	projectId string
	token     string
}

type Socket struct {
	url                  string
	isConnected          bool
	connectedOnce        bool
	options              websocketOptions
	pendingMsg           []model.WebsocketMessage
	Socket               *websocket.Conn
	SendMessage          chan model.WebsocketMessage
	registerCallbackMap  map[string]func(data interface{})
	onReconnectCallbacks []func()
	mux                  sync.RWMutex
}

func Init(url string, config *config.Config) *Socket {
	url = "ws://" + url + "/v1/api/" + config.Project + "/socket/json"
	if config.IsSecure {
		url = "wss://" + url + "/v1/api/" + config.Project + "/socket/json"
	}
	return &Socket{
		url:                 url,
		options:             websocketOptions{projectId: config.Project, token: config.Token},
		registerCallbackMap: map[string]func(data interface{}){},
		pendingMsg:          []model.WebsocketMessage{},
		mux:                 sync.RWMutex{},
	}
}

func (s *Socket) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(s.url, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	s.setSocket(conn)
	s.setConnected(true)

	if s.getConnectedOnce() {
		for _,value := range s.onReconnectCallbacks {
			value()
		}
	}

	if len(s.pendingMsg) != 0 {
		for _, payload := range s.pendingMsg {
			s.Socket.WriteJSON(payload)
		}
		s.pendingMsg = []model.WebsocketMessage{}
	}

	writeMessage := make(chan model.WebsocketMessage)
	defer close(writeMessage)
	s.setWriterChannel(writeMessage)

	// create a websocket reader & writer
	go s.writerRoutine()
	go s.read()

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

func (s Socket) writerRoutine() {
	for msg := range s.SendMessage {

		if !s.getConnected() {
			s.mux.Lock()
			s.pendingMsg = append(s.pendingMsg, msg)
			s.mux.Unlock()
		}

		if err := s.Socket.WriteJSON(msg); err != nil {
			log.Println(err)
		}
	}
}

func (s *Socket) Send(Type string, data interface{}) string {
	id := uuid.NewV1().String()
	s.SendMessage <- model.WebsocketMessage{ID: id, Type: Type, Data: data}
	return id
}

func (s *Socket) read() {
	msg := model.WebsocketMessage{}
	for {
		if s.getConnected() {
			if err := s.Socket.ReadJSON(&msg); err != nil {
				log.Println(err)
				s.setConnected(false)
				continue
			}
		} else {
			if err := s.connect(); err != nil {
				log.Println(err)
				continue
			}
		}

		if msg.ID != "" {
			cb, ok := s.registerCallbackMap[msg.ID]
			if ok {
				go cb(msg.Data)
				s.unregisterCallback(msg.ID)
			}
		}

		cb, ok := s.registerCallbackMap[msg.Type]
		if ok {
			go cb(msg.Data)
		}
	}
}
