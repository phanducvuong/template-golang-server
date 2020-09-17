package controllers

import (
	"go.mongodb.org/mongo-driver/bson"
	"context"
	// "fmt"
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

func InitUser(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	var userModel models.UserModel
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil || userModel.FBId == "" || userModel.Name == "" {
		errRes := util.ResponseUtil(3000, "can not get body")
		w.Write(errRes)
		return
	}

	findUser := db.Collection("users").FindOne(context.TODO(), bson.M{"fb_id": userModel.FBId})
	if findUser.Err() != nil {
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
	}

	w.Write(util.ResponseUtil(2000, "success"))
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