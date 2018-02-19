# gochat
Simple Chat API using Go and WebSocket

### Prerequisite
1. [Go](https://golang.org/)
```sh
$ go version
go version go1.9.4 darwin/amd64
```
2. [SQLite](https://www.sqlite.org/index.html)
```sh
$ sqlite3 version
SQLite version 3.16.0 2016-11-04 19:09:39
```

### How to Run
1. Clone this repository
```sh
$ git clone https://github.com/haruelrovix/gochat.git && cd gochat
```
2. Start SQLite using `chat.db`
```sh
$ sqlite3 chat.db
```
3. Create below schema
```sql
CREATE TABLE chat(
  pk integer primary key autoincrement,
  message text,
  timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
```
4. Check it again, just to make sure
```sh
sqlite> .schema
```
5. Open another terminal, run
```sh
$ go run main.go
```
6. If it asks to accept incoming network connections, allow it.
<img src="https://i.imgur.com/FqfijBf.png" alt="Accpet incoming network connections" width="30%" />

7. `gochat` listening on port 9000
```sh
[negroni] listening on :9000
```

### Test the API
To test `Post` and `Get`, I used [Insomnia REST Client](https://github.com/getinsomnia).
1. `Post` to `http://localhost:9000/chat`. It accepts one parameter: the string message to sent.
<img src="https://i.imgur.com/KJOAwln.png" alt="Post" width="50%" />

2. `Get` to `http://localhost:9000/chat`. It retrieves all previously-sent messages.
<img src="https://i.imgur.com/dDaS5Rh.png" alt="Get" width="50%" />

WebSocket can be tested using [Smart Websocket Client](https://chrome.google.com/webstore/detail/smart-websocket-client/omalebghpgejjiaoknljcfmglgbpocdp).

3. Connect to `ws://localhost:9000/ws`, send a message.
<img src="https://i.imgur.com/85ef1Xo.png" alt="WebSocket" width="50%" />

4. `Get` again, the message sent through WebSocket is also recorded.
```json
[
 {
  "message": "Hi!",
	"sent": "2018-02-18T18:38:12Z"
 },
 {
  "message": "asl pls",
  "sent": "2018-02-18T18:42:40Z"
 },
 {
  "message": "hobby?",
	"sent": "2018-02-18T18:42:53Z"
 },
 {
  "message": "a lonely programmer",
	"sent": "2018-02-18T18:48:33Z"
 }
]
```

### Logging
[Negroni](https://github.com/urfave/negroni), Idiomatic HTTP Middleware for Golang.
```java
[negroni] listening on :9000
[negroni] 2018-02-19T01:38:12+07:00 | 0 |        10.765667ms | localhost:9000 | POST /chat
[negroni] 2018-02-19T01:42:40+07:00 | 0 |        2.117104ms | localhost:9000 | POST /chat
[negroni] 2018-02-19T01:42:53+07:00 | 0 |        5.173429ms | localhost:9000 | POST /chat
[negroni] 2018-02-19T01:42:57+07:00 | 200 |      5.160381ms | localhost:9000 | GET /chat
[negroni] 2018-02-19T01:43:04+07:00 | 200 |      177.932µs | localhost:9000 | GET /chat
[negroni] 2018-02-19T01:50:52+07:00 | 200 |      1.622574ms | localhost:9000 | GET /chat
```

### Debugging
VS Code and [Delve](https://github.com/derekparker/delve), a debugger for the Go programming language.

<img src="https://i.imgur.com/kldvHPj.png" alt="Debugging" width="50%" />
