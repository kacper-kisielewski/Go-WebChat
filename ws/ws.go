package ws

import (
	"Website/jwt"
	"Website/settings"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/gorilla/websocket"
	"github.com/kyokomi/emoji"
	"github.com/microcosm-cc/bluemonday"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  settings.ReadBufferSize,
		WriteBufferSize: settings.WriteBufferSize,
	}
	clients   []Client
	sanitizer = bluemonday.StrictPolicy()
)

//Client struct
type Client struct {
	Username      string
	Conn          *websocket.Conn
	LastMessageAt time.Time
}

//SendTo sends a message to a client
func (c *Client) SendTo(message, authorUsername string) error {
	return c.Conn.WriteJSON(Message{
		Message:        emoji.Sprint(sanitizer.Sanitize(message)),
		AuthorUsername: authorUsername,
	})
}

func newClient(w http.ResponseWriter, req *http.Request, username string) (Client, error) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return Client{}, err
	}

	return Client{username, conn, time.Unix(0, 0)}, nil
}

//Message struct
type Message struct {
	Message        string
	AuthorUsername string
}

//ChatHandler handles server chat
func ChatHandler(w http.ResponseWriter, req *http.Request) {
	username, _, err := jwt.GetUsernameAndEmailFromToken(req.Header.Get("Sec-Websocket-Protocol"))
	if err != nil {
		return
	}

	client, err := newClient(w, req, username)
	if err != nil {
		return
	}

	addClient(client)
	defer removeClient(client)

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			return
		}

		if time.Now().Sub(client.LastMessageAt) <= settings.ChatMessageCooldown {
			client.SendTo(
				fmt.Sprintf("Please wait %s before sending a message", settings.ChatMessageCooldown.String()),
				settings.ChatSystemUsername)
			continue
		}
		client.LastMessageAt = time.Now()
		if !checkMessage(string(message)) {
			continue
		}

		brodcast(string(message), client.Username)
	}
}

func checkMessage(message string) bool {
	if len(message) > settings.MaximumChatMessageLength {
		return false
	}

	if len(strings.TrimSpace(message)) == 0 {
		return false
	}

	return true
}

func brodcast(message, authorUsername string) {
	for _, client := range clients {
		log.Printf("[%s] %s", authorUsername, message)
		if client.SendTo(message, authorUsername) != nil {
			removeClient(client)
		}
	}
}

func addClient(client Client) {
	clients = append(clients, client)
	client.SendTo(fmt.Sprintf(
		"Welcome to %s - %s online",
		settings.SiteName,
		pluralize.NewClient().Pluralize("user", len(clients), true),
	), settings.ChatSystemUsername)
	log.Printf(
		"New connection: %s [%s] --> %s",
		client.Conn.RemoteAddr().String(),
		client.Username,
		client.Conn.LocalAddr().String(),
	)
}

func removeClient(client Client) {
	for i, _client := range clients {
		if _client == client {
			clients = append(clients[:i], clients[i+1:]...)
			log.Printf("Connection to %s [%s] closed",
				client.Conn.RemoteAddr().String(),
				client.Username)
			client.Conn.Close()
			return
		}
	}
}
