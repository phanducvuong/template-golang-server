package controllers

import (
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
	"rank-server-pikachu/app/models"
	"rank-server-pikachu/app/util"
	"rank-server-pikachu/app/functions"
)

type updateScoreModel struct {
	FBId			string						`json:"fb_id"`
	Time			int64							`json:"time"`
	HighScore	int64							`json:"high_score"`
	Combo			int								`json:"combo"`
	BestCombo	int								`json:"best_combo"`
	IDLevel		int								`json:"id_level"`
}

type challengeBody struct {
	FBIdA				string						`json:"fb_id_a"`
	FBIdB				string						`json:"fb_id_b"`
	ModeScene		int32							`json:"mode_scene"`
	NumOfItem		int32							`json:"num_of_item"`
	Row					int16							`json:"row"`
	Col					int16							`json:"col"`
	DataBoard		string						`json:"data_board"`
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
		rankModel	:= models.RankModel {
			FBId			: userModel.FBId,
			Name			: userModel.Name,
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
		errorRes 	:= util.ResponseUtil(3000, "can not update score")
		w.Write(errorRes)
		return
	}

	result := db.Collection("ranks").FindOne(context.TODO(), bson.M{"fb_id": updateScoreModel.FBId})
	if result.Err() != nil {
		w.Write(util.ResponseUtil(3000, "user not exist " + updateScoreModel.FBId))
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

	filter := bson.M{"fb_id": bson.M{ "$eq": updateScoreModel.FBId }}
	update := bson.M{"$set": bson.M{ "data": rankModel.Data }}
	_, err = db.Collection("ranks").UpdateOne(context.TODO(), filter, update)
	
	if err != nil {
		w.Write(util.ResponseUtil(3000, "can not update score user!"))
		return
	}

	w.Write(util.ResponseUtil(2000, "update score success!"))
}

func InitChallenge(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	timee := time.Now().UnixNano() / 1e6
	fmt.Printf("%d\n", timee)
	fmt.Printf("%s\n", time.Now())

	w.Write(util.ResponseUtil(2000, "ccc"))
}