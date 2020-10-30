package main

import (
	"Website/db"
	"Website/jwt"
	"Website/settings"
	"Website/ws"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

//Declare constants
const (
	testUsername = "User"
	testEmail    = "user@example.com"
	testPassword = "Password1234"
)

var (
	testMessage = []byte("This is a test messsage")

	accessToken string
)

//LoginBody struct
type LoginBody struct {
	email    string
	password string
}

func TestIndex(t *testing.T) {
	assert.Equal(t, http.StatusOK, performRequest(
		setupRouter(), "GET", "/", nil,
	).Code)
}

func TestAuthentication(t *testing.T) {
	for i := 0; i < 2; i++ {
		db.RegisterUser(testUsername, testEmail, testPassword)
	}

	db.RegisterUser(strings.ToUpper(testUsername), testEmail, testPassword)
	db.RegisterUser(testUsername, strings.ToUpper(testEmail), testPassword)
	db.RegisterUser(strings.ToUpper(testUsername), strings.ToUpper(testEmail), testPassword)

	var (
		user      db.User
		userCount int64
	)

	db.DB.Model(&user).Where("email ILIKE ? OR username ILIKE ?", testEmail, testUsername).Count(&userCount)
	assert.Equal(t, 1, int(userCount))

	router := setupRouter()

	var (
		resp     *httptest.ResponseRecorder
		respJSON map[string]string
		respBody []byte
		ok       bool
	)

	resp = sendLoginRequest(router, testEmail, testPassword)
	respBody, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.Code)

	json.Unmarshal(respBody, &respJSON)
	accessToken, ok = respJSON["access_token"]
	assert.True(t, ok)

	user, err := jwt.GetUserFromToken(accessToken)
	assert.Nil(t, err)
	assert.Equal(t, testEmail, user.Email)
	assert.Equal(t, testUsername, user.Username)

	resp = sendLoginRequest(router, testEmail, testPassword+"a")
	respBody, _ = ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Equal(t, settings.LoginInvalidCredientialsMessage, string(respBody))
}

func TestWebsocketChat(t *testing.T) {
	url, _ := url.Parse("ws://127.0.0.1:8000/chat")

	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	assert.Nil(t, err)

	wsConn, _, err := websocket.NewClient(conn, url, http.Header{
		"Sec-Websocket-Protocol": {accessToken},
	}, settings.ReadBufferSize, settings.WriteBufferSize)
	assert.Nil(t, err)

	var message ws.Message
	wsConn.WriteMessage(1, testMessage)

	assert.Nil(t, wsConn.ReadJSON(&message))
	assert.Equal(t, string(testMessage), message.Message)
	assert.Equal(t, testUsername, message.AuthorUsername)
}

func sendLoginRequest(router *gin.Engine, email, password string) *httptest.ResponseRecorder {
	return performRequest(
		router, "POST", "/auth/login", url.Values{
			"email":    {email},
			"password": {password},
		},
	)
}

func performRequest(r http.Handler, method, path string, form url.Values) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
