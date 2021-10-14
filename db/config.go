package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	IOT     = "iot"    // 디비 계정
	MFile   = "mfiles" // db 명
	Watch   = "watch"  // db 명
	Project = "project"
	Url     = "urls"
)

var (
	client    = GetMgoCli()
	iotDB     = client.Database(IOT).Collection(MFile)
	socketDB  = client.Database(IOT).Collection(Watch)
	projectDB = client.Database(IOT).Collection(Project)
	urlDB     = client.Database(IOT).Collection(Url)
)

var collection = Collection()

func Collection() *mongo.Collection {
	var (
		client     = GetMgoCli()
		collection *mongo.Collection
	)
	collection = client.Database(IOT).Collection(MFile)
	collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "filename", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	return collection
}
