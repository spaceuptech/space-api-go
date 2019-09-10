package transport

import (
	"log"

	"github.com/spaceuptech/space-api-go/api/proto"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

// Transport is the objct which handles all communication with the server
type Transport struct {
	con  *websocket.Conn
	stub proto.SpaceCloudClient
	conn *grpc.ClientConn
}

type CallBackFunction func(string, interface{})

//Init initialises a new transport
func Init(addr string, sslEnabled bool) (*Transport, error) {
	var w *websocketConnection
	w.Init()
	err := w.Connect(addr, sslEnabled)
	if err != nil {
		log.Println("Error in establishing websocket connection", err)
	}
	return &Transport{nil, nil, nil}, nil
}

// GetConn returns the underlying gRPC client connection
func (t *Transport) GetConn() *grpc.ClientConn {
	return t.conn
}

func (t *Transport) GetWebsockConn() *websocket.Conn {
	return t.con
}

func (t *Transport) GetStub() proto.SpaceCloudClient {
	return t.stub
}
