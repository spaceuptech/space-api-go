package transport

import (
	"net/http"

	"github.com/spaceuptech/space-api-go/api/proto"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

// Transport is the objct which handles all communication with the server
type Transport struct {
	// Transport variables
	sslEnabled bool
	addr       string

	// Client drivers
	httpClient *http.Client
	con        *websocket.Conn
	stub       proto.SpaceCloudClient
	conn       *grpc.ClientConn
}

type CallBackFunction func(string, interface{})

// New initialises a new transport
func New(addr string, sslEnabled bool) *Transport {

	return &Transport{
		sslEnabled: sslEnabled,
		addr:       addr,
		httpClient: &http.Client{},
	}
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
