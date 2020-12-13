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

type updateScoreModel struct {
	FBId      string `json:"fb_id"`
	Time      int64  `json:"time"`
	HighScore int64  `json:"high_score"`
	Combo     int    `json:"combo"`
	BestCombo int    `json:"best_combo"`
	IDLevel   int    `json:"id_level"`
}

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
		rankModel := models.RankModel{
			FBId: userModel.FBId,
			Name: userModel.Name,
		}

		_, errInsertRank := db.Collection("ranks").InsertOne(context.TODO(), rankModel)
		_, errInsertUser := db.Collection("users").InsertOne(context.TODO(), userModel)
		if errInsertUser != nil || errInsertRank != nil {
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
	var updateScoreModel updateScoreModel

	err := json.NewDecoder(r.Body).Decode(&updateScoreModel)
	if err != nil || updateScoreModel.FBId == "" {
		fmt.Println(err)
		w.Write(util.ResponseUtil(3000, "can not update score"))
		return
	}

	result := db.Collection("ranks").FindOne(context.TODO(), bson.M{"fb_id": updateScoreModel.FBId})
	if result.Err() != nil {
		w.Write(util.ResponseUtil(3000, "user not exist "+updateScoreModel.FBId))
		return
	} //user no have data in rankModel

	rankModel := models.RankModel{}
	err = result.Decode(&rankModel)
	if err != nil {
		errDecode := util.ResponseUtil(3000, "Decode json failed!")
		w.Write(errDecode)
		return
	}

	functions.UpdateScoreUser(&rankModel, updateScoreModel.IDLevel, updateScoreModel.Time,
		updateScoreModel.HighScore, updateScoreModel.Combo, updateScoreModel.BestCombo)

	filter := bson.M{"fb_id": bson.M{"$eq": updateScoreModel.FBId}}
	update := bson.M{"$set": bson.M{"data": rankModel.Data}}
	_, err = db.Collection("ranks").UpdateOne(context.TODO(), filter, update)

	if err != nil {
		w.Write(util.ResponseUtil(3000, "can not update score user!"))
		return
	}

	w.Write(util.ResponseUtil(2000, "update score success!"))
}

func GetLeaderboardByLevel(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	level, err := strconv.Atoi(query.Get("level"))

	if err != nil || level < 0 {
		w.Write(util.ResponseUtil(3000, "Invalid level!"))
		return
	}

	cursor, err := db.Collection("ranks").Find(context.TODO(), bson.M{
		"data": bson.M{
			"$elemMatch": bson.M{"id_level": level},
		}})
	if err != nil {
		w.Write(util.ResponseUtil(3000, "error"))
		return
	}

	var arrRank []functions.Leaderboard
	for cursor.Next(context.TODO()) {
		var elem models.RankModel
		err = cursor.Decode(&elem)

		arrRank = append(arrRank, functions.GetLeaderboard(elem, level))
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

	if err != nil || level < 0 {
		w.Write(util.ResponseUtil(3000, "Invalid level!"))
		return
	}

	cursor, err := db.Collection("ranks").Find(context.TODO(), bson.M{
		"data": bson.M{
			"$elemMatch": bson.M{"id_level": level},
		},
	}, options.Find().SetSort(bson.M{"high_score": -1}).SetLimit(3))

	if err != nil {
		w.Write(util.ResponseUtil(3000, "error"))
		return
	}

	var arrRank []functions.Leaderboard
	for cursor.Next(context.TODO()) {
		var elem models.RankModel
		err = cursor.Decode(&elem)

		arrRank = append(arrRank, functions.GetLeaderboard(elem, level))
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
