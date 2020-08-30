package models

type ChallengeModel struct {
	FBIdChallenger			string							`json:"fb_id_challenger"`
	FBIdChallenged			string							`json:"fb_id_challenged"`
	KeyChallenge				string							`json:"key_challenge"`
	ModeScene						int									`json:"mode_scene"`
	NumberItem					int									`json:"number_item"`
	Board								[]BoardModel				`json:"board"`
	ResultChallenger		ChallengerModel			`json:"result_challenger"`
	ResultChallenged		ChallengerModel			`json:"result_challenged"`
}