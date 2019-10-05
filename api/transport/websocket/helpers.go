package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/spaceuptech/space-api-go/api/model"
)

func (s *Socket) setSocket(socket *websocket.Conn) {
	s.mux.Lock()
	s.Socket = socket
	s.mux.Unlock()
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

func (s *Socket) setWriterChannel(ch chan model.WebsocketMessage) {
	s.mux.Lock()
	s.SendMessage = ch
	s.mux.Unlock()
}

func (s *Socket) RegisterCallback(Type string, function func(data interface{})) {
	s.mux.Lock()
	s.registerCallbackMap[Type] = function
	s.mux.Unlock()
}

func (s *Socket) unregisterCallback(Type string) {
	s.mux.Lock()
	delete(s.registerCallbackMap, Type)
	s.mux.Unlock()
}

func (s *Socket) RegisterOnReconnectCallback(function func()) {
	s.mux.Lock()
	s.onReconnectCallbacks = append(s.onReconnectCallbacks, function)
	s.mux.Unlock()
}

func (s *Socket) setConnectedOnce(value bool) {
	s.mux.Lock()
	s.connectedOnce = value
	s.mux.Unlock()
}

func (s *Socket) getConnectedOnce() bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.connectedOnce
}
