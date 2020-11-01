package db

import (
	"Website/settings"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//AuthenticateUser checks whether credientials are correct
func AuthenticateUser(hashedPassword, password []byte) bool {
	if bcrypt.CompareHashAndPassword(hashedPassword, password) != nil {
		return false
	}
	return true
}

//RegisterUser adds new user to the database
func RegisterUser(username, email string, password []byte) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, settings.Salt)
	if err != nil {
		return err
	}

	user := &User{
		Username:       username,
		Email:          strings.ToLower(email),
		HashedPassword: hashedPassword,
	}
	if UserExists(username, user.Email) {
		return errors.New("User already exists")
	}
	DB.Create(user)

	return nil
}
