package models

type BoardModel struct {
	Row			int			`json:"row"`
	Col			int			`json:"col"`
	Data		string	`json:"data"`
}