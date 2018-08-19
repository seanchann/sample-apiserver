/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package storage

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/seanchann/sample-server/pkg/apis/sample"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

// REST implements  for Test
type REST struct {
}

//NewREST new a rest api for Test
func NewREST() *REST {
	return &REST{}
}

//New creates a new Test request
func (r *REST) New() runtime.Object {
	return &api.Test{}
}

//Create send a request to asterisk
func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (out runtime.Object, err error) {
	_, ok := obj.(*api.Test)
	if !ok {
		return nil, fmt.Errorf("invalid Test object: %#v", obj)
	}

	return &metav1.Status{
		Status:  metav1.StatusSuccess,
		Message: fmt.Sprintf("Test request succeeded"),
		Code:    http.StatusOK,
	}, nil
}

//NamespaceScoped implement scoper interface.
func (r *REST) NamespaceScoped() bool {
	return false
}
