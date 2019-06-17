package mgo

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/utils"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/realtime"
)

// Mongo is the client responsible to commuicate with the Mongo crud module
type Mongo struct {
	config *config.Config
	db     string
}

// Init returns a Mongo client object
func Init(config *config.Config) *Mongo {
	return &Mongo{db: utils.Mongo, config: config}
}

// Insert returns a helper to fire a insert request
func (s *Mongo) Insert(col string) *Insert {
	return initInsert(context.TODO(), s.db, col, s.config)
}

// Get returns a helper to fire a get all request
func (s *Mongo) Get(col string) *Get {
	return initGet(context.TODO(), s.db, col, utils.All, s.config)
}

// GetOne returns a helper to fire a get one request
func (s *Mongo) GetOne(col string) *Get {
	return initGet(context.TODO(), s.db, col, utils.One, s.config)
}

// Count returns a helper to fire a get count request
func (s *Mongo) Count(col string) *Get {
	return initGet(context.TODO(), s.db, col, utils.Count, s.config)
}

// Distinct returns a helper to fire a get distinct request
func (s *Mongo) Distinct(col string) *Get {
	return initGet(context.TODO(), s.db, col, utils.Distinct, s.config)
}

// Update returns a helper to fire an update all request
func (s *Mongo) Update(col string) *Update {
	return initUpdate(context.TODO(), s.db, col, utils.All, s.config)
}

// UpdateOne returns a helper to fire an update one request
func (s *Mongo) UpdateOne(col string) *Update {
	return initUpdate(context.TODO(), s.db, col, utils.One, s.config)
}

// Upsert returns a helper to fire an upsert request
func (s *Mongo) Upsert(col string) *Update {
	return initUpdate(context.TODO(), s.db, col, utils.Upsert, s.config)
}

// Delete returns a helper to fire a delete all request
func (s *Mongo) Delete(col string) *Delete {
	return initDelete(context.TODO(), s.db, col, utils.All, s.config)
}

// DeleteOne returns a helper to fire a delete one request
func (s *Mongo) DeleteOne(col string) *Delete {
	return initDelete(context.TODO(), s.db, col, utils.One, s.config)
}

// Aggr returns a helper to fire a aggregation request
func (s *Mongo) Aggr(col string) *Aggr {
	return initAggr(context.TODO(), s.db, col, utils.All, s.config)
}

// AggrOne returns a helper to fire a aggregation request
func (s *Mongo) AggrOne(col string) *Aggr {
	return initAggr(context.TODO(), s.db, col, utils.One, s.config)
}

// BeginBatch returns a helper to fire a batch request
func (s *Mongo) BeginBatch() *Batch {
	return initBatch(context.TODO(), s.db, s.config)
}

// LiveQuery returns a helper to fire a liveQuery request
func (s *Mongo) LiveQuery(col string) *realtime.LiveQuery {
	return realtime.Init(s.config, s.db, col)
}

// Profile fires a profile request
func (s *Mongo) Profile(id string) (*model.Response, error) {
	m := &proto.Meta{DbType: s.db, Project: s.config.Project, Token: s.config.Token}
	return s.config.Transport.Profile(context.TODO(), m, id)
}

// Profiles fires a profiles request
func (s *Mongo) Profiles() (*model.Response, error) {
	m := &proto.Meta{DbType: s.db, Project: s.config.Project, Token: s.config.Token}
	return s.config.Transport.Profiles(context.TODO(), m)
}

// SignIn fires a signIn request
func (s *Mongo) SignIn(email, password string) (*model.Response, error) {
	m := &proto.Meta{DbType: s.db, Project: s.config.Project, Token: s.config.Token}
	return s.config.Transport.SignIn(context.TODO(), m, email, password)
}

// SignUp fires a signUp request
func (s *Mongo) SignUp(email, name, password, role string) (*model.Response, error) {
	m := &proto.Meta{DbType: s.db, Project: s.config.Project, Token: s.config.Token}
	return s.config.Transport.SignUp(context.TODO(), m, email, name, password, role)
}

// EditProfile fires a editProfile request
func (s *Mongo) EditProfile(id string, values model.ProfileParams) (*model.Response, error) {
	m := &proto.Meta{DbType: s.db, Project: s.config.Project, Token: s.config.Token}
	return s.config.Transport.EditProfile(context.TODO(), m, id, values)
}
