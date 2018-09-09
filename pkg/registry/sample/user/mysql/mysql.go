package mysql

import (
	"fmt"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"
	mysqlregistry "apistack/pkg/registry/generic/registry/mysqls"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/pages"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"
	"gofreezer/pkg/storage/storagebackend"

	"apistack/examples/apiserver/pkg/api"
	"apistack/examples/apiserver/pkg/registry/core/user"
)

type REST struct {
	*mysqlregistry.Store
}

// NewREST returns a RESTStorage object that will work with testtype.
func NewREST(opts generic.RESTOptions) *REST {
	prefix := "/" + opts.ResourcePrefix
	newListFunc := func() runtime.Object { return &api.UserList{} }

	storageConfig := opts.StorageConfig
	storageConfig.Type = storagebackend.StorageTypeMysql
	storageInterface, dFunc := generic.NewRawStorage(storageConfig)

	store := &registry.Store{
		NewFunc: func() runtime.Object { return &api.User{} },
		// NewListFunc returns an object capable of storing results of an etcd list.
		NewListFunc: newListFunc,
		// Produces a path that etcd understands, to the root of the resource
		// by combining the namespace in the context with the given prefix.
		KeyRootFunc: func(ctx freezerapi.Context) string {
			return prefix
		},
		// Produces a path that etcd understands, to the resource by combining
		// the namespace in the context with the given prefix.
		KeyFunc: func(ctx freezerapi.Context, name string) (string, error) {
			return registry.NoNamespaceKeyFunc(ctx, prefix, name)
		},
		// Retrieve the name field of the resource.
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.User).Name, nil
		},
		// Used to match objects based on labels/fields for list.
		PredicateFunc: func(label labels.Selector, field fields.Selector, page pages.Selector) storage.SelectionPredicate {
			return storage.SelectionPredicate{
				Label: label,
				Field: field,
				Page:  page,
				GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
					user, ok := obj.(*api.User)
					if !ok {
						return nil, nil, fmt.Errorf("unexpected type of given object")
					}
					return labels.Set(user.ObjectMeta.Labels), fields.Set{}, nil
				},
			}
		},

		CreateStrategy:      user.Strategy,
		UpdateStrategy:      user.Strategy,
		DeleteStrategy:      user.Strategy,
		ReturnDeletedObject: true,

		Storage:     storageInterface,
		DestroyFunc: dFunc,
	}

	return &REST{mysqlregistry.NewStore(*store)}
}
