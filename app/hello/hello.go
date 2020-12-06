package hello

import (
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"encoding/json"
	"net/http"
)

type Person struct {
	Name 	string 	`json:"name"`
	Age		int64		`json:"age"`
	Ls		[]int		`json:"ls"`
}

func GetProfile(w http.ResponseWriter, r *http.Request) {

	person := Person{"Alex", 13, []int{1,2,3}}
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

func TestPostData(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p.Name)
	fmt.Println(p.Age)
	fmt.Println(p.Ls)

	w.Write([]byte("ok"))
}

func GetAllData(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	cursor, err := db.Collection("test").Find(context.TODO(), bson.M{})
	if err != nil {
		w.Write([]byte("find value failed!"))
		return
	}

	var arrPerson []Person
	for cursor.Next(context.TODO()) {
		var elem Person
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
			w.Write([]byte("get document failed!"))
			return
		}

		arrPerson = append(arrPerson, elem)
	}

	if cursor.Err() != nil {
		w.Write([]byte("Get document failed!"))
		return
	}

	cursor.Close(context.TODO())
	jsArr, err := json.Marshal(arrPerson)
	if err != nil {
		w.Write([]byte("parse json failed!"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsArr)
}