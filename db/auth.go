package db

import (
	"Website/settings"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//GetUserByEmail returns user model from email
func GetUserByEmail(email string) User {
	var user User

	DB.First(&user, `LOWER(email) = ?`, strings.ToLower(email)).Scan(&user)

	return user
}

//AuthenticateUser checks whether credientials are correct
func AuthenticateUser(email string, password []byte) bool {
	user := GetUserByEmail(email)

	if bcrypt.CompareHashAndPassword(user.HashedPassword, password) != nil {
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

//UserExists checks whether user already exists in the database
func UserExists(username, email string) bool {
	var (
		user      User
		userCount int64
	)

	DB.Model(&user).Where(
		`username ILIKE ? OR email ILIKE ?`, username, email,
	).Count(&userCount)
	return int(userCount) >= 1
}
