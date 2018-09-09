package user

import (
	"fmt"

	"gofreezer/pkg/api"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/pages"
	"gofreezer/pkg/runtime"
	apistorage "gofreezer/pkg/storage"
	"gofreezer/pkg/util/validation/field"

	"apistack/pkg/registry/generic"

	userapi "apistack/examples/apiserver/pkg/api"
	"apistack/examples/apiserver/pkg/api/validation"
)

// storageClassStrategy implements behavior for StorageClass objects
type userStrategy struct {
	runtime.ObjectTyper
	api.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = userStrategy{api.Scheme, api.SimpleNameGenerator}

func (userStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (userStrategy) PrepareForCreate(ctx api.Context, obj runtime.Object) {
	_ = obj.(*userapi.User)
}

func (userStrategy) Validate(ctx api.Context, obj runtime.Object) field.ErrorList {
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
func (userStrategy) PrepareForUpdate(ctx api.Context, obj, old runtime.Object) {
	_ = obj.(*userapi.User)
	_ = old.(*userapi.User)
	// newuser := obj.(*userapi.User)
	// olduser := old.(*userapi.User)
	// olduser.Name = newuser.Name
	// glog.Infof("got newuser : %+v \r\n", *newuser)
	// glog.Infof("got olduser : %+v \r\n", *olduser)
}

func (userStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateUser(obj.(*userapi.User))
	return append(errorList, validation.ValidateUserUpdate(obj.(*userapi.User), old.(*userapi.User))...)
}

func (userStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchUser returns a generic matcher for a given label and field selector.
func MatchUser(label labels.Selector, field fields.Selector, page pages.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label: label,
		Field: field,
		Page:  page,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			cls, ok := obj.(*userapi.User)
			if !ok {
				return nil, nil, fmt.Errorf("given object is not of type TestType")
			}

			return labels.Set(cls.ObjectMeta.Labels), StorageClassToSelectableFields(cls), nil
		},
	}
}

// StorageClassToSelectableFields returns a label set that represents the object
func StorageClassToSelectableFields(user *userapi.User) fields.Set {
	return generic.ObjectMetaFieldsSet(&user.ObjectMeta, false)
}
