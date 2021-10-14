package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileInfo struct {
	Guid         string   `json:"guid" bson:"guid"`
	FileName     string   `json:"filename" bson:"filename"`
	VideoMD5     string   `json:"videomd5" bson:"videomd5"`
	SubFile      []string `json:"subfile" bson:"subfile"`
	SubFileCount int32    `json:"subfilecount" bson:"subfilecount"`
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
