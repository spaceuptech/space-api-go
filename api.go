package api

import (
	"context"
	"time"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/db"
	"github.com/spaceuptech/space-api-go/filestore"
	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/transport"
	"github.com/spaceuptech/space-api-go/transport/websocket"
)

// API is the main API object to communicate with space cloud
type API struct {
	config *config.Config
	socket *websocket.Socket
}

// New initialised a new instance of the API object
func New(project, url string, sslEnabled bool) *API {
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

// DB creates a db client instance
func (api *API) DB(dbAlias string) *db.DB {
	return db.New(dbAlias, api.config)
}

// Call invokes the specified function on the backend
func (api *API) Call(service, endpoint string, params interface{}, timeout int) (*model.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	return api.config.Transport.Call(ctx, api.config.Token, api.config.Project, service, endpoint, params, timeout)
}

// FileStore creates a FileStore instance
func (api *API) Filestore() *filestore.Filestore {
	return filestore.New(api.config)
}
