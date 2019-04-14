package mgo

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
		u.find = GenerateFind(conds[0])
	} else {
		u.find = GenerateFind(utils.And(conds...))
	}
	return u
}

// Set the value of the provided fields in the document
func (u *Update) Set(obj utils.M) *Update {
	u.update["$set"] = obj
	return u
}

// Push adds an item to an array
func (u *Update) Push(obj utils.M) *Update {
	u.update["$push"] = obj
	return u
}

// Remove removes the specified field from a document
func (u *Update) Remove(fields ...string) *Update {
	obj := make(utils.M, len(fields))
	for _, field := range fields {
		obj[field] = 1
	}
	u.update["$unset"] = obj
	return u
}

// Rename renames the specified field
func (u *Update) Rename(obj utils.M) *Update {
	u.update["$rename"] = obj
	return u
}

// Inc increments the value of the field by the specified amount
func (u *Update) Inc(obj utils.M) *Update {
	u.update["$inc"] = obj
	return u
}

// Mul multiplies the value of the field by the specified amount
func (u *Update) Mul(obj utils.M) *Update {
	u.update["$mul"] = obj
	return u
}

// Max updates the field if the specified value is greater than the existing field value
func (u *Update) Max(obj utils.M) *Update {
	u.update["$max"] = obj
	return u
}

// Min updates the field if the specified value is lesser than the existing field value
func (u *Update) Min(obj utils.M) *Update {
	u.update["$min"] = obj
	return u
}

// CurrentTimestamp sets the value of a field to current timestamp
func (u *Update) CurrentTimestamp(fields ...string) *Update {
	objTemp, p := u.update["$currentDate"]
	if !p {
		objTemp = utils.M{}
	}

	obj := objTemp.(utils.M)
	for _, field := range fields {
		obj[field] = utils.M{"$type": "timestamp"}
	}

	u.update["$currentDate"] = obj
	return u
}

// CurrentDate sets the value of a field to current date
func (u *Update) CurrentDate(fields ...string) *Update {
	objTemp, p := u.update["$currentDate"]
	if !p {
		objTemp = utils.M{}
	}

	obj := objTemp.(utils.M)
	for _, field := range fields {
		obj[field] = utils.M{"$type": "date"}
	}

	u.update["$currentDate"] = obj
	return u
}

// Apply executes the operation and returns the result
func (u *Update) Apply() (*model.Response, error) {
	return u.config.Transport.Update(u.ctx, u.meta, u.op, u.find, u.update)
}
