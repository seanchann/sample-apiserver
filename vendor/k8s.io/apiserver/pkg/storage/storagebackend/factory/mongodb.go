package factory

import (
	"time"

	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/mongodbs/mongodb"
	"k8s.io/apiserver/pkg/storage/storagebackend"

	"gopkg.in/mgo.v2"

	"github.com/golang/glog"
)

//dial mongo db with admin db, admin user, admin passwd
func newMongoDBClient(cfg storagebackend.MongoExtendConfig) (*mgo.Session, error) {
	// We need this object to establish a session to our MongoDB.
	adminDB := cfg.AdminCred[0]
	adminName := cfg.AdminCred[1]
	adminPwd := cfg.AdminCred[2]
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    cfg.ServerList,
		Timeout:  60 * time.Second,
		Database: adminDB,
		Username: adminName,
		Password: adminPwd,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		glog.Fatalf("CreateSession:%s\n", err)
		return nil, err
	}
	mongoSession.SetMode(mgo.Monotonic, true)

	generalDB := cfg.GeneralCred[0]
	generalUser := cfg.GeneralCred[1]
	generalPwd := cfg.GeneralCred[2]
	dbhandle := mongoSession.DB(generalDB)

	user := &mgo.User{
		Username: generalUser,
		Password: generalPwd,
		Roles:    []mgo.Role{mgo.RoleReadWrite},
	}

	err = dbhandle.UpsertUser(user)
	if err != nil {
		return nil, err
	}

	cred := &mgo.Credential{
		Username: generalUser,
		Password: generalPwd,
		Source:   generalDB,
	}

	err = mongoSession.Login(cred)
	if err != nil {
		return nil, err
	}

	return mongoSession, nil
}

func newMongoStorage(c storagebackend.Config) (storage.Interface, DestroyFunc, error) {

	client, err := newMongoDBClient(c.Mongodb)
	if err != nil {
		return nil, nil, err
	}

	destroyFunc := func() {
		client.Close()
	}

	return mongodb.New(client, c.Mongodb.GeneralCred[0], c.Codec), destroyFunc, nil
}
