package controllers

import (
	repository "ayo-indonesia-api/app/repositories"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

// SignIn godoc
// @Summary SignIn
// @Description SignIn
// @Tags Auth
// @Accept json
// @Produce json
// @Param signin body reqres.SignInRequest true "SignIn user"
// @Success 200
// @Router /v1/auth/signin [post]
// @Security ApiKeyAuth
func SignIn(c *gin.Context) {
	var req reqres.SignInRequest

	// Bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
		return
	}

	user, accessToken, err := repository.SignIn(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"error":  err.Error(),
		})
		return
	}

	userData, _ := repository.GetUserByIDPlain(int(user.ID))

	userResponse := reqres.GlobalUserAuthResponse{
		ID:       int(userData.ID),
		Fullname: userData.Fullname,
		Email:    userData.Email,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":         userResponse,
			"access_token": accessToken,
			"expiration":   time.Now().Add(72 * time.Hour).Format("2006-01-02 15:04:05"),
		},
	})
}
