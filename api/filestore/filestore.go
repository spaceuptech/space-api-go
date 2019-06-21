package filestore

import (
	"context"
	"io"
	
	"github.com/spaceuptech/space-api-go/api/config"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/model"
)

// Filestore contains the values for the filestore instance
type Filestore struct {
	config  *config.Config
}

// Init initializes the filestore module
func Init(config *config.Config) *Filestore {
	return &Filestore{config}
}

func (f *Filestore) CreateFolder(path, name string) (*model.Response, error) {
	m := &proto.Meta{Project: f.config.Project, Token: f.config.Token}
	return f.config.Transport.CreateFolder(context.TODO(), m, path, name)
}

func (f *Filestore) DeleteFile(path string) (*model.Response, error) {
	m := &proto.Meta{Project: f.config.Project, Token: f.config.Token}
	return f.config.Transport.DeleteFile(context.TODO(), m, path)
}

func (f *Filestore) ListFiles(path string) (*model.Response, error) {
	m := &proto.Meta{Project: f.config.Project, Token: f.config.Token}
	return f.config.Transport.ListFiles(context.TODO(), m, path)
}

func (f *Filestore) UploadFile(path, name string, reader io.Reader) (*model.Response, error) {
	m := &proto.Meta{Project: f.config.Project, Token: f.config.Token}
	return f.config.Transport.UploadFile(context.TODO(), m, path, name, reader)
}

func (f *Filestore) DownloadFile(path string, writer io.Writer) error {
	m := &proto.Meta{Project: f.config.Project, Token: f.config.Token}
	return f.config.Transport.DownloadFile(context.TODO(), m, path, writer)
}
