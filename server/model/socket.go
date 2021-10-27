package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SocketType struct {
	Guid         string
	Addr         string
	Status       string
	Wear         int32
	TakeMedicine int32
}

var info = DBType{
	Database:   "iot",
	Collection: "watches",
}
var db = Collection(&info)

func WatchSocketStatus(s *SocketType) bool {
	var (
		filter bson.M
		doc    bson.M
	)
	if s.Status == "on" {
		filter = bson.M{"sn": s.Guid}
		doc = bson.M{"wifi": 1, "addr": s.Addr, "wifiState": s.Status}
	} else {
		filter = bson.M{"addr": s.Addr}
		doc = bson.M{"wifiState": "off", "addr": nil}
	}

	update := bson.D{{"$set", doc}}
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	var result bson.M
	err := db.FindOneAndUpdate(context.TODO(), filter, update, &opt).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

func WatchUpdate(s *SocketType) bool {
	filter := bson.M{"sn": s.Guid}
	update := bson.D{{"$set", bson.M{"sn": s.Guid, "wear": s.Wear, "takeMedicine": s.TakeMedicine}}}
	var result bson.M
	err := db.FindOneAndUpdate(context.TODO(), filter, update).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	return true
}
