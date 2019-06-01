package api

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/mgo"
	"github.com/spaceuptech/space-api-go/api/sql"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/transport"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// API is the main API object to communicate with space cloud
type API struct {
	config *config.Config
}

// Init initialised a new instance of the API object
func Init(project, host, port string, sslEnabled bool) (*API, error) {
	t, err := transport.Init(host, port, sslEnabled)
	if err != nil {
		return nil, err
	}
	c := &config.Config{Project: project, Transport: t}

	return &API{c}, err
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
func (api *API) Mongo() *mgo.Mongo {
	return mgo.Init(api.config)
}

// MySQL returns a mysql client instance
func (api *API) MySQL() *sql.SQL {
	return sql.Init(utils.MySQL, api.config)
}

// Postgres creates a postgres client instance
func (api *API) Postgres() *sql.SQL {
	return sql.Init(utils.Postgres, api.config)
}

// Call invokes the specified function on the backend
func (api *API) Call(service, function string, params utils.M, timeout int) (*model.Response, error) {
	return api.config.Transport.Call(context.TODO(), api.config.Token, service, function, params, timeout)
}
