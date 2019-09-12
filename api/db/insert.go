package db

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/transport"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Insert contains the methods for the create operation
type Insert struct {
	ctx      context.Context
	meta     *proto.Meta
	op       string
	obj      interface{}
	config   *config.Config
	httpMeta *model.Meta
}

func initInsert(ctx context.Context, db, col string, config *config.Config) *Insert {
	m := &proto.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token}
	meta := &model.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token}
	return &Insert{ctx: ctx, meta: m, config: config, httpMeta: meta}
}

// Docs sets the documents to be inserted into the database
func (i *Insert) Docs(docs interface{}) *Insert {
	i.op = utils.All
	i.obj = docs
	return i
}

// Doc sets the document to be inserted into the database
func (i *Insert) Doc(doc interface{}) *Insert {
	i.op = utils.One
	i.obj = doc
	return i
}

// Apply executes the operation and returns the result
func (i *Insert) Apply() (*model.Response, error) {
	transport.Send("create", i.createCreateReq(), i.httpMeta)
	return i.config.Transport.Insert(i.ctx, i.httpMeta, i.meta, i.op, i.obj)
}

func (i *Insert) getProject() string {
	return i.config.Project
}
func (i *Insert) getDb() string {
	return i.httpMeta.DbType
}
func (i *Insert) getToken() string {
	return i.config.Token
}
func (i *Insert) getCollection() string {
	return i.httpMeta.Col
}
func (i *Insert) getOperation() string {
	return i.op
}
func (i *Insert) getDoc() interface{} {
	return i.obj
}

func (i *Insert) createCreateReq() *model.CreateRequest {
	return &model.CreateRequest{Document: i.obj, Operation: i.op}
}
