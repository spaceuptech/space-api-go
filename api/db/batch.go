package db

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/transport"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// Delete contains the methods for the delete operation
type Batch struct {
	ctx      context.Context
	meta     *proto.Meta
	db       string
	config   *config.Config
	reqs     []model.AllRequest
	httpMeta *model.Meta
}

type Request interface {
	getProject() string
	getDb() string
	getToken() string
	getCollection() string
	getOperation() string
}

func initBatch(ctx context.Context, db string, config *config.Config) *Batch {
	m := &proto.Meta{DbType: db, Project: config.Project, Token: config.Token}
	meta := &model.Meta{DbType: db, Project: config.Project, Token: config.Token}
	return &Batch{ctx, m, db, config, []model.AllRequest{}, meta}
}

// Add adds a delete request to batch
func (b *Batch) Add(request *Request) error {
	req := *request
	if b.config.Project != req.getProject() {
		return errors.New("Cannot Batch Requests of Different Projects")
	}
	if b.db != req.getDb() {
		return errors.New("Cannot Batch Requests of Different Database Types")
	}
	if b.config.Token != req.getToken() {
		return errors.New("Cannot Batch Requests using Different Tokens")
	}
	allReq := model.AllRequest{}
	allReq.Col = req.getCollection()
	allReq.Operation = req.getOperation()
	if req, ok := req.(*Insert); ok {
		docJSON, err := json.Marshal(req.getDoc())
		if err != nil {
			return err
		}
		allReq.Document = docJSON
		allReq.Type = utils.Create
	}
	if req, ok := req.(*Update); ok {
		// findJSON, err := json.Marshal(req.getFind())
		// if err != nil {
		// 	return err
		// }
		allReq.Find = map[string]interface{}{"find": req.getFind()}
		// updateJSON, err := json.Marshal(req.getUpdate())
		// if err != nil {
		// 	return err
		// }
		allReq.Update = map[string]interface{}{"update": req.getUpdate()}
		allReq.Type = utils.Update
	}
	if req, ok := req.(*Delete); ok {
		// findJSON, err := json.Marshal(req.getFind())
		// if err != nil {
		// 	return err
		// }
		allReq.Find = map[string]interface{}{"find": req.getFind()}
		allReq.Type = utils.Delete
	}
	b.reqs = append(b.reqs, allReq)
	return nil
}

// Apply executes the operation and returns the result
func (b *Batch) Apply() {
	transport.Send("batch", b.createBatchReq(), b.httpMeta)
	//	return b.config.Transport.Batch(b.ctx, b.meta,&[]proto.AllRequest)
}

func (b *Batch) createBatchReq() *model.BatchRequest {
	return &model.BatchRequest{Requests: b.reqs}
}
