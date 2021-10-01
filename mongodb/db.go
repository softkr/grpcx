package db

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileInfo struct {
	Guid         string   `json:"guid" bson:"guid"`
	FileName     string   `json:"filename" bson:"filename"`
	VideoMD5     string   `json:"videomd5" bson:"videomd5"`
	SubFile      []string `json:"subfile" bson:"subfile"`
	SubFileCount int32    `json:"subfilecount" bson:"subfilecount"`
}

type WatchesType struct {
	Project      string `json:"project" bson:"project"`
	Sn           string `json:"sn" bson:"sn"`
	TakeMedicine int32  `json:"takeMedicine" bson:"takeMedicine"`
	Wear         int32  `json:"wear" bson:"wear"`
	Addr         string `json:"addr" bson:"addr"`
	WifiState    string `json:"wifiState" bson:"wifiState"`
	Old_sn       string `json:"old_sn" bson:"old_sn"`
	Wifi         string `json:"wifi" bson:"wifi"`
}

const (
	iot     = "iot"
	MFile   = "mfiles"
	Watches = "watches"
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

func Collection() *mongo.Collection {
	var (
		client     = GetMgoCli()
		collection *mongo.Collection
	)
	collection = client.Database(iot).Collection(MFile)
	collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "filename", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	return collection
}

func WatchCollection() *mongo.Collection {
	var (
		client     = GetMgoCli()
		collection *mongo.Collection
	)
	collection = client.Database(iot).Collection(Watches)
	return collection
}

var (
	collection      = Collection()
	watchCollection = WatchCollection()
)

// 프로젝트 구함
func GetProject(guid string) string {
	filter := bson.M{
		"sn": guid,
	}
	var result WatchesType
	watchCollection.FindOne(context.TODO(), filter).Decode(&result)
	return result.Project
}

func Insert(data *FileInfo) {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"guid": data.FileName,
	}
	query := bson.M{
		"$setOnInsert": bson.M{
			"guid":     data.Guid,
			"filename": data.FileName,
			"videomd5": data.VideoMD5,
			"subfile":  data.SubFile,
		},
	}
	collection.UpdateOne(context.TODO(), filter, query, opts)
	// if err != nil {
	// 	log.Print(err)
	// }
}

func Update(value string) int32 {
	filter := bson.M{
		"subfile": value,
	}
	var result FileInfo
	collection.FindOneAndUpdate(context.TODO(), filter, bson.M{
		"$inc": bson.M{
			"subfilecount": 1,
		},
	}).Decode(&result)
	return result.SubFileCount
}

func Find(value string) (data FileInfo) {
	filter := bson.M{
		"subfile": value,
	}
	var result FileInfo
	collection.FindOne(context.TODO(), filter).Decode(&result)
	return result
}

func DeleteOne(value string) {
	filter := bson.M{
		"videomd5": value,
	}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

func Count(value string) int32 {
	filter := bson.M{
		"videomd5": value,
	}
	var result FileInfo
	collection.FindOne(context.TODO(), filter).Decode(&result)
	return result.SubFileCount - int32(len(result.SubFile))
}
