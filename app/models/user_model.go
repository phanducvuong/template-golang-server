package models

type UserModel struct {
	FBId						string		`json:"fb_id" bson:"fb_id"`
	Name						string		`json:"name" bson:"name"`
	Avatar					string		`json:"avatar" bson:"avatar"`
	Achievement			[]string	`'json:"achievement" bson:"achievement"`
	ScoreChallenge	int				`json:"score_challenge" bson:"score_challenge"`
	NotiChallenge		[]string	`json:"noti_challenge" bson:"noti_challenge"`
}