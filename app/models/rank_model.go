package models

type RankModel struct {
	FBId			string				`json:"fb_id"`
	Data			[]LevelModel	`json:"data"`
}