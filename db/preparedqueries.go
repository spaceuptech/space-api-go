package db

import (
	"context"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/types"
)

// PreparedQuery contains the methods for the PreparedQueries operation
type PreparedQuery struct {
	args     map[string]interface{}
	config   *config.Config
	httpMeta *types.Meta
}

func initPreparedQuery(db, id string, config *config.Config) *PreparedQuery {
	meta := &types.Meta{ID: id, DbType: db, Project: config.Project, Token: config.Token, Operation: types.PreparedQueries}
	return &PreparedQuery{config: config, httpMeta: meta}
}

//Args sets the Arguments to be passed to prepared query
func (p *PreparedQuery) Args(args map[string]interface{}) *PreparedQuery {
	p.args = args
	return p
}

// Apply executes the operation and returns the result
func (p *PreparedQuery) Apply(ctx context.Context) (*types.Response, error) {
	return p.config.Transport.DoDBRequest(ctx, p.httpMeta, p.preparedQueryeReq())
}

func (p *PreparedQuery) preparedQueryeReq() *types.PreparedQueryRequest {
	return &types.PreparedQueryRequest{Params: p.args}
}
