package model

const ServiceRegister string = "service-register"
const ServiceRequest string = "service-request"

type Function func(params interface{}, auth map[string]interface{}) (interface{}, error)

// ServiceRegisterRequest is a structure providing options for service register
type ServiceRegisterRequest struct {
	Service string `json:"service"`
	Project string `json:"project"`
	Token   string `json:"token"`
}

// WebsocketMessage is the body for a websocket request
type WebsocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	ID   string      `json:"id"` // the request id
}

//FunctionsPayload is the struct transmitted via the broker
type FunctionsPayload struct {
	ID      string                 `json:"id"`
	Auth    map[string]interface{} `json:"auth"`
	Params  interface{}            `json:"params"`
	Service string                 `json:"service"`
	Func    string                 `json:"func"`
	Error   string                 `json:"error"`
}
