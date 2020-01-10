package db

import (
	"context"
	"strings"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/utils"
)

// Get contains the methods for the get operation
type Get struct {
	readOptions *model.ReadOptions
	op          string
	find        utils.M
	config      *config.Config
	meta        *model.Meta
}

func initGet(db, col, op string, config *config.Config) *Get {
	meta := &model.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token, Operation: utils.Read}
	f := make(utils.M)
	return &Get{&model.ReadOptions{}, op, f, config, meta}
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
	a := int64(skip)
	g.readOptions.Skip = &a
	return g
}

// Limit limits the number of results returned
func (g *Get) Limit(limit int) *Get {
	a := int64(limit)
	g.readOptions.Limit = &a
	return g
}

// Key sets the key for the distinct query
func (g *Get) Key(key string) *Get {
	g.readOptions.Distinct = &key
	return g
}

// Apply executes the operation and returns the result
func (g *Get) Apply(ctx context.Context) (*model.Response, error) {
	return g.config.Transport.DoDBRequest(ctx, g.meta, g.createReadReq())
}

func (g *Get) createReadReq() *model.ReadRequest {
	return &model.ReadRequest{Find: g.find, Operation: g.op, Options: g.readOptions}
}
