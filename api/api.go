package api

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/db"
	"github.com/spaceuptech/space-api-go/api/transport/websocket"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/filestore"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/pubsub"
	"github.com/spaceuptech/space-api-go/api/service"
	"github.com/spaceuptech/space-api-go/api/transport"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// API is the main API object to communicate with space cloud
type API struct {
	config *config.Config
	socket *websocket.Socket
}

// Init initialised a new instance of the API object
func Init(project, url string, sslEnabled bool) *API {
	t := transport.New(url, sslEnabled)
	c := &config.Config{Project: project, Transport: t}
	w := websocket.Init(url, c)
	return &API{config: c, socket: w}
}

// SetToken sets the JWT token to be used in each request
func (api *API) SetToken(token string) {
	api.config.Token = token
}

// SetProjectID sets the project id to be used by the API
func (api *API) SetProjectID(project string) {
	api.config.Project = project
}

// Mongo returns a mongo db client instance
func (api *API) Mongo() *db.DB {
	return db.Init(utils.Mongo, api.config)
}

// MySQL returns a mysql client instance
func (api *API) MySQL() *db.DB {
	return db.Init(utils.MySQL, api.config)
}

// Postgres creates a postgres client instance
func (api *API) Postgres() *db.DB {
	return db.Init(utils.Postgres, api.config)
}

// Call invokes the specified function on the backend
func (api *API) Call(service, function string, params utils.M, timeout int) (*model.Response, error) {
	return api.config.Transport.Call(context.TODO(), api.config.Token, service, function, params, timeout)
}

// Service creates a Service instance
func (api *API) Service(config *config.Config, serviceName string) (*service.Service, error) {
	val, err := service.Init(config, serviceName, api.socket)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// FileStore creates a FileStore instance
func (api *API) Filestore() *filestore.Filestore {
	return filestore.Init(api.config)
}

// Pubsub creates a Pubsub instance
func (api *API) Pubsub(url string) *pubsub.Pubsub {
	return pubsub.Init(url, api.config.Project, api.socket)
}
