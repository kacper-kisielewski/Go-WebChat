package ws

import (
	"Website/jwt"
	"Website/settings"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/microcosm-cc/bluemonday"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  settings.ReadBufferSize,
		WriteBufferSize: settings.WriteBufferSize,
	}
	clients   []*Client
	sanitizer = bluemonday.StrictPolicy()
)

//Message struct
type Message struct {
	Message        string
	AuthorUsername string
}

//ChatHandler handles server chat
func ChatHandler(w http.ResponseWriter, req *http.Request, channel string) {
	username, _, err := jwt.GetUsernameAndEmailFromToken(req.Header.Get("Sec-Websocket-Protocol"))
	if err != nil {
		return
	}

	client, err := newClient(w, req, username, channel)
	if err != nil {
		return
	}

	addClient(&client)
	defer removeClient(&client)

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

		brodcast(string(message), client)
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

func brodcast(message string, author Client) {
	for _, client := range clients {
		if client.Channel != author.Channel {
			continue
		}
		log.Printf("[#%s] [%s] %s", author.Channel, author.Username, message)
		if client.SendTo(message, author.Username) != nil {
			removeClient(client)
		}
	}
}
