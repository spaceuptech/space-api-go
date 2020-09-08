package filestore

import (
	"context"
	"io"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/types"
)

// Filestore contains the values for the filestore instance
type Filestore struct {
	config *config.Config
}

// New initializes the filestore module
func New(config *config.Config) *Filestore {
	return &Filestore{config}
}

// todo implement this
func (f *Filestore) CreateFolder(ctx context.Context, path, name string) (*types.Response, error) {
	return f.config.Transport.CreateFolder(ctx, f.config.Project, path, name, f.config.Token)
}

func (f *Filestore) DeleteFile(ctx context.Context, path string, meta interface{}) (*types.Response, error) {
	return f.config.Transport.DeleteFile(ctx, meta, f.config.Project, path, f.config.Token)
}

func (f *Filestore) DeleteFolder(ctx context.Context, path string, meta interface{}) (*types.Response, error) {
	return f.config.Transport.DeleteFile(ctx, meta, f.config.Project, path, f.config.Token)
}

func (f *Filestore) ListFiles(ctx context.Context, listWhat, path string) (*types.Response, error) {
	return f.config.Transport.List(ctx, f.config.Project, listWhat, path, f.config.Token)
}

func (f *Filestore) UploadFile(ctx context.Context, path, name string, meta interface{}, reader io.Reader) (*types.Response, error) {
	return f.config.Transport.UploadFile(ctx, f.config.Project, path, name, meta, reader, f.config.Token)
}

func (f *Filestore) DownloadFile(ctx context.Context, path string, writer io.Writer) error {
	return f.config.Transport.DownloadFile(ctx, f.config.Project, path, writer, f.config.Token)
}

func (f *Filestore) DoesFileOrFolderExists(ctx context.Context, path string) (*types.Response, error) {
	return f.config.Transport.DoesExists(ctx, f.config.Project, path, f.config.Token)
}
