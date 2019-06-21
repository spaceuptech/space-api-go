package sql

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
func (get *Get) Where(conds ...utils.M) *Get {
	if len(conds) == 1 {
		get.find = utils.GenerateFind(conds[0])
	} else {
		get.find = utils.GenerateFind(utils.And(conds...))
	}
	return get
}

// Select returns fields selectively
func (get *Get) Select(sel map[string]int32) *Get {
	get.readOptions.Select = sel
	return get
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
func (get *Get) Skip(skip int) *Get {
	get.readOptions.Skip = int64(skip)
	return get
}

// Limit limits the number of results returned
func (get *Get) Limit(limit int) *Get {
	get.readOptions.Limit = int64(limit)
	return get
}

// Apply executes the operation and returns the result
func (get *Get) Apply() (*model.Response, error) {
	return get.config.Transport.Read(get.ctx, get.meta, get.find, get.op, get.readOptions)
}
