package db

import (
	"context"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/utils"
)

// Aggr contains the methods for the aggregation operation
type Aggr struct {
	op       string
	pipeline []interface{}
	config   *config.Config
	meta     *model.Meta
}

func initAggr(db, col, op string, config *config.Config) *Aggr {
	meta := &model.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token, Operation: utils.Aggregate}
	p := make([]interface{}, 0)
	return &Aggr{op, p, config, meta}
}

// Pipe sets the pipeline to run on the backend
func (a *Aggr) Pipe(pipeline []interface{}) *Aggr {
	a.pipeline = pipeline
	return a
}

// Apply executes the operation and returns the result
func (a *Aggr) Apply(ctx context.Context) (*model.Response, error) {
	return a.config.Transport.DoDBRequest(ctx, a.meta, a.createAggrReq())
}

func (a *Aggr) createAggrReq() *model.AggregateRequest {
	return &model.AggregateRequest{Pipeline: a.pipeline, Operation: a.op}
}
