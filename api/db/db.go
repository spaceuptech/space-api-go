package db

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/realtime"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// DB is the client responsible to commuicate with the DB crud module
type DB struct {
	config *config.Config
	db     string
}

// Init returns a DB client object
func Init(db string, config *config.Config) *DB {
	return &DB{config, db}
}

// Insert returns a helper to fire a insert request
func (d *DB) Insert(col string) *Insert {
	return initInsert(context.TODO(), d.db, col, d.config)
}

// Get returns a helper to fire a get all request
func (d *DB) Get(col string) *Get {
	return initGet(context.TODO(), d.db, col, utils.All, d.config)
}

// GetOne returns a helper to fire a get one request
func (d *DB) GetOne(col string) *Get {
	return initGet(context.TODO(), d.db, col, utils.One, d.config)
}

// Count returns a helper to fire a get count request
func (d *DB) Count(col string) *Get {
	return initGet(context.TODO(), d.db, col, utils.Count, d.config)
}

// Distinct returns a helper to fire a get distinct request
func (d *DB) Distinct(col string) *Get {
	return initGet(context.TODO(), d.db, col, utils.Distinct, d.config)
}

// Update returns a helper to fire an update all request
func (d *DB) Update(col string) *Update {
	return initUpdate(context.TODO(), d.db, col, utils.All, d.config)
}

// UpdateOne returns a helper to fire an update one request
func (d *DB) UpdateOne(col string) *Update {
	return initUpdate(context.TODO(), d.db, col, utils.One, d.config)
}

// Upsert returns a helper to fire an upsert request
func (d *DB) Upsert(col string) *Update {
	return initUpdate(context.TODO(), d.db, col, utils.Upsert, d.config)
}

// Delete returns a helper to fire a delete all request
func (d *DB) Delete(col string) *Delete {
	return initDelete(context.TODO(), d.db, col, utils.All, d.config)
}

// DeleteOne returns a helper to fire a delete one request
func (d *DB) DeleteOne(col string) *Delete {
	return initDelete(context.TODO(), d.db, col, utils.One, d.config)
}

// Aggr returns a helper to fire a aggregation (all) request
func (d *DB) Aggr(col string) *Aggr {
	return initAggr(context.TODO(), d.db, col, utils.All, d.config)
}

// AggrOne returns a helper to fire a aggregation (one) request
func (d *DB) AggrOne(col string) *Aggr {
	return initAggr(context.TODO(), d.db, col, utils.One, d.config)
}

// BeginBatch returns a helper to fire a batch request
func (d *DB) BeginBatch() *Batch {
	return initBatch(context.TODO(), d.db, d.config)
}

// LiveQuery returns a helper to fire a liveQuery request
func (d *DB) LiveQuery(col string) *realtime.LiveQuery {
	return realtime.Init(d.config, d.db, col)
}

// Profile fires a profile request
func (d *DB) Profile(id string) (*model.Response, error) {
	m := &proto.Meta{DbType: d.db, Project: d.config.Project, Token: d.config.Token}
	return d.config.Transport.Profile(context.TODO(), m, id)
}

// Profiles fires a profiles request
func (d *DB) Profiles() (*model.Response, error) {
	m := &proto.Meta{DbType: d.db, Project: d.config.Project, Token: d.config.Token}
	return d.config.Transport.Profiles(context.TODO(), m)
}

// SignIn fires a signIn request
func (d *DB) SignIn(email, password string) (*model.Response, error) {
	m := &proto.Meta{DbType: d.db, Project: d.config.Project, Token: d.config.Token}
	return d.config.Transport.SignIn(context.TODO(), m, email, password)
}

// SignUp fires a signUp request
func (d *DB) SignUp(email, name, password, role string) (*model.Response, error) {
	m := &proto.Meta{DbType: d.db, Project: d.config.Project, Token: d.config.Token}
	return d.config.Transport.SignUp(context.TODO(), m, email, name, password, role)
}

// EditProfile fires a editProfile request
func (d *DB) EditProfile(id string, values model.ProfileParams) (*model.Response, error) {
	m := &proto.Meta{DbType: d.db, Project: d.config.Project, Token: d.config.Token}
	return d.config.Transport.EditProfile(context.TODO(), m, id, values)
}
