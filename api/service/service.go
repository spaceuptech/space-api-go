package service

import (
	"encoding/json"
	"log"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/transport"
	"github.com/spaceuptech/space-api-go/api/utils"
)

type ServiceFunction func(*model.Message, *model.Message, transport.CallBackFunction)

// Service contains the values for the service instance
type Service struct {
	config  *config.Config
	service string
	id      string
	funcs   map[string]ServiceFunction
	// options interface{}
	// store   interface{}
}

// RealtimeRequest is the object sent for realtime requests
type Payload struct {
	Auth     []byte `json:"auth"`
	Token    string `json:"token"`
	Project  string `json:"project"`
	Type     string `json:"type"`
	Id       string `json:"id"`
	Service  string `json:"service"`
	Params   []byte `json:"params"`
	Function string `json:"function"`
	Error    string `json:"error"`
}

func New(config *config.Config, service string) *Service {
	return nil
}

// func Init(config *config.Config, service string) *Service {
// 	id := uuid.NewV1().String()
// 	var w transport.WebsocketConnection
// 	var cb transport.CallBackFunctions
// 	w.RegisterOnReconnectCallback(cb("service-register", {} ))

// 	return &Service{config, service, id, make(map[string]ServiceFunction)}
// }

func (s *Service) RegisterFunc(funcName string, function ServiceFunction) {
	s.funcs[funcName] = function
}

func (s *Service) ServiceRequest(req Payload) {
	// function := req.Function
	// params := req.Params
	// auth := req.Auth
	// if (!auth || Object.keys(auth) == 0) auth = null;

}

// Start is used to start the particular service (is Blocking)
func (s *Service) Start() {

	con := s.config.Transport.GetWebsockConn()
	for {
		if con != nil {
			log.Println("Connected to Space Cloud")

			c := make(chan *Payload, 10)
			go func() {
				for payload := range c {
					if err := con.WriteJSON(payload); err != nil {
						log.Println("Error in sendin data  ", err)
					}
				}
			}()
			c <- &Payload{Service: s.service, Type: utils.TypeServiceRegister, Id: s.id, Project: s.config.Project, Token: s.config.Token}
			for {
				in := Payload{}
				err := con.ReadJSON(&in)
				if err != nil {
					log.Println("Error reading json.", err)
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
									c <- &Payload{Service: s.service, Type: utils.TypeServiceRequest, Id: in.Id, Error: "Error parsing the result"}
								} else {
									c <- &Payload{Service: s.service, Type: utils.TypeServiceRequest, Id: in.Id, Params: answer}
								}
							}
						})
					} else {
						c <- &Payload{Service: s.service, Type: utils.TypeServiceRequest, Id: in.Id, Error: "Function Not Registered"}
					}
				}
			}
			close(c)
		} else {
			log.Println("Not connected. Attempting to Reconnect...")

		}
	}
}
