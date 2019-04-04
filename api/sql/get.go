package sql

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/mgo"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/transport"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Get contains the methods for the get operation
type Get struct {
	meta        *proto.Meta
	readOptions *proto.ReadOptions
	op          string
	find        utils.M
	config      *config.Config
}

func initGet(db, col string, config *config.Config) *Get {
	m := &proto.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token}
	r := new(proto.ReadOptions)
	op := "all"
	f := make(map[string]interface{})
	return &Get{m, r, op, f, config}
}

// Where sets the where clause for the request
func (get *Get) Where(conds ...utils.M) *Get {
	if len(conds) == 1 {
		get.find = mgo.GenerateFind(conds[0])
	} else {
		get.find = mgo.GenerateFind(utils.And(conds...))
	}
	return get
}

// Select returns fields selectively
func (get *Get) Select(sel map[string]int32) *Get {
	get.readOptions.Select = sel
	return get
}

// Sort sorts the result
func (get *Get) Sort(order map[string]int32) *Get {
	get.readOptions.Sort = order
	return get
}

// Skip skips some of the result
func (get *Get) Skip(skip int) *Get {
	get.readOptions.Skip = int64(skip)
	return get
}

// Limit limits the number of results returned
func (get *Get) Limit(limit int) *Get {
	get.readOptions.Limit = int64(limit)
	return get
}

func (get *Get) one() (*model.Response, error) {
	get.op = "one"
	return transport.Read(context.TODO(), get.config.Stub, get.meta, get.find, get.op, get.readOptions)
}

func (get *Get) all() (*model.Response, error) {
	get.op = "all"
	return transport.Read(context.TODO(), get.config.Stub, get.meta, get.find, get.op, get.readOptions)
}
