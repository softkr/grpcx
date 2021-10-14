package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type WatchesType struct {
	Project      string `json:"project" bson:"project"`
	Guid         string `json:"guid" bson:"guid"`
	TakeMedicine int32  `json:"takeMedicine" bson:"takeMedicine"`
	Wear         int32  `json:"wear" bson:"wear"`
	Addr         string `json:"addr" bson:"addr"`
	WifiState    string `json:"wifiState" bson:"wifiState"`
	OldSn        string `json:"old_sn" bson:"old_sn"`
	Wifi         string `json:"wifi" bson:"wifi"`
}

type WatchSocketType struct {
	Watch_sn              string `json:"watch_sn" bson:"watch_sn"`
	Step_count            int    `json:"step_count" bson:"step_count"`
	Camera_shooting_count int    `json:"camera_shooting_count" bson:"camera_shooting_count"`
}

// GetProjectDB 프로젝트 구함
func FindWatch(guid string) string {
	filter := bson.M{
		"sn": guid,
	}
	var result WatchesType
	socketDB.FindOne(context.TODO(), filter).Decode(&result)
	return result.Project
}

func WatchUpdate(watch_sn string, step_count, camera_shooting_count int32) string {
	filter := bson.M{"sn": watch_sn}
	update := bson.D{{"$set", bson.M{"sn": watch_sn, "wear": step_count, "takeMedicine": camera_shooting_count}}}
	var result WatchesType
	socketDB.FindOneAndUpdate(context.TODO(), filter, update).Decode(&result)
	return result.Project
}

func WatchSocketStatus(guid, addr, status string) string {
	var (
		filter bson.M
		doc    bson.M
	)
	fmt.Println(status)
	if status == "on" {
		fmt.Println("on")
		filter = bson.M{"sn": guid}
		doc = bson.M{"wifi": 1, "addr": addr, "wifiState": status}
	} else {
		fmt.Println("off")
		filter = bson.M{"addr": addr}
		doc = bson.M{"wifiState": "off", "addr": nil}
	}

	update := bson.D{{"$set", doc}}
	var result WatchesType
	socketDB.FindOneAndUpdate(context.TODO(), filter, update).Decode(&result)
	fmt.Println(result.Project)
	return result.Project
}
