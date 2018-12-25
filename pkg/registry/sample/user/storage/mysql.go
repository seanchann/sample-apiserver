/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package storage

import (
	api "github.com/seanchann/sample-apiserver/pkg/apis/sample"
	"github.com/seanchann/sample-apiserver/pkg/registry/sample/user"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
)

type mysqlREST struct {
	*genericregistry.Store
}

// newMysqlREST returns a RESTStorage object that will work with testtype.
func newMysqlREST(optsGetter generic.RESTOptionsGetter) *mysqlREST {

	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &api.User{} },
		NewListFunc:              func() runtime.Object { return &api.UserList{} },
		DefaultQualifiedResource: api.Resource("users"),

		PredicateFunc:       user.MatchUser,
		CreateStrategy:      user.Strategy,
		UpdateStrategy:      user.Strategy,
		DeleteStrategy:      user.Strategy,
		ReturnDeletedObject: true,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: user.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		panic(err) // TODO: Propagate error up
	}

	return &mysqlREST{store}
}
