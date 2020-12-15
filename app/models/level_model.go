package models

type LevelModel struct {
	FbID			string	`json:"fb_id" bson:"fb_id"`
	Name			string	`json:"name"`
	Time			int64		`json:"time"`
	HighScore	int64		`json:"high_score" bson:"high_score"`
	Combo			int			`json:"combo"`
	BestCombo	int			`json:"best_combo" bson:"best_combo"`
	IDLevel		int			`json:"id_level" bson:"id_level"`
}