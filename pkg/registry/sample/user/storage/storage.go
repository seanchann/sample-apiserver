package storage

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	api "github.com/seanchann/sample-apiserver/pkg/apis/sample"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
)

//UserREST impl user rest
type UserREST struct {
	etcd     *etcdREST
	mysql    *mysqlREST
	dynamodb *dynamoREST
}

var mysqlbackend = true
var etcdbackend = false
var dynamodbbackend = false

//NewREST new rest
func NewREST(optsGetter generic.RESTOptionsGetter) *UserREST {
	return &UserREST{
		mysql: newMysqlREST(optsGetter),
	}
}

//New impl New
func (r *UserREST) New() runtime.Object {
	if mysqlbackend {
		return r.mysql.New()
	}

	return nil
}

//NewList impl NewList
func (r *UserREST) NewList() runtime.Object {
	if mysqlbackend {
		return r.mysql.NewList()
	}
	return nil
}

//Get impl get
func (r *UserREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {

	if mysqlbackend {
		return r.mysql.Get(ctx, name, options)
	}

	return nil, apierrors.NewInternalError(fmt.Errorf("not enable any backend"))
}

//List impl list
func (r *UserREST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	glog.V(5).Infof("list options %+v", *options)
	if mysqlbackend {
		return r.mysql.List(ctx, options)
	}

	return nil, apierrors.NewInternalError(fmt.Errorf("not enable any backend"))
}

//Update impl update
func (r *UserREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	if mysqlbackend {
		glog.Infof("**************update mysql db \r\n")
		// oldObj, err := r.Get(ctx, name, nil)
		// if err != nil {
		// 	return nil, false, err
		// }

		// olduser := oldObj.(*api.User)
		// glog.Infof("**************update mysql db oldobj:%+v \r\n", *olduser)

		obj, falg, err := r.mysql.Update(ctx, name, objInfo, createValidation, updateValidation, forceAllowCreate, options)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, falg, err
		}
		objUser := obj.(*api.User)
		glog.Infof("**************update mysql db resultObj:%+v \r\n", *objUser)
		return objUser, falg, err
	}

	return nil, false, apierrors.NewInternalError(fmt.Errorf("not enable any backend"))
}

//Delete impl delete
func (r *UserREST) Delete(ctx context.Context, name string, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	glog.V(5).Infof("delete %v resoure", name)
	if mysqlbackend {
		return r.mysql.Delete(ctx, name, options)
	}

	return nil, false, apierrors.NewInternalError(fmt.Errorf("not enable any backend"))
}

//Create imple create
func (r *UserREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	glog.Infof("create user %v", obj)
	if mysqlbackend {
		return r.mysql.Create(ctx, obj, createValidation, options)
	}

	return nil, apierrors.NewInternalError(fmt.Errorf("not enable any backend"))
}

//Watch impl watch
func (r *UserREST) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	if mysqlbackend {
		return r.mysql.Watch(ctx, options)
	}

	return nil, apierrors.NewInternalError(fmt.Errorf("not enable any backend"))
}

//NamespaceScoped implement scoper interface.
func (r *UserREST) NamespaceScoped() bool {
	return false
}
