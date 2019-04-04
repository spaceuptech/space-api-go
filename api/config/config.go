package config

import (
	"github.com/spaceuptech/space-api-go/api/proto"
)

// Config holds the config of the API object
type Config struct {
	Project    string
	Host, Port string
	Token      string
	Stub       proto.SpaceCloudClient
}
