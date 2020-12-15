package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"rank-server-pikachu/app/functions"
	"rank-server-pikachu/app/models"
	"rank-server-pikachu/app/util"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type challengeBody struct {
	FBIdA     string `json:"fb_id_a"`
	FBIdB     string `json:"fb_id_b"`
	ModeScene *int   `json:"mode_scene"`
	NumOfItem *int   `json:"num_of_item"`
	Row       *int   `json:"row"`
	Col       *int   `json:"col"`
	DataBoard string `json:"data_board"`
	Time      *int64 `json:"time"`
	HighScore *int64 `json:"high_score"`
	Combo     *int   `json:"combo"`
	BestCombo *int   `json:"best_combo"`
}

func InitUser(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	var userModel models.UserModel
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil || userModel.FBId == "" || userModel.Name == "" {
		errRes := util.ResponseUtil(3000, "can not get body")
		w.Write(errRes)
		return
	}

	if !functions.ChkUserExist(db, userModel.FBId) {
		_, errInsertUser := db.Collection("users").InsertOne(context.TODO(), userModel)
		if errInsertUser != nil {
			errorInsert := util.ResponseUtil(3000, "can't insert document to db")
			w.Write(errorInsert)
			return
		}

		w.Write(util.ResponseUtil(2000, "success"))
		return
	}

	w.Write(util.ResponseUtil(3000, "User is exist!"))
}

func UpdateScoreUser(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	var updateScoreModel models.LevelModel

	err := json.NewDecoder(r.Body).Decode(&updateScoreModel)
	if err != nil || updateScoreModel.FbID == "" {
		w.Write(util.ResponseUtil(3000, "Can't update score!"))
		return
	}

	if !functions.ChkUserExist(db, updateScoreModel.FbID) {
		w.Write(util.ResponseUtil(3000, "User is not exist!"))
		return
	} //user is not exist

	if !functions.ChkLevelUserExist(db, updateScoreModel.FbID, updateScoreModel.IDLevel) {
		_,err := db.Collection("levels").InsertOne(context.TODO(), updateScoreModel);
		if err != nil {
			w.Write(util.ResponseUtil(3000, "Update level failed!"))
			return
		}
		w.Write(util.ResponseUtil(2000, "Success!"))
		return
	} //level is not exist

	filter 	:= bson.M{"$and": []interface{}{
		bson.M{"id_level"	: bson.M{"$eq": updateScoreModel.IDLevel}},
		bson.M{"fb_id"		: bson.M{"$eq": updateScoreModel.FbID}},
	}}
	update	:= bson.M{"$set": bson.M{
		"fb_id": updateScoreModel.FbID,
		"name": updateScoreModel.Name,
		"time": updateScoreModel.Time,
		"high_score": updateScoreModel.HighScore,
		"combo": updateScoreModel.Combo,
		"best_combo": updateScoreModel.BestCombo,
		"id_level": updateScoreModel.IDLevel,
	}}

	_, err = db.Collection("levels").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.Write(util.ResponseUtil(3000, "Update level failed!"))
		return
	}
	w.Write(util.ResponseUtil(2000, "Update score success!"))
}

func GetLeaderboardByLevel(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	query 			:= r.URL.Query()
	level, err 	:= strconv.Atoi(query.Get("level"))

	if err != nil || level < 0 {
		w.Write(util.ResponseUtil(3000, "Invalid level!"))
		return
	}

	filter 	:= bson.M{"id_level": level}
	ops			:= options.Find().SetSort(bson.M{"high_score": -1}).SetLimit(100)
	cursor, err := db.Collection("levels").Find(context.TODO(), filter, ops)
	if err != nil {
		w.Write(util.ResponseUtil(3000, "error"))
		return
	}

	var arrRank []models.LevelModel
	for cursor.Next(context.TODO()) {
		var elem models.LevelModel
		err = cursor.Decode(&elem)

		arrRank = append(arrRank, elem)
	}

	if cursor.Err() != nil {
		w.Write(util.ResponseUtil(3000, "Get data failed!"))
		return
	}

	cursor.Close(context.TODO())
	jsArr, err := json.Marshal(arrRank)
	if err != nil {
		w.Write(util.ResponseUtil(3000, "Parse json failed!"))
		return
	}

	w.Write(util.ResponseUtil(2000, string(jsArr)))
}

func GetLDBTop3(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	query 			:= r.URL.Query()
	level, err 	:= strconv.Atoi(query.Get("level"))
	fbID				:= query.Get("fb_id")

	if err != nil || fbID == "" || level < 0 {
		w.Write(util.ResponseUtil(3000, "Invalid level!"))
		return
	}

	filter 			:= bson.M{"id_level": level}
	opts				:= options.Find().SetSort(bson.M{"high_score": -1}).SetLimit(3)
	cursor, err := db.Collection("levels").Find(context.TODO(), filter, opts)

	//find high score specific user
	var userLevel = new(models.LevelModel)
	filter			 = bson.M{"$and": []interface{}{
		bson.M{"fb_id": fbID},
		bson.M{"id_level": level},
	}}
	userFind 		:= db.Collection("levels").FindOne(context.TODO(), filter)
	if userFind.Err() == nil {
		userFind.Decode(&userLevel)
	}

	if err != nil {
		w.Write(util.ResponseUtil(3000, "error"))
		return
	}

	// var arrRank []models.LevelModel
	var arrRank []models.LevelModel
	for cursor.Next(context.TODO()) {
		var elem models.LevelModel
		err = cursor.Decode(&elem)
		arrRank = append(arrRank, elem)
	}
	arrRank = append(arrRank, *userLevel)
	defer cursor.Close(context.TODO())

	if cursor.Err() != nil {
		w.Write(util.ResponseUtil(3000, "Get data failed!"))
		return
	}

	jsArr, err := json.Marshal(arrRank)
	if err != nil {
		w.Write(util.ResponseUtil(3000, "Parse json failed!"))
		return
	}

	w.Write(util.ResponseUtil(2000, string(jsArr)))
}

func InitChallenge(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	var challengeBody challengeBody
	err := json.NewDecoder(r.Body).Decode(&challengeBody)
	if err != nil ||
		challengeBody.FBIdA == "" || challengeBody.FBIdB == "" || challengeBody.DataBoard == "" ||
		challengeBody.ModeScene == nil || challengeBody.NumOfItem == nil || challengeBody.Row == nil || challengeBody.Col == nil ||
		challengeBody.Time == nil || challengeBody.HighScore == nil || challengeBody.Combo == nil || challengeBody.BestCombo == nil {

		w.Write(util.ResponseUtil(3000, "Check info challenge!"))
		return

	}

	chkUserExist, dataUser := functions.FindUserByFBId(db, challengeBody.FBIdB)
	if !functions.ChkUserExist(db, challengeBody.FBIdA) || !chkUserExist {
		w.Write(util.ResponseUtil(3000, "FBId User Not Exist!"))
		return
	}

	timee := time.Now().UnixNano() / 1e6
	keyChallenge := fmt.Sprintf("%s_%s_%d", challengeBody.FBIdA, challengeBody.FBIdB, timee)
	dataChallenge := models.ChallengeModel{
		FBIdChallenger: challengeBody.FBIdA,
		FBIdChallenged: challengeBody.FBIdB,
		KeyChallenge:   keyChallenge,
		ModeScene:      *challengeBody.ModeScene,
		NumberItem:     *challengeBody.NumOfItem,
		Board: models.BoardModel{
			Row:  *challengeBody.Row,
			Col:  *challengeBody.Col,
			Data: challengeBody.DataBoard,
		},
		ResultChallenger: models.ResultChallenge{
			Time:      *challengeBody.Time,
			HighScore: *challengeBody.HighScore,
			Combo:     *challengeBody.Combo,
			BestCombo: *challengeBody.BestCombo,
		},
		ResultChallenged: models.ResultChallenge{
			Time:      0,
			HighScore: 0,
			Combo:     0,
			BestCombo: 0,
		},
	}

	//update NotiChallenge user challenged
	dataUser.NotiChallenge = append(dataUser.NotiChallenge, keyChallenge)
	filter := bson.M{"fb_id": bson.M{"$eq": challengeBody.FBIdB}}
	update := bson.M{"$set": bson.M{"noti_challenge": dataUser.NotiChallenge}}
	_, err = db.Collection("users").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.Write(util.ResponseUtil(3000, "Update UserModel Failed!"))
		return
	}

	resultInsert := functions.InitChallenge(db, dataChallenge)
	if !resultInsert {
		w.Write(util.ResponseUtil(3000, "Init Challenge Failed!"))
		return
	}
	w.Write(util.ResponseUtil(2000, "Success!"))
}
