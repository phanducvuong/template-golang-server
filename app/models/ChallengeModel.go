package models

type ChallengeModel struct {
	FBIdChallenger			string							`json:"fb_id_challenger" bson:"fb_id_challenger"`
	FBIdChallenged			string							`json:"fb_id_challenged" bson:"fb_id_challenged"`
	KeyChallenge				string							`json:"key_challenge" bson:"key_challenge"`
	ModeScene						int									`json:"mode_scene" bson:"mode_scene"`
	NumberItem					int									`json:"number_item" bson:"number_item"`
	Board								[]BoardModel				`json:"board"`
	ResultChallenger		int64								`json:"high_scorer" bson:"high_scorer"`
	ResultChallenged		int64								`json:"high_scored" bson:"high_scored"`
}