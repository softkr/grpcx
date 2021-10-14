package db

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
