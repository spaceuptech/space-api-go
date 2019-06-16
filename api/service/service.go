package service

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log"
	"io"
	"encoding/json"

	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/utils"
)

type CallBackFunction func(string, interface{})()
type ServiceFunction func(*model.Message, *model.Message, CallBackFunction)()

// Service contains the methods for the service instance
type Service struct {
	config  *config.Config
	service string
	id      string
	funcs   map[string]ServiceFunction
}

func Init(config *config.Config, service string) *Service {
	id := uuid.NewV1().String()
	return &Service{config, service, id, make(map[string]ServiceFunction)}
}

func (s *Service) RegisterFunc(funcName string, function ServiceFunction) {
	s.funcs[funcName] = function
}

func (s *Service) Start() {
	registerRequest := &proto.FunctionsPayload{Service: s.service, Type: utils.TypeServiceRegister, Id: s.id, Project: s.config.Project, Token: s.config.Token}
	conn := s.config.Transport.GetConn()
	for {
		state := conn.GetState()
		if state.String() == "READY" {
			log.Println("Connected to Space Cloud")
			stream, err := s.config.Transport.GetStub().Service(context.TODO())
			if err != nil {
				continue
			}
			c := make(chan *proto.FunctionsPayload)
			go func() {
				for payload := range c {
					if err := stream.Send(payload); err != nil {
						log.Println("Failed to send a request:", err)
					}
				}
			}()
			c <- registerRequest
			for {
				in, err := stream.Recv()
				if err == io.EOF {
					close(c)
					return
				}
				if err != nil {
					log.Println("Failed to receive a request:", err)
					break
				}
				if in.Type == utils.TypeServiceRegister {
					if in.Id == s.id {
						temp := make(map[string]bool)
						json.Unmarshal(in.Params, &temp)
						if temp["ack"] {
							log.Println("Service registered")
						} else {
							close(c)
							panic("Could Not Register Service")
						}
					}
				} else if in.Type == utils.TypeServiceRequest {
					function, found := s.funcs[in.Function]
					if found {
						go function(&model.Message{in.Params}, &model.Message{in.Auth}, func(kind string, res interface{}) {
							if kind != "response" {
								close(c)
								panic("Type must be 'response'")
							} else {
								answer, err := json.Marshal(res)
								if err != nil {
									log.Println("Could not marshal the result", res)
									c <- &proto.FunctionsPayload{Service: s.service, Type: utils.TypeServiceRequest, Id: in.Id, Error: "Error parsing the result"}
								} else {
									c <- &proto.FunctionsPayload{Service: s.service, Type: utils.TypeServiceRequest, Id: in.Id, Params: answer}
								}
							}
						})
					} else {
						c <- &proto.FunctionsPayload{Service: s.service, Type: utils.TypeServiceRequest, Id: in.Id, Error: "Function Not Registered"}
					}
				}
			}
			close(c)
		} else {
			log.Println("Not connected. Attempting to Reconnect...")
			conn.WaitForStateChange(context.TODO(), state)
		}
	}
}