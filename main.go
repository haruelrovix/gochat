package main

import (
	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/negroni"

	"./lib"
)

// Config
const database = "chat.db"
const port = ":9000"

var db *sql.DB

func main() {
	db, _ = sql.Open("sqlite3", database)
	lib.DB = db

	// Let's start registering a couple of URL paths and handlers
	router := newRouter()

	// Provide some default middlewares
	n := negroni.Classic()
	n.UseHandler(router)

	// Run the negroni stack as an HTTP server on port 9000
	n.Run(port)
}

func newRouter() *mux.Router {
	// If the route path is "/path/", accessing "/path" will perform a redirect
	// to the former and vice versa.
	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("GET").
		Path("/ws").
		Name("Communication Channel").
		HandlerFunc(lib.ServeWebSocket)

	router.
		Methods("POST").
		Path("/chat").
		Name("Send a message").
		HandlerFunc(lib.SendMessage)

	router.
		Methods("GET").
		Path("/chat").
		Name("Retrieve messages").
		HandlerFunc(lib.GetMessages)

	return router
}
