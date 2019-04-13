package mgo

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Delete contains the methods for the delete operation
type Delete struct {
	ctx    context.Context
	meta   *proto.Meta
	op     string
	find   utils.M
	config *config.Config
}

func initDelete(ctx context.Context, db, col, op string, config *config.Config) *Delete {
	m := &proto.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token}
	f := make(utils.M)
	return &Delete{ctx, m, op, f, config}
}

// Where sets the where clause for the request
func (d *Delete) Where(conds ...utils.M) *Delete {
	if len(conds) == 1 {
		d.find = GenerateFind(conds[0])
	} else {
		d.find = GenerateFind(utils.And(conds...))
	}
	return d
}

// Apply executes the operation and returns the result
func (d *Delete) Apply() (*model.Response, error) {
	return d.config.Transport.Delete(d.ctx, d.meta, d.op, d.find)
}
