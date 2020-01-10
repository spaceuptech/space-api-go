package db

import (
	"context"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/proto"
	"github.com/spaceuptech/space-api-go/realtime"
	"github.com/spaceuptech/space-api-go/utils"
)

// DB is the client responsible to commuicate with the DB crud module
type DB struct {
	config *config.Config
	db     string
}

// New returns a DB client object
func New(db string, config *config.Config) *DB {
	return &DB{config, db}
}

// Insert returns a helper to fire a insert request
func (d *DB) Insert(col string) *Insert {
	return initInsert(d.db, col, d.config)
}

// Get returns a helper to fire a get all request
func (d *DB) Get(col string) *Get {
	return initGet(d.db, col, utils.All, d.config)
}

// GetOne returns a helper to fire a get one request
func (d *DB) GetOne(col string) *Get {
	return initGet(d.db, col, utils.One, d.config)
}

// Count returns a helper to fire a get count request
func (d *DB) Count(col string) *Get {
	return initGet(d.db, col, utils.Count, d.config)
}

// Distinct returns a helper to fire a get distinct request
func (d *DB) Distinct(col string) *Get {
	return initGet(d.db, col, utils.Distinct, d.config)
}

// Update returns a helper to fire an update all request
func (d *DB) Update(col string) *Update {
	return initUpdate(d.db, col, utils.All, d.config)
}

// UpdateOne returns a helper to fire an update one request
func (d *DB) UpdateOne(col string) *Update {
	return initUpdate(d.db, col, utils.One, d.config)
}

// Upsert returns a helper to fire an upsert request
func (d *DB) Upsert(col string) *Update {
	return initUpdate(d.db, col, utils.Upsert, d.config)
}

// Delete returns a helper to fire a delete all request
func (d *DB) Delete(col string) *Delete {
	return initDelete(d.db, col, utils.All, d.config)
}

// DeleteOne returns a helper to fire a delete one request
func (d *DB) DeleteOne(col string) *Delete {
	return initDelete(d.db, col, utils.One, d.config)
}

// Aggr returns a helper to fire a aggregation (all) request
func (d *DB) Aggr(col string) *Aggr {
	return initAggr(d.db, col, utils.All, d.config)
}

// AggrOne returns a helper to fire a aggregation (one) request
func (d *DB) AggrOne(col string) *Aggr {
	return initAggr(d.db, col, utils.One, d.config)
}

// BeginBatch returns a helper to fire a batch request
func (d *DB) BeginBatch() *Batch {
	return initBatch(d.db, d.config)
}

// LiveQuery returns a helper to fire a liveQuery request
func (d *DB) LiveQuery(col string) *realtime.LiveQuery {
	return realtime.New(d.config, d.db, col)
}

// TODO: add support for the user management module
// // Profile fires a profile request
// func (d *DB) Profile(ctx context.Context, id string) (*model.Response, error) {
// 	m := &proto.Meta{DbType: d.db, Project: d.config.Project, Token: d.config.Token}
// 	return d.config.Transport.Profile(ctx, m, id)
// }
//
// // Profiles fires a profiles request
// func (d *DB) Profiles(ctx context.Context) (*model.Response, error) {
// 	m := &proto.Meta{DbType: d.db, Project: d.config.Project, Token: d.config.Token}
// 	return d.config.Transport.Profiles(ctx, m)
// }
//
// // SignIn fires a signIn request
// func (d *DB) SignIn(ctx context.Context,email, password string) (*model.Response, error) {
// 	m := &proto.Meta{DbType: d.db, Project: d.config.Project, Token: d.config.Token}
// 	return d.config.Transport.SignIn(ctx, m, email, password)
// }
//
// // SignUp fires a signUp request
// func (d *DB) SignUp(ctx context.Context, email, name, password, role string) (*model.Response, error) {
// 	m := &proto.Meta{DbType: d.db, Project: d.config.Project, Token: d.config.Token}
// 	return d.config.Transport.SignUp(ctx, m, email, name, password, role)
// }
//
// // EditProfile fires a editProfile request
// func (d *DB) EditProfile(ctx context.Context, id string, values model.ProfileParams) (*model.Response, error) {
// 	m := &proto.Meta{DbType: d.db, Project: d.config.Project, Token: d.config.Token}
// 	return d.config.Transport.EditProfile(ctx, m, id, values)
// }
