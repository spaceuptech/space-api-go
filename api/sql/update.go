package sql

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Update contains the methods for the update operation
type Update struct {
	ctx          context.Context
	meta         *proto.Meta
	op           string
	find, update utils.M
	config       *config.Config
}

func initUpdate(ctx context.Context, db, col, op string, config *config.Config) *Update {
	m := &proto.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token}
	f := make(utils.M)
	u := make(utils.M)
	return &Update{ctx, m, op, f, u, config}
}

// Where sets the where clause for the request
func (u *Update) Where(conds ...utils.M) *Update {
	if len(conds) == 1 {
		u.find = utils.GenerateFind(conds[0])
	} else {
		u.find = utils.GenerateFind(utils.And(conds...))
	}
	return u
}

// Set the value of the provided fields in the document
func (u *Update) Set(obj utils.M) *Update {
	u.update["$set"] = obj
	return u
}

// Apply executes the operation and returns the result
func (u *Update) Apply() (*model.Response, error) {
	return u.config.Transport.Update(u.ctx, u.meta, u.op, u.find, u.update)
}

func (u *Update) getProject() (string) {
	return u.config.Project
}
func (u *Update) getDb() (string) {
	return u.meta.DbType
}
func (u *Update) getToken() (string) {
	return u.config.Token
}
func (u *Update) getCollection() (string) {
	return u.meta.Col
}
func (u *Update) getOperation() (string) {
	return u.op
}
func (u *Update) getUpdate() (interface{}) {
	return u.update
}
func (u *Update) getFind() (interface{}) {
	return u.find
}