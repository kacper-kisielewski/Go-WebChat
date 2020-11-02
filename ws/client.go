package ws

import (
	"Website/settings"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/gorilla/websocket"
	"github.com/kyokomi/emoji"
)

//Client struct
type Client struct {
	Username      string
	Channel       string
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

func newClient(w http.ResponseWriter, req *http.Request, username, channel string) (Client, error) {
	for _, client := range getClientsInChannel(channel) {
		if client.Username == username {
			return Client{}, errors.New("User already in this channel")
		}
	}

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return Client{}, err
	}

	return Client{username, channel, conn, time.Unix(0, 0)}, nil
}

func addClient(client Client) {
	clients = append(clients, client)

	log.Printf(
		"New connection: %s [%s (#%s)] --> %s",
		client.Conn.RemoteAddr().String(),
		client.Username,
		client.Channel,
		client.Conn.LocalAddr().String(),
	)

	client.SendTo(fmt.Sprintf(
		"Welcome to #%s - %s here",
		client.Channel,
		pluralize.NewClient().Pluralize(
			"user", len(getClientsInChannel(client.Channel))-1, true),
	), settings.ChatSystemUsername)
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

func getClientsInChannel(channel string) []Client {
	var clientList []Client

	for _, client := range clients {
		if client.Channel == channel {
			clientList = append(clientList, client)
		}
	}
	return clientList
}
