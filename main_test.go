package main

import (
	"Website/db"
	"Website/jwt"
	"Website/settings"
	"Website/ws"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var (
	testUsername = faker.Username()
	testEmail    = faker.Email()
	testPassword = faker.Password()

	testMessage = []byte("This is a test messsage")

	accessToken string
)

func TestIndex(t *testing.T) {
	assert.Equal(t, http.StatusOK, performRequest(
		setupRouter(), "GET", "/", nil,
	).Code)
}

func TestRegister(t *testing.T) {
	router := setupRouter()

	for i := 0; i < 2; i++ {
		sendRegisterRequest(router, testUsername, testEmail, testPassword)
	}

	sendRegisterRequest(router, strings.ToUpper(testUsername), testEmail, testPassword)
	sendRegisterRequest(router, testUsername, strings.ToUpper(testEmail), testPassword)
	sendRegisterRequest(router, strings.ToUpper(testUsername), strings.ToUpper(testEmail), testPassword)

	sendRegisterRequest(router, testUsername+"-", faker.Email(), faker.Password())

	for i := 1; i < 4; i++ {
		assert.NotEqual(t, http.StatusOK, sendRegisterRequest(
			router, strings.Repeat("a", settings.MaximumUsernameLength+i), faker.Email(), faker.Password(),
		))
	}

	for i := 0; i < settings.MinimumUsernameLength; i++ {
		assert.NotEqual(t, http.StatusOK, sendRegisterRequest(
			router, strings.Repeat("a", i), faker.Email(), faker.Password(),
		))
	}

	var (
		user      db.User
		userCount int64
	)

	db.DB.Model(&user).Where("email ILIKE ? OR username ILIKE ?", testEmail, testUsername).Count(&userCount)
	assert.Equal(t, 1, int(userCount))

	db.DB.Model(&user).Where("LENGTH(username) > ?", settings.MaximumUsernameLength).Count(&userCount)
	assert.Equal(t, 0, int(userCount))

	db.DB.Model(&user).Where("LENGTH(username < ?", settings.MinimumUsernameLength).Count(&userCount)
	assert.Equal(t, 0, int(userCount))
}

func TestLogin(t *testing.T) {
	var (
		router = setupRouter()

		resp     *httptest.ResponseRecorder
		respBody []byte
	)

	resp = sendLoginRequest(router, testEmail, testPassword)

	accessToken = strings.Split(resp.HeaderMap.Get("Set-Cookie"), "=")[1]
	assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)

	resp = sendLoginRequest(router, testEmail, testPassword+"a")
	respBody, _ = ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Equal(t, settings.LoginInvalidCredientialsMessage, string(respBody))
}

func TestAccessToken(t *testing.T) {
	user, err := jwt.GetUserFromToken(accessToken)
	assert.Nil(t, err)
	assert.Equal(t, strings.ToLower(testEmail), user.Email)
	assert.Equal(t, testUsername, user.Username)
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

func sendRegisterRequest(router *gin.Engine, username, email, password string) *httptest.ResponseRecorder {
	return performRequest(
		router, "POST", "/auth/register", url.Values{
			"username": {username},
			"email":    {email},
			"password": {password},
		},
	)
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
