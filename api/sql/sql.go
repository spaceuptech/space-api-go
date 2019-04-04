package sql

import (
	"github.com/spaceuptech/space-api-go/api/config"
)

// SQL is the client responsible to commuicate with the SQL crud module
type SQL struct {
	config *config.Config
	db     string
}

// Init returns a SQL client object
func Init(db string, config *config.Config) *SQL {
	return &SQL{db: db, config: config}
}
