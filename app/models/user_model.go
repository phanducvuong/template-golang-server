package models

type UserModel struct {
	FBId						string		`json:"fb_id"`
	Name						string		`json:"name"`
	Avatar					string		`json:"avatar"`
	Achievement			[]string	`'json:"achievement"`
	ScoreChallenge	int				`json:"score_challenge"`
	NotiChallenge		[]string	`json:"noti_challenge"`
}