package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type UrlsTypes struct {
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}

// 프로젝트 구함
func urls(guid string) string {
	filter := bson.M{
		"sn": guid,
	}
	var result WatchesType
	urlDB.FindOne(context.TODO(), filter).Decode(&result)
	return result.Project
}
