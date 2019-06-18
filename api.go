package api

import "github.com/spaceuptech/space-api-go/api"

// New initialised a new instance of the API object
func New(project, url string, sslEnabled bool) (*api.API, error) {
	return api.Init(project, url, sslEnabled)
}
