package websocket

import (
	"github.com/gorilla/websocket"

	"github.com/spaceuptech/space-api-go/model"
)

func (s *Socket) setSocket(socket *websocket.Conn) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.socket = socket
}

func (s *Socket) checkIsConnecting() bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.isConnecting {
		return false
	}

	s.isConnecting = true
	return true
}

func (s *Socket) resetIsConnecting() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.isConnecting = false
}

func (s *Socket) getConnected() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.isConnect
}

func (s *Socket) setConnected(value bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.isConnect = value
}

func (s *Socket) setWriterChannel(ch chan model.WebsocketMessage) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.sendMessage = ch
}

func (s *Socket) unregisterCallback(Type string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.registerCallbackMap, Type)
}

func (s *Socket) setConnectedOnce(value bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.connectedOnce = value
}

func (s *Socket) isConnectedOnce() bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.connectedOnce
}

func (s *Socket) getRegisteredCallBack(Type string) (func(data interface{}), bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	fn, ok := s.registerCallbackMap[Type]
	return fn, ok
}

func (s *Socket) addPendingMsg(msg model.WebsocketMessage) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.pendingMsg = append(s.pendingMsg, msg)
}
