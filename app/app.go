package app

import (
	"rank-server-pikachu/app/controllers"
	"github.com/gorilla/mux"
	"rank-server-pikachu/app/hello"
	"context"
	"fmt"
	"log"
	"net/http"
	"rank-server-pikachu/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Router 	*mux.Router
	DB 			*mongo.Database
}

func (app *App) Initialize() {
	clientOption 	:= options.Client().ApplyURI(config.URLMongodb)
	client, err 	:= mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	app.DB			= client.Database("pikachu-db")
	app.Router 	= mux.NewRouter()

	app.setRouter()
	fmt.Println("Connect mongodb success")
}

func (app *App) setRouter() {
	app.Router.HandleFunc("/", hello.GetProfile).Methods("GET")
	app.Router.HandleFunc("/get", app.handleRequest(hello.GetDB)).Methods("GET")
	app.Router.HandleFunc("/get-all", app.handleRequest(hello.GetAllData)).Methods("GET")
	app.Router.HandleFunc("/libgdx", app.handleRequest(hello.TestPostData)).Methods("POST")

	//route user
	app.Router.HandleFunc(config.PathAPI + "/user/init", app.handleRequest(controllers.InitUser)).Methods("POST")
	app.Router.HandleFunc(config.PathAPI + "/user/update-score", app.handleRequest(controllers.UpdateScoreUser)).Methods("POST")
	app.Router.HandleFunc(config.PathAPI + "/user/init-challenge", app.handleRequest(controllers.InitChallenge)).Methods("POST")
	app.Router.HandleFunc(config.PathAPI + "/get-leaderboard-by-level", app.handleRequest(controllers.GetLeaderboardByLevel)).Methods("GET")
}

type RequestHandlerFunc func(db *mongo.Database, w http.ResponseWriter, r *http.Request)
func (app *App) handleRequest(handler RequestHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(app.DB, w, r)
	}
}

func (app *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, app.Router))
}
