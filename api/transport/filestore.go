package transport

import (
	"context"
	"io"
	"errors"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
	"github.com/spaceuptech/space-api-go/api/utils"
)

// CreateFolder triggers the gRPC CreateFolder function on space cloud
func (t *Transport) CreateFolder(ctx context.Context, meta *proto.Meta, path, name string) (*model.Response, error) {
	req := proto.CreateFolderRequest{Name: name, Path:path, Meta: meta}
	res, err := t.stub.CreateFolder(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// DeleteFile triggers the gRPC DeleteFile function on space cloud
func (t *Transport) DeleteFile(ctx context.Context, meta *proto.Meta, path string) (*model.Response, error) {
	req := proto.DeleteFileRequest{Path:path, Meta: meta}
	res, err := t.stub.DeleteFile(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// ListFiles triggers the gRPC ListFiles function on space cloud
func (t *Transport) ListFiles(ctx context.Context, meta *proto.Meta, path string) (*model.Response, error) {
	req := proto.ListFilesRequest{Path:path, Meta: meta}
	res, err := t.stub.ListFiles(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// UploadFile triggers the gRPC UploadFile function on space cloud
func (t *Transport) UploadFile(ctx context.Context, meta *proto.Meta, path, name string, reader io.Reader) (*model.Response, error) {
	buf := make([]byte, utils.PayloadSize)
	stream, err := t.stub.UploadFile(context.TODO())
	if err != nil {
		return nil, err
	}
	req := proto.UploadFileRequest{Path:path, Name:name, Meta:meta}
	if err = stream.Send(&req); err != nil {
		return nil, err // Ideally EOF should never occur
	}
	for {
		n, err := reader.Read(buf)
        if n > 0 {
			buf = buf[:n]
		}
        if err == io.EOF {
            break
        }
		if err != nil {
			return nil, err
		}
		req = proto.UploadFileRequest{Payload:buf}
		if err1 := stream.Send(&req); err1 != nil {
			if err1 == io.EOF {
				break // Ideally this should never occur
			} 
			return nil, err1
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// DownloadFile triggers the gRPC DownloadFile function on space cloud
func (t *Transport) DownloadFile(ctx context.Context, meta *proto.Meta, path string, writer io.Writer) error {
	stream, err := t.stub.DownloadFile(ctx, &proto.DownloadFileRequest{Path:path, Meta: meta})
	
	if err != nil {
		return err
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		if res.Status < 200 || res.Status > 300 {
			return errors.New(res.Error)
		}
		
		if _, err = writer.Write(res.Payload); err != nil {
			return err
		}
	}
	return nil
}
