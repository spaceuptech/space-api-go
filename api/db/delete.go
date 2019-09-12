package db

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/transport"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Delete contains the methods for the delete operation
type Delete struct {
	ctx      context.Context
	meta     *proto.Meta
	op       string
	find     utils.M
	config   *config.Config
	httpMeta *model.Meta
}

func initDelete(ctx context.Context, db, col, op string, config *config.Config) *Delete {
	m := &proto.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token}
	meta := &model.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token}
	f := make(utils.M)
	return &Delete{ctx, m, op, f, config, meta}
}

// Where sets the where clause for the request
func (d *Delete) Where(conds ...utils.M) *Delete {
	if len(conds) == 1 {
		d.find = utils.GenerateFind(conds[0])
	} else {
		d.find = utils.GenerateFind(utils.And(conds...))
	}
	return d
}

// Apply executes the operation and returns the result
func (d *Delete) Apply() (*model.Response, error) {
	transport.Send("delete", d.createDeleteReq(), d.httpMeta)
	return d.config.Transport.Delete(d.ctx, d.meta, d.op, d.find)
}

func (d *Delete) getProject() string {
	return d.config.Project
}
func (d *Delete) getDb() string {
	return d.httpMeta.DbType
}
func (d *Delete) getToken() string {
	return d.config.Token
}
func (d *Delete) getCollection() string {
	return d.httpMeta.Col
}
func (d *Delete) getOperation() string {
	return d.op
}
func (d *Delete) getFind() utils.M {
	return d.find
}

func (d *Delete) createDeleteReq() *model.DeleteRequest {
	return &model.DeleteRequest{Find: d.find, Operation: d.op}
}
