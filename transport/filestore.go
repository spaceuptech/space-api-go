package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/spaceuptech/space-api-go/model"
	"github.com/spaceuptech/space-api-go/utils"
)

// CreateFolder creates a folder/directory on selected file store
func (t *Transport) CreateFolder(ctx context.Context, project, path, name string) (*model.Response, error) {
	// Create an http request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.generateFileUploadURL(project), nil)
	if err != nil {
		return nil, err
	}

	if path == "" {
		path = "/"
	}

	// Set the url parameters
	q := req.URL.Query()
	q.Add("path", path)
	q.Add("fileType", "dir")
	q.Add("makeAll", "true")
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	// send the http request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer utils.CloseTheCloser(res.Body)

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

// DeleteFile deletes file/directory from selected file store
func (t *Transport) DeleteFile(ctx context.Context, meta interface{}, project, path string) (*model.Response, error) {
	// Clean query parameters
	if meta == nil {
		meta = map[string]int{}
	}
	metaJSON, _ := json.Marshal(meta)

	if path == "" {
		path = "/"
	}

	// Create an http request
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, t.generateFileUploadURL(project)+path, bytes.NewBuffer(metaJSON))
	if err != nil {
		return nil, err
	}

	// send the http request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer utils.CloseTheCloser(res.Body)

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

// List lists all the files/folders or both according to the mode under certain directory
func (t *Transport) List(ctx context.Context, project, mode, path string) (*model.Response, error) {
	if path == "" {
		path = "/"
	}
	// Create an http request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.generateFileUploadURL(project)+path, nil)
	if err != nil {
		return nil, err
	}

	// Set the url parameters
	q := req.URL.Query()
	q.Add("path", path)
	q.Add("op", "list")
	q.Add("mode", mode)
	req.URL.RawQuery = q.Encode()

	// send the http request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer utils.CloseTheCloser(res.Body)

	// Unmarshal the response
	result := utils.M{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return &model.Response{Status: res.StatusCode, Data: result}, nil
	}

	return &model.Response{Status: res.StatusCode, Error: result["error"].(string)}, nil
}

// UploadFile creates a file in selected file store
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

// DownloadFile downloads specified file from selected file store
func (t *Transport) DownloadFile(ctx context.Context, project, path string, writer io.Writer) error {
	if path == "" {
		path = "/"
	}
	// Create an http request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.generateFileUploadURL(project)+path, nil)
	if err != nil {
		return err
	}

	// Set the url parameters
	q := req.URL.Query()
	q.Add("path", path)
	q.Add("op", "dir")
	req.URL.RawQuery = q.Encode()

	// send the http request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer utils.CloseTheCloser(res.Body)

	_, err = io.Copy(writer, res.Body)
	if err != nil {
		return fmt.Errorf("error downloading file unable to write")
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("error downloading file status code - %v", res.StatusCode)
}

func (t *Transport) generateFileUploadURL(project string) string {
	scheme := "http"
	if t.sslEnabled {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s/v1/api/%s/files", scheme, t.addr, project)
}

// DoesExists checks if specified file exists in selected file store
func (t *Transport) DoesExists(ctx context.Context, project, path string) (*model.Response, error) {
	if path == "" {
		path = "/"
	}
	// Create an http request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.generateFileUploadURL(project)+path, nil)
	if err != nil {
		return nil, err
	}

	// Set the url parameters
	q := req.URL.Query()
	q.Add("path", path)
	q.Add("op", "exist")
	req.URL.RawQuery = q.Encode()

	// send the http request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer utils.CloseTheCloser(res.Body)

	// Unmarshal the response
	result := utils.M{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return &model.Response{Status: res.StatusCode, Data: result}, nil
	}

	return &model.Response{Status: res.StatusCode, Error: result["error"].(string)}, nil
}
