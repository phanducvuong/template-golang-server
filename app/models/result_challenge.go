package models

type ResultChallenge struct {
	Time			int64		`json:"time"`
	HighScore	int64		`json:"high_score" bson:"high_score"`
	Combo			int			`json:"combo"`
	BestCombo	int			`json:"best_combo" bson:"best_combo"`
}