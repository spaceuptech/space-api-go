package websocket

// RegisterOnReconnectCallback registers a callback
func (s *Socket) RegisterOnReconnectCallback(function func()) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.onReconnectCallbacks = append(s.onReconnectCallbacks, function)
}

// RegisterCallback registers a callback
func (s *Socket) RegisterCallback(Type string, function func(data interface{})) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.registerCallbackMap[Type] = function
}
