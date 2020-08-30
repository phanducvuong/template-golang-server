package hello

import (
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"encoding/json"
	"net/http"
)

type Person struct {
	Name 	string
	Age		int
}

func GetProfile(w http.ResponseWriter, r *http.Request) {

	person := Person{"Alex", 13}
	js, err := json.Marshal(person)
	if err != nil {
		w.Write([]byte("can't not send profile"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func GetDB(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	// var result Person
	findResult := db.Collection("test").FindOne(context.TODO(), bson.M{"name": "abc"})

	if findResult.Err() != nil {
		w.Write([]byte("No document\n"))
		return
	}

	result := Person{}
	err := findResult.Decode(&result)
	if err != nil {
		w.Write([]byte("no document 1\n"))
		return
	}

	js, err := json.Marshal(result)
	if err != nil {
		w.Write([]byte("parse json failed"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}