package db

import (
	"Website/settings"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//GetUserByEmail returns user model from email
func GetUserByEmail(email string) User {
	var user User

	DB.First(&user, `LOWER(email) = ?`, email).Scan(&user)

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
func RegisterUser(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), settings.Salt)
	if err != nil {
		return err
	}

	user := &User{
		Username:       username,
		Email:          strings.ToLower(email),
		HashedPassword: hashedPassword,
	}
	if err := ValidateUser(user); err != nil {
		return err
	}

	DB.Create(user)

	return nil
}

//ValidateUser validates user model
func ValidateUser(user *User) error {
	if UserExists(user.Username, user.Email) {
		return errors.New("User already exists")
	}

	if len(user.Username) > settings.MaximumPasswordLength {
		return fmt.Errorf("Username length cannot exceed %d characters", settings.MaximumPasswordLength)
	}
	if len(user.Username) <= settings.MinimumPasswordLength {
		return fmt.Errorf("Username has to be at least %d characters", settings.MinimumPasswordLength)
	}

	if strings.Count(user.Email, "@") != 1 {
		return errors.New("Invalid email address")
	}

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
