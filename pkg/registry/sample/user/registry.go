/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package user

import (
	"context"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/apiserver/pkg/registry/rest"

	api "github.com/seanchann/sample-apiserver/pkg/apis/sample"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	ListUsers(ctx context.Context, options *metainternalversion.ListOptions) (*api.UserList, error)
	WatchUsers(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error)
	GetUser(ctx context.Context, name string, options *metav1.GetOptions) (*api.User, error)
	CreateUser(ctx context.Context, user *api.User, createValidation rest.ValidateObjectFunc) error
	UpdateUser(ctx context.Context, svc *api.User, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc) error
	DeleteUser(ctx context.Context, name string) error
}

// storage puts strong typing around storage calls
type storage struct {
	rest.StandardStorage
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(s rest.StandardStorage) Registry {
	return &storage{s}
}

func (s *storage) ListUsers(ctx context.Context, options *metainternalversion.ListOptions) (*api.UserList, error) {
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*api.UserList), nil
}

func (s *storage) WatchUsers(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	return s.Watch(ctx, options)
}

func (s *storage) GetUser(ctx context.Context, name string, options *metav1.GetOptions) (*api.User, error) {
	obj, err := s.Get(ctx, name, options)
	if err != nil {
		return nil, err
	}
	return obj.(*api.User), nil
}

func (s *storage) CreateUser(ctx context.Context, user *api.User, createValidation rest.ValidateObjectFunc) error {
	_, err := s.Create(ctx, user, createValidation, nil)
	return err
}

func (s *storage) UpdateUser(ctx context.Context, user *api.User, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc) error {
	_, _, err := s.Update(ctx, user.Name, rest.DefaultUpdatedObjectInfo(user), createValidation, updateValidation, false, nil)
	return err
}

func (s *storage) DeleteUser(ctx context.Context, name string) error {
	_, _, err := s.Delete(ctx, name, nil)
	return err
}
