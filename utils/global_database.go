package utils

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

var GlobalMongoClient *mongo.Client
var GlobalDatabase *mongo.Database

type Mongo struct {
	Host     string
	UserName string
	Password string
	DbName   string
	AuthDB   string
	MaxConns uint64
	Port     int
	Key      string
}

func init() {
	mongo := Mongo{
		Host:     "",
		UserName: "",
		Password: "",
		DbName:   "",
		AuthDB:   "",
		MaxConns: 0,
		Port:     0,
		Key:      "",
	}
	if _, err := newClientFromConf(&mongo); err != nil {
		panic(fmt.Sprintf("[MONGODB] NEW CONNECT FAIL: %s \n", err))
	}
}

func newClientFromConf(mongoConf *Mongo) (*mongo.Client, error) {
	pwd := mongoConf.Password
	if mongoConf.MaxConns == 0 {
		mongoConf.MaxConns = 100
	}
	Readconcern := false
	Writeconcern := false

	opts := &options.ClientOptions{}

	opts.SetMaxPoolSize(mongoConf.MaxConns)
	opts.ApplyURI(mongoConf.Host)

	if Readconcern == true {
		opts.SetReadConcern(readconcern.Majority())
	}

	if Writeconcern == true {
		wc := writeconcern.New(writeconcern.WMajority())
		opts.SetWriteConcern(wc)
	}
	if mongoConf.UserName != "" {
		credential := options.Credential{
			Username:   mongoConf.UserName,
			Password:   pwd,
			AuthSource: mongoConf.AuthDB,
		}
		opts.SetAuth(credential)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
	defer cancel()
	mongoCli, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	GlobalMongoClient = mongoCli
	GlobalDatabase = GlobalMongoClient.Database(mongoConf.DbName)

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = mongoCli.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	return GlobalMongoClient, nil
}
