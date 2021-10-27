package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var fileDB = FileCollection()

// FindWatchType
type FindWatchType struct {
	Project string `bson:"project"`
}

// GetProjectDB 프로젝트 구함
func FindWatch(guid string) string {
	filter := bson.M{
		"sn": guid,
	}
	var result FindWatchType
	db.FindOne(context.Background(), filter).Decode(&result)
	return result.Project
}

type FileInfoType struct {
	Guid         string   `json:"guid" bson:"guid"`
	FileName     string   `json:"filename" bson:"filename"`
	VideoMD5     string   `json:"videomd5" bson:"videomd5"`
	SubFile      []string `json:"subfile" bson:"subfile"`
	SubFileCount int32    `json:"subfilecount" bson:"subfilecount"`
}

func FileInsert(data *FileInfoType) {
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
	fileDB.UpdateOne(context.TODO(), filter, query, opts)
}

func FileUpdate(value string) int32 {
	filter := bson.M{
		"subfile": value,
	}
	var result FileInfoType
	fileDB.FindOneAndUpdate(context.TODO(), filter, bson.M{
		"$inc": bson.M{
			"subfilecount": 1,
		},
	}).Decode(&result)
	return result.SubFileCount
}

func FileFind(value string) (data FileInfoType) {
	filter := bson.M{
		"subfile": value,
	}
	var result FileInfoType
	fileDB.FindOne(context.TODO(), filter).Decode(&result)
	return result
}

func FileDeleteOne(value string) {
	filter := bson.M{
		"videomd5": value,
	}
	_, err := fileDB.DeleteMany(context.TODO(), filter)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

func FileCount(value string) int32 {
	filter := bson.M{
		"videomd5": value,
	}
	var result FileInfoType
	fileDB.FindOne(context.TODO(), filter).Decode(&result)
	return result.SubFileCount - int32(len(result.SubFile))
}
