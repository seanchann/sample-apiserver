/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package mysql

import (
	"context"
	"database/sql"
	stderrs "errors"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/golang/glog"
)

func ScanRows(rows *sql.Rows, t *Table, obj runtime.Object) ([]*RowResult, error) {
	columns, _ := rows.Columns()
	count := len(columns)
	defaultVals := make([]interface{}, count)
	valuesPtrs := make([]interface{}, count)

	var listObj []*RowResult
	tableObj := reflect.Indirect(reflect.New(t.obj.Type()))
	for i, col := range columns {
		key, ok := t.columnToFreezerTagKey[col]
		if !ok {
			//use default value for scan
			valuesPtrs[i] = &defaultVals[i]
			continue
		}

		filedName := t.freezerTag[key].structField
		valuesPtrs[i] = tableObj.FieldByName(filedName).Addr().Interface()
	}
	for rows.Next() {
		err := rows.Scan(valuesPtrs...)

		if err != nil {
			glog.Errorf("scan table(%v) error %v\r\n", t.name, err)
			return nil, err
		}

		item := &RowResult{}
		err = t.CovertRowsToObject(item, obj, tableObj)
		if err != nil {
			glog.Errorf("scan table(%v) error %v\r\n", t.name, err)
			return nil, err
		}
		listObj = append(listObj, item)
	}

	return listObj, nil
}

func GetActualResourceKey(key string) string {
	var actual string
	if i := strings.LastIndexAny(key, "/"); i >= 0 {
		actual = key[i+1:]
	} else {
		actual = key
	}
	return actual
}

// WithTable returns a copy of parent in which the value associated with tablecontextKey is val.
func WithTable(parent context.Context, val interface{}) context.Context {
	internalCtx, ok := parent.(context.Context)
	if !ok {
		panic(stderrs.New("Invalid context type"))
	}
	return context.WithValue(internalCtx, tablecontextKey, val)
}

// UpdateNameWithResouceKey implements metadata.name
func UpdateNameWithResouceKey(obj runtime.Object, name string) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	accessor.SetName(name)
	return nil
}

func GetObjKind(objPtr runtime.Object) string {
	v, err := conversion.EnforcePtr(objPtr)
	if err != nil {
		return string("")
	}

	kind := v.Type().String()
	if i := strings.IndexAny(kind, "."); i >= 0 {
		kind = kind[i+1:]
	}
	return kind
}

func GetListItemObj(listObj runtime.Object) (listPtr interface{}, itemObj runtime.Object, err error) {
	listPtr, err = meta.GetItemsPtr(listObj)
	if err != nil {
		return
	}

	items, err := conversion.EnforcePtr(listPtr)
	if err != nil {
		return
	}
	if items.Kind() != reflect.Slice {
		err = fmt.Errorf("object(%v) not a slice", items.Kind())
		return
	}

	itemObj = reflect.New(items.Type().Elem()).Interface().(runtime.Object)

	return
}

func CloneRuntimeObj(objPtr runtime.Object) (runtime.Object, error) {
	v, err := conversion.EnforcePtr(objPtr)
	if err != nil {
		return nil, err
	}

	newObj := reflect.New(v.Type())
	storeObj := newObj.Interface().(runtime.Object)
	return storeObj, nil
}
