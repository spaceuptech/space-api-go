package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/proto"
	"github.com/spaceuptech/space-api-go/utils"
)

// CreateFolder triggers the gRPC CreateFolder function on space cloud
func (t *Transport) CreateFolder(ctx context.Context, meta *proto.Meta, path, name string) (*model.Response, error) {
	req := proto.CreateFolderRequest{Name: name, Path: path, Meta: meta}
	res, err := t.stub.CreateFolder(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: nil}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// DeleteFile triggers the gRPC DeleteFile function on space cloud
func (t *Transport) DeleteFile(ctx context.Context, meta *proto.Meta, path string) (*model.Response, error) {
	req := proto.DeleteFileRequest{Path: path, Meta: meta}
	res, err := t.stub.DeleteFile(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: nil}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// ListFiles triggers the gRPC ListFiles function on space cloud
func (t *Transport) ListFiles(ctx context.Context, meta *proto.Meta, path string) (*model.Response, error) {
	req := proto.ListFilesRequest{Path: path, Meta: meta}
	res, err := t.stub.ListFiles(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 && res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: nil}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// UploadFile triggers the gRPC UploadFile function on space cloud
func (t *Transport) UploadFile(ctx context.Context, project, path, name string, meta interface{}, reader io.Reader) (*model.Response, error) {
	r, writer := io.Pipe()

	// Create an http request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.generateFileUploadURL(project), r)
	if err != nil {
		return nil, err
	}

	// Create a multipart mwriter
	mwriter := multipart.NewWriter(writer)
	req.Header.Add("Content-Type", mwriter.FormDataContentType())

	// Create an error channel
	errchan := make(chan error)

	go func() {
		defer close(errchan)
		defer utils.CloseTheCloser(writer)

		w, err := mwriter.CreateFormFile("file", name)
		if err != nil {
			errchan <- err
			return
		}

		if written, err := io.Copy(w, reader); err != nil {
			errchan <- fmt.Errorf("error copying %s (%d bytes written): %v", path, written, err)
			return
		}

		_ = mwriter.WriteField("name", name)

		if err := mwriter.Close(); err != nil {
			errchan <- err
			return
		}
	}()

	// Clean query parameters
	if meta == nil {
		meta = map[string]int{}
	}
	metaJSON, _ := json.Marshal(meta)

	if path == "" {
		path = "/"
	}

	// Set the url parameters
	q := req.URL.Query()
	q.Add("meta", string(metaJSON))
	q.Add("path", path)
	q.Add("fileType", "file")
	q.Add("makeAll", "true")
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer utils.CloseTheCloser(res.Body)
	if err := <-errchan; err != nil {
		return nil, err
	}

	// Unmarshal the response
	result := utils.M{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return &model.Response{Status: res.StatusCode, Data: nil}, nil
	}

	return &model.Response{Status: res.StatusCode, Error: result["error"].(string)}, nil
}

// DownloadFile triggers the gRPC DownloadFile function on space cloud
func (t *Transport) DownloadFile(ctx context.Context, meta *proto.Meta, path string, writer io.Writer) error {
	stream, err := t.stub.DownloadFile(ctx, &proto.DownloadFileRequest{Path: path, Meta: meta})

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

func (t *Transport) generateFileUploadURL(project string) string {
	scheme := "http"
	if t.sslEnabled {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s/v1/api/%s/files", scheme, t.addr, project)
}
