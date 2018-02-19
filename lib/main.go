package lib

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type Chat struct {
	Message string    `json:"message"`
	Sent    time.Time `json:"sent"`
}

// Configure the upgrader and database
var upgrader = websocket.Upgrader{}
var DB *sql.DB

// API to retrieve message at realtime
func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
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
	tx, err := DB.Begin()
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
func SendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		return
	}

	// Dump message to database
	insertMessage(body)

	if err := json.NewEncoder(w).Encode(http.StatusOK); err != nil {
		panic(err)
	}
}

// API to retrieve messages
func GetMessages(w http.ResponseWriter, r *http.Request) {
	chat := []Chat{}
	rows, _ := DB.Query("select message, timestamp from chat")
	for rows.Next() {
		row := Chat{}
		rows.Scan(&row.Message, &row.Sent)
		chat = append(chat, row)
	}

	if err := json.NewEncoder(w).Encode(chat); err != nil {
		panic(err)
	}
}
