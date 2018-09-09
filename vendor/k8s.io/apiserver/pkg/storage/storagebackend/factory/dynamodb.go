package factory

import (
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/awsdynamodb/dynamodb"
	"k8s.io/apiserver/pkg/storage/storagebackend"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/glog"
)

//newDynamodbSession create session with config and credentials
func newDynamodbSession(cfg storagebackend.AWSDynamoDBConfig) (*session.Session, error) {

	config := aws.Config{
		Region: aws.String(cfg.Region),
	}

	if len(cfg.AccessID) > 0 && len(cfg.AccessKey) > 0 {
		config.Credentials = credentials.NewStaticCredentials(cfg.AccessID, cfg.AccessKey, cfg.Token)
	} else {
		config.Credentials = credentials.NewSharedCredentials("", "default")
	}

	// Specify profile for config and region for requests
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: config,
	})
	if err != nil {
		glog.Fatalf("create aws session error %v\r\n", err.Error())
		return nil, err
	}

	return sess, nil
}

func newDynamodbStorage(c storagebackend.Config) (storage.Interface, DestroyFunc, error) {

	client, err := newDynamodbSession(c.AWSDynamoDB)
	if err != nil {
		return nil, nil, err
	}

	destroyFunc := func() {
		//do nothing
	}

	return dynamodb.New(client, c.AWSDynamoDB.Table, c.Codec), destroyFunc, nil
}
