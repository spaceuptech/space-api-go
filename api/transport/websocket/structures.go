package websocket

type websocketOptions struct {
	projectId string
	token     string
}

type WriteMessageStructure struct {
	Type string
	Data interface{}
}

type ServiceRegisterRequest struct {
	Service string `json:"service"`
	Project string `json:"project"`
	Token   string `json:"token"`
}

type Message struct {
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
