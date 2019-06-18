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
	conn *grpc.ClientConn
}

// Init initialises a new transport
func Init(url string, sslEnabled bool) (*Transport, error) {
	dialOptions := []grpc.DialOption{}

	if sslEnabled {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(url, dialOptions...)
	if err != nil {
		return nil, err
	}

	stub := proto.NewSpaceCloudClient(conn)
	return &Transport{stub, conn}, nil
}

// GetStub returns the underlying gRPC stub
func (t *Transport) GetStub() proto.SpaceCloudClient {
	return t.stub
}

// GetConn returns the underlying gRPC client connection
func (t *Transport) GetConn() *grpc.ClientConn {
	return t.conn
}
