package util

import (
	"encoding/json"
	"rank-server-pikachu/app/models"
)

func ResponseUtil(statusCode int, msg string) []byte {
	res := models.ResModel {
		StatusCode	: statusCode,
		Msg					: msg,
	}
	jsRes,_		:= json.Marshal(res)
	return jsRes
}