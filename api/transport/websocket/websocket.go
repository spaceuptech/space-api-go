package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/spaceuptech/space-api-go/api/config"
	"log"
	"sync"
	"time"
)

const ServiceRegister string = "service-register"
const ServiceRequest string = "service-request"

type Socket struct {
	url             string
	options         websocketOptions
	isConnected     bool
	channels        map[string]chan FunctionsPayload
	serviceRegister map[string]chan map[string]bool
	WriteMessage    chan WriteMessageStructure
	Socket          *websocket.Conn
	pendingMsg      []WriteMessageStructure
	mux             sync.RWMutex
}

func (s *Socket) getConnected() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.isConnected
}

func (s *Socket) setConnected(value bool) {
	s.mux.Lock()
	s.isConnected = value
	s.mux.Unlock()
}

func (s *Socket) RegisterChannel(name string, channel chan FunctionsPayload) {
	_, ok := s.channels[name]
	if !ok {
		s.channels[name] = channel
	} else {
		// error  already exists
	}
}

func (s *Socket) UnRegisterChannel(name string) {
	delete(s.channels, name)
}

func Init(url string, config *config.Config) *Socket {
	// operation on url
	url = "ws://" + url + "/v1/api/" + config.Project + "/socket/json"
	if config.IsSecure {
		url = "wss://" + url + "/v1/api/" + config.Project + "/socket/json"
	}
	return &Socket{
		url:             url,
		options:         websocketOptions{projectId: config.Project},
		channels: map[string]chan FunctionsPayload{},
		serviceRegister: map[string]chan map[string]bool{},
		pendingMsg:      []WriteMessageStructure{},
		mux:             sync.RWMutex{},
	}
}

func (s *Socket) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(s.url, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	s.mux.Lock()
	s.Socket = conn
	s.setConnected(true)
	s.mux.Unlock()

	return nil
}

func (s *Socket) ServiceRegister(serviceName string) error {
	data, err := s.Request(ServiceRegister, ServiceRegisterRequest{
		Service: serviceName,
		Project: s.options.projectId,
		Token:   s.options.token,
	})
	if err != nil {
		return err
	}

	if !data {
		log.Println("Could not connect to service")
		return nil
	}
	log.Println("Service started successfully")
	return nil
}

func (s *Socket) Request(msgType string, data interface{}) (bool, error) {
	id := uuid.NewV1().String()

	if !s.getConnected() {
		// connect to server
		if err := s.connect(); err != nil {
			return false, err
		}
		go s.read()
	}

	// create a ws writer listening on writeMessage channel
	writeMessage := make(chan WriteMessageStructure)
	defer close(writeMessage)
	go s.writeMessage(writeMessage)
	s.WriteMessage = writeMessage

	if len(s.pendingMsg) != 0 {
		for _, msg := range s.pendingMsg {
			writeMessage <- msg
		}
		s.pendingMsg = []WriteMessageStructure{}
	}

	writeMessage <- WriteMessageStructure{
		Type: msgType,
		Data: Message{Type: msgType, Data: data, ID: id},
	}

	timer1 := time.NewTimer(1 * time.Second)
	defer timer1.Stop()

	// channel for receiving service register acknowledgement
	ch := make(chan map[string]bool)
	defer close(ch)

	s.mux.Lock()
	s.serviceRegister[msgType] = ch
	s.mux.Unlock()

	for {
		select {
		case <-timer1.C:
			return false, errors.New("response time elapsed")
		case msg := <-ch:
			return msg["ack"], nil
		}
	}
}

func (s *Socket) read() {
	msg := Message{}
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

		switch v := msg.Data.(type) {
		case map[string]bool:
			ch, ok := s.serviceRegister[msg.Type]
			if ok {
				ch <- v
			}
		case FunctionsPayload:
			if msg.ID != "" {
				ch, ok := s.channels[msg.ID]
				if ok {
					ch <- v
					s.UnRegisterChannel(msg.ID)
					return
				}
			}

			ch, ok := s.channels[msg.Type]
			if ok {
				ch <- v
			}
		}

	}
}

func (s Socket) writeMessage(sendMessage chan WriteMessageStructure) {
	for msg := range sendMessage {
		id := uuid.NewV1().String()
		data := Message{
			Type: msg.Type,
			Data: msg.Data,
			ID:   id,
		}
		if err := s.Socket.WriteJSON(data); err != nil {
			log.Println(err)
			s.pendingMsg = append(s.pendingMsg, msg)
			s.setConnected(false)
		}
	}
}
