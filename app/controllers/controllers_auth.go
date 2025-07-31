package controllers

import (
	repository "ayo-indonesia-api/app/repositories"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

// SignUp godoc
// @Summary SignUp
// @Description SignUp
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body reqres.SignUpRequest true "SignUp user"
// @Success 200 {object} map[string]interface{}
// @Router /v1/auth/signup [post]
// @Security ApiKeyAuth
func SignUp(c *gin.Context) {
	var request reqres.SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
		return
	}

	utils.StripTagsFromStruct(&request)

	if err := request.Validate(); err != nil {
		errVal := err.(validation.Errors)
		c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
		return
	}

	_, err := repository.GetUserByEmail(strings.ToLower(request.Email))
	if err == nil {
		c.JSON(http.StatusBadRequest, utils.Respond(http.StatusBadRequest, "bad request", "email sudah terdaftar"))
		return
	}

	inputUser := reqres.GlobalUserRequest{
		Fullname: request.Fullname,
		Email:    request.Email,
		Password: request.Password,
	}

	_, err = repository.CreateUser(&inputUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewBadRequestError([]map[string]interface{}{
			{
				"field": "Email",
				"error": err.Error(),
			},
		}))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "Registration Successful",
	})
}

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

// GetSignInUser godoc
// @Summary Get Sign In User
// @Description Get Sign In User
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/auth/user [get]
// @Security ApiKeyAuth
// @Security JwtToken
func GetSignInUser(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized access",
		})
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Invalid user ID type",
		})
		return
	}

	user, err := repository.GetUserByIDPlain(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to get user data",
			"error":   err.Error(),
		})
		return
	}

	data := reqres.GlobalUserAuthResponse{
		ID:       int(user.ID),
		Fullname: user.Fullname,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    data,
		"message": "Success to get user data",
	})
}
