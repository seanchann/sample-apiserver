package factory

import (
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/mysqls/mysql"
	"k8s.io/apiserver/pkg/storage/storagebackend"

	_ "github.com/go-sql-driver/mysql"
	dbmysql "github.com/jinzhu/gorm"
)

//connectionStr: user:password@tcp(host:port)/dbname
func newMysqlClient(connectionStr string) (*dbmysql.DB, error) {
	var err error
	connStr := string(connectionStr) + string("?parseTime=True")
	//connStr := string(connectionStr)
	db, err := dbmysql.Open(string("mysql"), connStr)
	if err != nil {
		return nil, err
	}
	db = db.Debug()

	return db, db.DB().Ping()
}

func newMysqlStorage(c storagebackend.Config) (storage.Interface, DestroyFunc, error) {
	endpoints := c.Mysql.ServerList

	client, err := newMysqlClient(endpoints[0])
	if err != nil {
		return nil, nil, err
	}

	destroyFunc := func() {
		client.Close()
	}

	return mysql.New(client, c.Codec, "v1"), destroyFunc, nil
}
