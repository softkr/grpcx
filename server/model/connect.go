package model

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type DBType struct {
	Database   string
	Collection string
}

const (
	IOT   = "iot"
	MFile = "mfiles"
)

var mgoCli *mongo.Client

func initEngine() {
	var err error

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))

	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
	}
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Println(err)
	}
}

func GetMgoCli() *mongo.Client {
	if mgoCli == nil {
		initEngine()
	}
	return mgoCli
}

func Collection(t *DBType) *mongo.Collection {
	var (
		client     = GetMgoCli()
		collection *mongo.Collection
	)
	collection = client.Database(IOT).Collection(t.Collection)
	collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "filename", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	return collection
}

func FileCollection() *mongo.Collection {
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
