package service

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/transport/websocket"
	"log"
	"sync"

	"github.com/spaceuptech/space-api-go/api/config"
)

type Service struct {
	options  *config.Config
	service  string
	client   *websocket.Socket
	function map[string]model.Function
	mux      sync.RWMutex
}

func Init(options *config.Config, serviceName string, client *websocket.Socket) (*Service, error) {
	return &Service{options: options, service: serviceName, client: client, function: map[string]model.Function{}}, nil
}

// RegisterFunc registers a functions
func (s *Service) RegisterFunc(functionName string, fn model.Function) {
	s.registerFunction(functionName, fn)
}

func (s *Service) serviceRequest(data interface{}) {
	req := model.FunctionsPayload{}
	if err := mapstructure.Decode(data, &req); err != nil {
		log.Println("Service request error ", err)
		s.client.Send(model.ServiceRequest, model.FunctionsPayload{ID: req.ID, Error: "Service request error" + err.Error()})
		return
	}

	functionName := req.Func
	params := req.Params
	auth := req.Auth

	if auth == nil || len(auth) == 0 {
		auth = nil
	}

	function, ok := s.getFunction(functionName)
	if !ok {
		log.Println("Service request no function registered on the service")
		s.client.Send(model.ServiceRequest, model.FunctionsPayload{ID: req.ID, Error: "No function registered on the service"})
		return
	}
	response, err := function(params, auth)
	if err != nil {
		log.Println("Service request error ", err)
		s.client.Send(model.ServiceRequest, model.FunctionsPayload{ID: req.ID, Error: "Service request error " + err.Error()})
		return
	}
	s.client.Send(model.ServiceRequest, model.FunctionsPayload{ID: req.ID, Params: response})
}

// Start
func (s *Service) Start() error {
	s.client.RegisterCallback(model.ServiceRequest, s.serviceRequest)
	s.client.RegisterOnReconnectCallback(func() {
		if err := s.serviceRegister(model.ServiceRegister, model.ServiceRegisterRequest{Service: s.service, Project: s.options.Project, Token: s.options.Token}); err != nil {
			log.Println(err)
		}
	})
	if err := s.serviceRegister(model.ServiceRegister, model.ServiceRegisterRequest{Service: s.service, Project: s.options.Project, Token: s.options.Token}); err != nil {
		return err
	}
	return nil
}
