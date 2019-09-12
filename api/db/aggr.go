package db

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/transport"
)

// Aggr contains the methods for the aggregation operation
type Aggr struct {
	ctx      context.Context
	op       string
	pipeline []interface{}
	config   *config.Config
	meta     *model.Meta
}

func initAggr(ctx context.Context, db, col, op string, config *config.Config) *Aggr {
	meta := &model.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token}
	p := []interface{}{}
	return &Aggr{ctx, op, p, config, meta}
}

// Pipe sets the pipeline to run on the backend
func (a *Aggr) Pipe(pipeline []interface{}) *Aggr {
	a.pipeline = pipeline
	return a
}

// Apply executes the operation and returns the result
func (a *Aggr) Apply() (*model.Response, error) {
	transport.Send("aggr", a.createAggrReq(), a.meta)
	return &model.Response{}, nil
}

func (a *Aggr) createAggrReq() *model.AggregateRequest {
	return &model.AggregateRequest{Pipeline: a.pipeline, Operation: a.op}
}
