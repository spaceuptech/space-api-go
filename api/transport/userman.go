package transport

import (
	"context"

	"github.com/spaceuptech/space-api-go/api/model"
	"github.com/spaceuptech/space-api-go/api/proto"
)

// Profile triggers the gRPC profile function on space cloud
func (t *Transport) Profile(ctx context.Context, meta *proto.Meta, id string) (*model.Response, error) {
	req := proto.ProfileRequest{Id: id, Meta: meta}
	res, err := t.stub.Profile(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 || res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// Profiles triggers the gRPC profiles function on space cloud
func (t *Transport) Profiles(ctx context.Context, meta *proto.Meta) (*model.Response, error) {
	req := proto.ProfilesRequest{Meta: meta}
	res, err := t.stub.Profiles(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 || res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// SignIn triggers the gRPC signIn function on space cloud
func (t *Transport) SignIn(ctx context.Context, meta *proto.Meta, email, password string) (*model.Response, error) {
	req := proto.SignInRequest{Email: email, Password: password, Meta: meta}
	res, err := t.stub.SignIn(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 || res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// SignUp triggers the gRPC signUp function on space cloud
func (t *Transport) SignUp(ctx context.Context, meta *proto.Meta, email, name, password, role string) (*model.Response, error) {
	req := proto.SignUpRequest{Email: email, Name: name, Password: password, Role: role, Meta: meta}
	res, err := t.stub.SignUp(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 || res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}

// EditProfile triggers the gRPC editProfile function on space cloud
func (t *Transport) EditProfile(ctx context.Context, meta *proto.Meta, id string, values model.ProfileParams) (*model.Response, error) {
	req := proto.EditProfileRequest{Id: id, Meta: meta}
	if values.Name != "" {
		req.Name = values.Name
	}
	if values.Email != "" {
		req.Email = values.Email
	}
	if values.Password != "" {
		req.Password = values.Password
	}
	res, err := t.stub.EditProfile(ctx, &req)
	if err != nil {
		return nil, err
	}

	if res.Status >= 200 || res.Status < 300 {
		return &model.Response{Status: int(res.Status), Data: res.Result}, nil
	}

	return &model.Response{Status: int(res.Status), Error: res.Error}, nil
}