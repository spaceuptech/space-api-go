package db

import (
	"context"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/types"
)

// Get contains the methods for the get operation
type Get struct {
	readOptions *types.ReadOptions
	op          string
	find        types.M
	aggregate   map[string]map[string]string
	group       []interface{}
	config      *config.Config
	meta        *types.Meta
}

func initGet(db, col, op string, config *config.Config) *Get {
	meta := &types.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token, Operation: types.Read}
	f := make(types.M)
	return &Get{&types.ReadOptions{}, op, f, map[string]map[string]string{}, make([]interface{}, 0), config, meta}
}

// Where sets the where clause for the request
func (g *Get) Where(conds ...types.M) *Get {
	if len(conds) == 1 {
		g.find = types.GenerateFind(conds[0])
	} else {
		g.find = types.GenerateFind(types.And(conds...))
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
	g.readOptions.Sort = order
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

// Key sets the key for the distinct query
func (g *Get) GroupBy(values ...string) *Get {
	v := []interface{}{}
	for _, value := range values {
		v = append(v, value)
	}
	g.group = v
	return g
}

// Key sets the key for the distinct query
func (g *Get) AggregateCount(col string) *Get {
	g.aggregate = map[string]map[string]string{"aggregate": {"count": col}}
	return g
}

// Key sets the key for the distinct query
func (g *Get) AggregateMax(col string) *Get {
	g.aggregate = map[string]map[string]string{"aggregate": {"max": col}}
	return g
}

// Key sets the key for the distinct query
func (g *Get) AggregateMin(col string) *Get {
	g.aggregate = map[string]map[string]string{"aggregate": {"min": col}}
	return g
}

// Key sets the key for the distinct query
func (g *Get) AggregateAverage(col string) *Get {
	g.aggregate = map[string]map[string]string{"aggregate": {"avg": col}}
	return g
}

// Key sets the key for the distinct query
func (g *Get) AggregateSum(col string) *Get {
	g.aggregate = map[string]map[string]string{"aggregate": {"sum": col}}
	return g
}

// Apply executes the operation and returns the result
func (g *Get) Apply(ctx context.Context) (*types.Response, error) {
	return g.config.Transport.DoDBRequest(ctx, g.meta, g.createReadReq())
}

func (g *Get) createReadReq() *types.ReadRequest {
	return &types.ReadRequest{Find: g.find, Operation: g.op, Options: g.readOptions, Aggregate: g.aggregate, GroupBy: g.group}
}
