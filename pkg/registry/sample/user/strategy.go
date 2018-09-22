/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package user

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	apistorage "k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"

	"github.com/seanchann/apimaster/pkg/api/legacyscheme"

	userapi "github.com/seanchann/sample-apiserver/pkg/apis/sample"
	"github.com/seanchann/sample-apiserver/pkg/apis/sample/validation"
)

// storageClassStrategy implements behavior for StorageClass objects
type userStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = userStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

func (userStrategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

func (userStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (userStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	_ = obj.(*userapi.User)
}

func (userStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	user := obj.(*userapi.User)
	return validation.ValidateUser(user)
}

// Canonicalize normalizes the object after validation.
func (userStrategy) Canonicalize(obj runtime.Object) {
}

func (userStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (userStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	_ = obj.(*userapi.User)
	_ = old.(*userapi.User)
}

func (userStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateUser(obj.(*userapi.User))
	return append(errorList, validation.ValidateUserUpdate(obj.(*userapi.User), old.(*userapi.User))...)
}

func (userStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// userToSelectableFields returns a field set that represents the object
func userToSelectableFields(user *userapi.User) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&user.ObjectMeta, true)
	specificFieldsSet := fields.Set{}
	return generic.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, bool, error) {
	cls, ok := obj.(*userapi.User)
	if !ok {
		return nil, nil, false, fmt.Errorf("given object is not of type TestType")
	}

	return labels.Set(cls.ObjectMeta.Labels), userToSelectableFields(cls), false, nil
}

// MatchUser returns a generic matcher for a given label and field selector.
func MatchUser(label labels.Selector, field fields.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}
