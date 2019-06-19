package mgo

import (
	"context"
	"strings"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Get contains the methods for the get operation
type Get struct {
	ctx         context.Context
	meta        *proto.Meta
	readOptions *proto.ReadOptions
	op          string
	find        utils.M
	config      *config.Config
}

func initGet(ctx context.Context, db, col, op string, config *config.Config) *Get {
	m := &proto.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token}
	r := new(proto.ReadOptions)
	f := make(utils.M)
	return &Get{ctx, m, r, op, f, config}
}

// Where sets the where clause for the request
func (g *Get) Where(conds ...utils.M) *Get {
	if len(conds) == 1 {
		g.find = utils.GenerateFind(conds[0])
	} else {
		g.find = utils.GenerateFind(utils.And(conds...))
	}
	return g
}

// Select returns fields selectively
func (g *Get) Select(sel map[string]int32) *Get {
	g.readOptions.Select = sel
	return g
}

// Sort sorts the result
func (g *Get) Sort(order ...string) *Get {
	ord := make(map[string]int32)
	for _, o := range order {
		if strings.HasPrefix(o, "-") {
			ord[o[1:]] = -1
		} else {
			ord[o] = 1
		}
	}
	g.readOptions.Sort = ord
	return g
}

// Skip skips some of the result
func (g *Get) Skip(skip int) *Get {
	g.readOptions.Skip = int64(skip)
	return g
}

// Limit limits the number of results returned
func (g *Get) Limit(limit int) *Get {
	g.readOptions.Limit = int64(limit)
	return g
}

// Key sets the key for the distinct query
func (g *Get) Key(key string) *Get {
	g.readOptions.Distinct = key
	return g
}

// Apply executes the operation and returns the result
func (g *Get) Apply() (*model.Response, error) {
	return g.config.Transport.Read(g.ctx, g.meta, g.find, g.op, g.readOptions)
}
