package api

import (
	"google.golang.org/grpc"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/proto"
)

// API is the main API object to communicate with space cloud
type API struct {
	config *config.Config
}

// Init create a new instance of the API object
func Init(project, host, port string) (*API, error) {
	conn, err := grpc.Dial(host + ":" + port)
	if err != nil {
		return nil, err
	}

	stub := proto.NewSpaceCloudClient(conn)

	c := &config.Config{Project: project, Host: host, Port: port, Stub: stub}

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
func (api *API) Mongo() {

}

// MySQL returns a mysql client instance
func (api *API) MySQL() {

}

// Postgres creates a postgres client instance
func (api *API) Postgres() {

}
