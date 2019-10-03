package config

import (
	"github.com/spaceuptech/space-api-go/api/transport"
)

// Config holds the config of the API object
type Config struct {
	Project   string
	Token     string
	IsSecure bool
	Transport *transport.Transport
}
