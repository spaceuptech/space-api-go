package transport

import (
	"crypto/tls"

	"github.com/spaceuptech/space-api-go/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Transport is the objct which handles all communication with the server
type Transport struct {
	stub proto.SpaceCloudClient
}

// Init initialises a new transport
func Init(host, port string, sslEnabled bool) (*Transport, error) {
	dialOptions := []grpc.DialOption{}

	if sslEnabled {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	conn, err := grpc.Dial(host+":"+port, dialOptions...)
	if err != nil {
		return nil, err
	}

	stub := proto.NewSpaceCloudClient(conn)
	return &Transport{stub}, nil
}
