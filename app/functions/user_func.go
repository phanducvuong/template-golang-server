package functions

import (
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	// "fmt"
	"rank-server-pikachu/app/models"
)

type Leaderboard struct {
	Name			string			`json:"name"`
	Score			int64				`json:"score"`
	HighScore	int64				`json:"high_score"`
}

// func UpdateScoreUser(levelModel *models.LevelModel, fbId string, name string, idLv int, time int64, highScore int64, combo int, bestCombo int) {
// 	for index, value := range rankUser.Data {
// 		if value.IDLevel == idLv {
			
// 			value.Time 			= time
// 			value.HighScore = highScore
// 			value.Combo			= combo
// 			value.BestCombo	= bestCombo

// 			rankUser.Data[index] = value
// 			return
// 		}
// 	}

// 	newLevel := models.LevelModel {
// 		Time			: time,
// 		HighScore	: highScore,
// 		Combo			: combo,
// 		BestCombo	: bestCombo,
// 		IDLevel		: idLv,
// 	}
// 	rankUser.Data = append(rankUser.Data, newLevel)
// }

func ChkUserExist(db *mongo.Database, fbID string) bool {
	findUser := db.Collection("users").FindOne(context.TODO(), bson.M{ "fb_id": fbID });
	if findUser.Err() != nil {
		return false
	}
	return true
}

func ChkLevelUserExist(db *mongo.Database, userID string, idLv int) bool {
	findUser := db.Collection("levels").FindOne(context.TODO(), bson.M{ "$and": []interface{}{
		bson.M{"user_id": userID},
		bson.M{"id_level": idLv},
	}});
	if findUser.Err() != nil {
		return false
	}
	return true
}

func FindUserByFBId(db *mongo.Database, fbID string) (bool, models.UserModel) {
	findUser := db.Collection("users").FindOne(context.TODO(), bson.M{ "fb_id": fbID });
	if findUser.Err() != nil {
		return false, models.UserModel{}
	}
	var userModel models.UserModel
	findUser.Decode(&userModel)
	return true, userModel
}

func InitChallenge(db *mongo.Database, data models.ChallengeModel) bool {
	_, errInsert := db.Collection("challenges").InsertOne(context.TODO(), data)
	if errInsert != nil {
		return false
	}
	return true
}

// func GetLeaderboard(data models.LevelModel, idLevel int) Leaderboard {
// 	var tmp Leaderboard
// 	v := findLevel(data.Data, idLevel)
// 	if v != nil {
// 		tmp.Name 	= data.Name
// 		tmp.Score = v.HighScore
// 	}
// 	return tmp
// }

// func findLevel(data []models.LevelModel, idLevel int) *models.LevelModel {
// 	for _, value := range data {
// 		if value.IDLevel == idLevel {
// 			return &value
// 		}
// 	}
// 	return nil
// }