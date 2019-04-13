package transport

import (
	"github.com/spaceuptech/space-api-go/api/proto"
	"google.golang.org/grpc"
)

// Transport is the objct which handles all communication with the server
type Transport struct {
	stub proto.SpaceCloudClient
}

// Init initialises a new transport
func Init(host, port string) (*Transport, error) {
	conn, err := grpc.Dial(host + ":" + port)
	if err != nil {
		return nil, err
	}

	stub := proto.NewSpaceCloudClient(conn)
	return &Transport{stub}, nil
}
