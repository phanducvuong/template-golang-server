package models

type LevelModel struct {
	Time			int			`'json:"time"`
	HighScore	int64		`json:"high_score"`
	Combo			int			`json:"combo"`
	BestCombo	int			`json:"best_combo"`
	IDLevel		int			`json:"id_level"`
}