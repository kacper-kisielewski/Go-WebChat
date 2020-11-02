package views

import (
	"Website/settings"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("username", UsernameValidator)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("email", EmailValidator)
	}
}

//UsernameValidator validator
var UsernameValidator validator.Func = func(fl validator.FieldLevel) bool {
	username, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if len(username) > settings.MaximumUsernameLength {
		return false
	}

	if len(username) < settings.MinimumUsernameLength {
		return false
	}

	for _, char := range username {
		if !strings.Contains(settings.UsernameWhitelistedCharacters, string(char)) {
			return false
		}
	}

	return true
}

//EmailValidator validator
var EmailValidator validator.Func = func(fl validator.FieldLevel) bool {
	email, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if checkmail.ValidateFormat(email) != nil {
		return false
	}

	if len(email) < settings.MinimumEmailLength {
		return false
	}

	return true
}

//IsValidChannelName checks whether channel name is valid
func IsValidChannelName(name string) bool {
	if len(name) > settings.MaximumChannelNameLength {
		return false
	}

	if len(name) < settings.MinimumChannelNameLength {
		return false
	}

	if settings.ChannelNameRegex.FindString(name) != "" {
		return false
	}

	return true
}
