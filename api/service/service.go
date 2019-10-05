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
	function map[string]function
	mux      sync.RWMutex
}

type function func(params interface{}, auth map[string]interface{}) (interface{}, error)

func Init(options *config.Config, serviceName string, client *websocket.	Socket) (*Service, error) {
	client.RegisterOnReconnectCallback(func() {
		v, err := client.Request(websocket.ServiceRegister, model.ServiceRegisterRequest{Service: serviceName, Project: options.Project, Token: options.Token})
		if err != nil {
			log.Println(err)
			return
		}
		responseData, ok := v.(map[string]bool)
		if !ok {
			log.Println("service initialization wrong type found")
			return
		}
		data, ok := responseData["ack"]
		if ok {
			log.Println("service initialization didn't received acknowledgement from knowledge")
			return
		}
		if !data {
			log.Println("Could not connect to service")
			return
		}
		log.Println("Service started successfully")
	})

	return &Service{options: options, service: serviceName, client: client, function: map[string]function{}}, nil
}

func (s *Service) RegisterFunc(functionName string, fn function) {
	s.function[functionName] = fn
}

func (s *Service) serviceRequest(data interface{}) {
	reqMap, ok := data.(map[string]interface{})
	if !ok {
		log.Println("service request type not found")
	}

	req := model.FunctionsPayload{}
	if err := mapstructure.Decode(reqMap, &req); err != nil {
		log.Println(err)
	}

	functionName := req.Func
	params := req.Params
	auth := req.Auth

	if auth == nil || len(auth) == 0 {
		auth = nil
	}

	s.mux.RLock()
	funcInfo := s.function
	s.mux.RUnlock()

	function, ok := funcInfo[functionName]
	if !ok {
		s.client.Send(websocket.ServiceRequest, model.FunctionsPayload{ID: req.ID, Error: "No function registered on the service"})
		return
	}
	response, err := function(params, auth)
	if err != nil {
		log.Println(err)
	}
	s.client.Send(websocket.ServiceRequest, model.FunctionsPayload{ID: req.ID, Params: response})
}

func (s *Service) Start() error {
	s.client.RegisterCallback(websocket.ServiceRequest, s.serviceRequest)

	v, err := s.client.Request(websocket.ServiceRegister, model.ServiceRegisterRequest{Service: s.service, Project: s.options.Project, Token: s.options.Token,})
	if err != nil {
		return err
	}
	responseData, ok := v.(map[string]bool)
	if !ok {

	}
	data, ok := responseData["ack"]
	if ok {

	}
	if !data {
		log.Println("Could not connect to service")
		return nil
	}
	log.Println("Service started successfully")

	return nil
}
