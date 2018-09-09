package rest

import (
	"apistack/examples/apiserver/pkg/api"
	"apistack/examples/apiserver/pkg/registry/core/user/dynamodb"
	"apistack/examples/apiserver/pkg/registry/core/user/etcd"
	"apistack/examples/apiserver/pkg/registry/core/user/mysql"
	"fmt"
	"time"

	"github.com/golang/glog"

	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"

	freezerapi "gofreezer/pkg/api"
)

type UserREST struct {
	etcd     *etcd.REST
	mysql    *mysql.REST
	dynamodb *dynamodb.REST
}

var mysqlbackend = true
var etcdbackend = false
var dynamodbbackend = false

func NewREST(etcdHandler *etcd.REST, mysqlHandler *mysql.REST, dynamodb *dynamodb.REST) *UserREST {
	return &UserREST{
		etcd:     etcdHandler,
		mysql:    mysqlHandler,
		dynamodb: dynamodb,
	}
}

func (*UserREST) New() runtime.Object {
	return &api.User{}
}

func (*UserREST) NewList() runtime.Object {
	return &api.UserList{}
}

func (r *UserREST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {

	if mysqlbackend {
		obj, err := r.mysql.Get(ctx, name)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, err
		}

		// olduser := obj.(*api.User)
		// olduser.Name = olduser.Spec.DetailInfo.Name
		return obj, err
	}

	if etcdbackend {
		etcdobj, err := r.etcd.Get(ctx, name)
		if err != nil {
			glog.Errorf("got from etcd error %v\r\n", err)
			return nil, err
		}
		return etcdobj, err
	}

	if dynamodbbackend {
		obj, err := r.dynamodb.Get(ctx, name)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, err
		}
		user := obj.(*api.User)
		freezerapi.FillObjectMetaSystemFields(ctx, &user.ObjectMeta)
		return obj, err
	}
	return nil, errors.NewInternalError(fmt.Errorf("not enable any backend"))
}

func (r *UserREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	glog.V(5).Infof("list options %+v", *options)
	if mysqlbackend {
		obj, err := r.mysql.List(ctx, options)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, err
		}
		return obj, err
	}

	if etcdbackend {
		etcdobj, err := r.etcd.List(ctx, options)
		if err != nil {
			glog.Errorf("got from etcd error %v\r\n", err)
			return nil, err
		}
		return etcdobj, err
	}

	if dynamodbbackend {
		obj, err := r.dynamodb.List(ctx, options)
		if err != nil {
			glog.Errorf("got from dynamodb error %v\r\n", err)
			return nil, err
		}
		return obj, err
	}

	return nil, errors.NewInternalError(fmt.Errorf("not enable any backend"))
}

func (r *UserREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	if mysqlbackend {
		glog.Infof("**************update mysql db \r\n")
		oldObj, err := r.Get(ctx, name)
		if err != nil {
			return nil, false, err
		}

		olduser := oldObj.(*api.User)
		glog.Infof("**************update mysql db oldobj:%+v \r\n", *olduser)

		newObj, err := objInfo.UpdatedObject(ctx, oldObj)
		if err != nil {
			return nil, false, err
		}
		newuser := newObj.(*api.User)
		glog.Infof("**************newobj:%+v \r\n", *newuser)
		olduser.Name = olduser.Spec.DetailInfo.Name
		newuser.Spec.DetailInfo.ID = olduser.Spec.DetailInfo.ID
		newuser.Name = olduser.Spec.DetailInfo.Name
		newuser.Spec.DetailInfo.RegDBTime = olduser.Spec.DetailInfo.RegDBTime
		newuser.Spec.DetailInfo.LastResetPwdTime = unversioned.NewTime(time.Now())
		newuser.Spec.DetailInfo.LastCheckInTime = unversioned.NewTime(time.Now())

		glog.Infof("**************update mysql db newobj:%+v \r\n", *newuser)

		obj, falg, err := r.mysql.Update(ctx, name, rest.DefaultUpdatedObjectInfo(newuser, api.Scheme))
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, falg, err
		}
		objUser := obj.(*api.User)
		glog.Infof("**************update mysql db resultObj:%+v \r\n", *objUser)
		return objUser, falg, err
	}

	if etcdbackend {
		etcdobj, falg, err := r.etcd.Update(ctx, name, objInfo)
		if err != nil {
			glog.Errorf("got from etcd error %v\r\n", err)
			return nil, falg, err
		}
		return etcdobj, falg, err
	}

	if dynamodbbackend {
		glog.Infof("**************update dynamo db \r\n")
		oldObj, err := r.Get(ctx, name)
		if err != nil {
			return nil, false, err
		}

		olduser := oldObj.(*api.User)
		glog.Infof("**************update dynamo db oldobj:%v \r\n", *olduser)

		newObj, err := objInfo.UpdatedObject(ctx, oldObj)
		if err != nil {
			return nil, false, err
		}
		newuser := newObj.(*api.User)
		newuser.Spec.DetailInfo = olduser.Spec.DetailInfo
		glog.Infof("**************update dynamo db newobj%v\r\n", *newuser)

		obj, flag, err := r.dynamodb.Update(ctx, name, rest.DefaultUpdatedObjectInfo(newuser, api.Scheme))
		if err != nil {
			glog.Errorf("got from dynamodb error %v\r\n", err)
			return nil, flag, err
		}
		return obj, flag, err
	}

	return nil, false, errors.NewInternalError(fmt.Errorf("not enable any backend"))
}

func (r *UserREST) Delete(ctx freezerapi.Context, name string, options *freezerapi.DeleteOptions) (runtime.Object, error) {
	glog.V(5).Infof("delete %v resoure", name)
	if mysqlbackend {
		obj, err := r.mysql.Delete(ctx, name, options)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, err
		}
		return obj, err
	}

	if etcdbackend {
		etcdobj, err := r.etcd.Delete(ctx, name, options)
		if err != nil {
			glog.Errorf("got from etcd error %v\r\n", err)
			return nil, err
		}
		return etcdobj, err
	}

	if dynamodbbackend {
		obj, err := r.dynamodb.Delete(ctx, name, options)
		if err != nil {
			glog.Errorf("got from dynamodb error %v\r\n", err)
			return nil, err
		}
		return obj, err
	}

	return nil, errors.NewInternalError(fmt.Errorf("not enable any backend"))
}

func (r *UserREST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	if mysqlbackend {
		user := obj.(*api.User)
		user.Spec.DetailInfo.LastCheckInTime = unversioned.NewTime(time.Now())
		user.Spec.DetailInfo.LastResetPwdTime = user.Spec.DetailInfo.LastCheckInTime

		resObj, err := r.mysql.Create(ctx, obj)
		if err != nil {
			glog.Errorf("got from mysql error %v\r\n", err)
			return nil, err
		}
		return resObj, err
	}

	if etcdbackend {
		resObj, err := r.etcd.Create(ctx, obj)
		if err != nil {
			glog.Errorf("got from etcd error %v\r\n", err)
			return nil, err
		}
		return resObj, err
	}

	if dynamodbbackend {
		resObj, err := r.dynamodb.Create(ctx, obj)
		if err != nil {
			glog.Errorf("got from dynamodb error %v\r\n", err)
			return nil, err
		}
		return resObj, err
	}

	return nil, errors.NewInternalError(fmt.Errorf("not enable any backend"))
}
