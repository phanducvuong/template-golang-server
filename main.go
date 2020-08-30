package main

import (
	"rank-server-pikachu/app"
)

func main() {
	a := app.App{}
	a.Initialize()
	a.Run(":3000")
}
