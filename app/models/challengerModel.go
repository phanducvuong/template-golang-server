package models

type ChallengerModel struct {
	HighScore			int64			`json:"high_score" bson:"high_score"`
}