package repository

import (
	"ayo-indonesia-api/app/middlewares"
	"ayo-indonesia-api/app/models"
	"ayo-indonesia-api/config"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func SignIn(email, password string) (user models.GlobalUser, token string, err error) {
	err = config.DB.
		Where("email = '" + strings.ToLower(email) + "'").First(&user).Error
	if err != nil {
		return
	}
	err = middlewares.VerifyPassword(password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		err = errors.New("incorrect password")
		return
	}

	token, err = middlewares.AuthMakeToken(user)
	if err != nil {
		return
	}
	return
}
