package main

import (
	"rank-server-pikachu/app"
)

func main() {

	// gin --appPort 5000 --port 3000   : cmd hotreload golang with gin

	a := app.App{}
	a.Initialize()
	a.Run(":5000")
}
