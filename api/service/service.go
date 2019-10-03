package service

import (
	"sync"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/transport/websocket"
)

type Service struct {
	options  *config.Config
	service  string
	client   *websocket.Socket
	function functionInfo
	mux      sync.RWMutex
}

type functionInfo struct {
	storeFunction map[string]func(params interface{}, auth interface{})
	response      FunctionResponse
}

type FunctionResponse struct {
	Type    string
	Message interface{}
}

func Init(options *config.Config, serviceName string, client *websocket.Socket) (*Service, error) {
	if err := client.ServiceRegister(serviceName); err != nil {
		return nil, err
	}
	return &Service{
		options: options,
		service: serviceName,
		client:  client,
		function: functionInfo{
			storeFunction: map[string]func(params interface{}, auth interface{}){},
			response:      FunctionResponse{},
		},
	}, nil
	// service request ?
}

func (s *Service) RegisterFunc(functionName string, function func(params interface{}, auth interface{}), response FunctionResponse) {
	s.function.storeFunction[functionName] = function
	s.function.response = response
}

func (s *Service) serviceRequest(ch chan websocket.FunctionsPayload) {
	for req := range ch {
		functionName := req.Func
		params := req.Params
		auth := req.Auth

		if auth == nil || len(auth) == 0 {
			auth = nil
		}

		s.mux.RLock()
		funcInfo := s.function
		s.mux.RUnlock()

		function, ok := funcInfo.storeFunction[functionName]
		if !ok {
			s.client.WriteMessage <- websocket.WriteMessageStructure{
				Type: websocket.ServiceRequest,
				Data: websocket.FunctionsPayload{ID: req.ID, Error: "No function registered on the function"},
			}
		}

		function(params, auth)

		switch funcInfo.response.Type {
		case "response":
			s.mux.Lock()
			s.client.WriteMessage <- websocket.WriteMessageStructure{
				Type: websocket.ServiceRequest,
				Data: websocket.FunctionsPayload{ID: req.ID, Params: req.Params},
			}
			s.mux.Unlock()
		}
	}
}

func (s *Service) Start() error {
	if err := s.client.ServiceRegister(s.service); err != nil {
		return err
	}

	ch := make(chan websocket.FunctionsPayload)
	go s.serviceRequest(ch)
	s.client.RegisterChannel(websocket.ServiceRequest, ch)

	return nil
}
