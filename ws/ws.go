package ws

import (
	"Website/jwt"
	"Website/settings"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  settings.ReadBufferSize,
		WriteBufferSize: settings.WriteBufferSize,
	}
	clients []*websocket.Conn
)

//Message struct
type Message struct {
	Message        string
	AuthorUsername string
}

//ChatHandler handles server chat
func ChatHandler(w http.ResponseWriter, req *http.Request) {
	user, err := jwt.GetUserFromToken(req.Header.Get("Sec-Websocket-Protocol"))
	if err != nil {
		return
	}

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	addClient(conn)
	defer removeClient(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		brodcast(string(message), user.Username)
	}
}

func brodcast(message, authorUsername string) {
	for _, client := range clients {
		if client.WriteJSON(Message{
			Message:        message,
			AuthorUsername: authorUsername,
		}) != nil {
			removeClient(client)
		}
	}
}

func addClient(client *websocket.Conn) {
	clients = append(clients, client)
	log.Printf(
		"New connection: %s --> %s",
		client.RemoteAddr().String(),
		client.LocalAddr().String(),
	)
}

func removeClient(client *websocket.Conn) {
	for i, _client := range clients {
		if _client == client {
			clients = append(clients[:i], clients[i+1:]...)
			log.Printf("Connection to %s closed", client.RemoteAddr().String())
			client.Close()
			return
		}
	}
}
