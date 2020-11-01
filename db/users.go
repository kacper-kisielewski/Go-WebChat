package db

import "strings"

//GetUserByEmail returns user model from email
func GetUserByEmail(email string) User {
	var user User

	DB.First(&user, `LOWER(email) = ?`, strings.ToLower(email)).Scan(&user)

	return user
}

//GetUserByUsername returns user model from username
func GetUserByUsername(username string) User {
	var user User

	DB.First(&user, `username ILIKE ?`, username).Scan(&user)

	return user
}

//UserExists checks whether user exists in the database
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
