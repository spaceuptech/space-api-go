package service

import (
	"github.com/spaceuptech/space-api-go/api/model"
	"log"
)

func (s *Service) registerFunction(functionName string, fn model.Function) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.function[functionName] = fn
}

func (s *Service) getFunction(functionName string) (model.Function, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	fn, ok := s.function[functionName]
	return fn, ok
}

func (s *Service) serviceRegister(Type string, msg model.ServiceRegisterRequest) error {
	v, err := s.client.Request(model.ServiceRegister, model.ServiceRegisterRequest{Service: s.service, Project: s.options.Project, Token: s.options.Token})
	if err != nil {
		log.Println("Service register error:", err)
		return err
	}
	responseData, ok := v.(map[string]interface{})
	if !ok {
		log.Println("Service register error wrong type found")
		return err
	}
	ack, ok := responseData["ack"]
	if ok {
		log.Println("Service register error didn't received acknowledgement from server")
		return err
	}
	successful, ok := ack.(bool)
	if !successful {
		log.Println("Service register error could not connect to service")
		return err
	}
	log.Println("Service started successfully")
	return nil
}
