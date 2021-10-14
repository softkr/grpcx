package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type ProjectType struct {
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
	Url  string `json:"url" bson:"url"`
}

// SetProject
func SetProject(project *ProjectType) {

	res, err := projectDB.InsertOne(context.TODO(), project)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

// SelectProject
func SelectProject() {

	cursor, err := projectDB.Find(context.TODO(), bson.D{{"code", "JAPAN"}})
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}

}

// UpdateProject
func UpdateProject() {
	objID, _ := primitive.ObjectIDFromHex("61667c750e0cd4d50ae40821")
	filter := bson.D{{"_id", objID}}
	update := bson.D{{"$set", bson.D{{"name", "일본"}}}}

	result, err := projectDB.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}
	fmt.Println(result)
}

// DeleteProject
func DeleteProject() {
	objID, _ := primitive.ObjectIDFromHex("61667c750e0cd4d50ae40821")
	projectDB.DeleteOne(context.TODO(), bson.D{{"_id", objID}})
}
