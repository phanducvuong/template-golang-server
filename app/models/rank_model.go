package models

type RankModel struct {
	FBId			string				`json:"fb_id" bson:"fb_id"`
	Name			string				`json:"name" bson:"name"`
	Data			[]LevelModel	`json:"data" bson:"data"`
}