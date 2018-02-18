package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/negroni"
)

type Chat struct {
	Message string    `json:"message"`
	Sent    time.Time `json:"sent"`
}

// Configure the upgrader and database
var upgrader = websocket.Upgrader{}
var db *sql.DB

func main() {
	db, _ = sql.Open("sqlite3", "chat.db")

	// Let's start registering a couple of URL paths and handlers
	router := newRouter()

	// Provide some default middlewares
	n := negroni.Classic()
	n.UseHandler(router)

	// Run the negroni stack as an HTTP server on port 9000
	n.Run(":9000")
}

// API to retrieve message at realtime
func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	// Handle request origin not allowed
	allowAllOrigin := func(r *http.Request) bool { return true }
	upgrader.CheckOrigin = allowAllOrigin

	// Upgrade the HTTP server connection to the WebSocket protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a long lived connection that will retrieve message at realtime
	for {
		// Call the connection's WriteMessage and ReadMessage methods to send and
		// receive messages as a slice of bytes
		// https://godoc.org/github.com/gorilla/websocket
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

		// Dump message to database
		insertMessage(p)
	}
}

func insertMessage(message []byte) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into chat(message) values(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(message)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

// API to send message
func sendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	insertMessage(body)
}

// API to retrieve messages
func getMessages(w http.ResponseWriter, r *http.Request) {
	chat := []Chat{}
	rows, _ := db.Query("select message, timestamp from chat")
	for rows.Next() {
		row := Chat{}
		rows.Scan(&row.Message, &row.Sent)
		chat = append(chat, row)
	}

	if err := json.NewEncoder(w).Encode(chat); err != nil {
		panic(err)
	}
}

func newRouter() *mux.Router {
	// If the route path is "/path/", accessing "/path" will perform a redirect
	// to the former and vice versa.
	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("GET").
		Path("/ws").
		Name("Communication Channel").
		HandlerFunc(serveWebSocket)

	router.
		Methods("POST").
		Path("/chat").
		Name("Send a message").
		HandlerFunc(sendMessage)

	router.
		Methods("GET").
		Path("/chat").
		Name("Retrieve messages").
		HandlerFunc(getMessages)

	return router
}
