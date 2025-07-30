package middlewares

import (
	"ayo-indonesia-api/app/models"
	"ayo-indonesia-api/app/utils"
	"ayo-indonesia-api/config"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		bearerToken := strings.Split(authorizationHeader, " ")

		if len(bearerToken) != 2 {
			c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Incorrect Authorization Token"))
			c.Abort()
			return
		}

		tokenStr := bearerToken[1]

		UserID, err := ValidateToken(tokenStr)
		if err != nil {
			fmt.Println("Token Validation Error:", err)
			c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError(err.Error()))
			c.Abort()
			return
		}

		c.Set("user_id", UserID)
		c.Next()
	}
}

func ValidateToken(tokenString string) (userID int, err error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
		err = nil
	}
	tokenStringbyt, err := hex.DecodeString(tokenString)
	if err != nil {
		err = errors.New("incorrect token format")
		return
	}
	str := string(tokenStringbyt)
	newtStr := strings.Replace(string(str), config.LoadConfig().AppKey, "", -1)
	decoded, err := base64.StdEncoding.DecodeString(newtStr)
	if err != nil {
		err = errors.New("incorrect token format")
		return
	}
	newStr := strings.Replace(string(decoded), config.LoadConfig().AppKey, "", -1)
	newdecoded, err := base64.StdEncoding.DecodeString(newStr)
	if err != nil {
		err = errors.New("incorrect token format")
		return
	}
	parts := strings.Split(string(newdecoded), "&")
	expiredAt, _ := strconv.Atoi(parts[1])
	if expiredAt < int(time.Now().In(location).Unix()) {
		err = errors.New("incorrect token format")
		return
	}
	userID, _ = strconv.Atoi(parts[0])

	return
}

func AuthMakeToken(user models.GlobalUser) (string, error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
		err = nil
	}

	ExpiresAt := time.Now().In(location).Add(7 * 24 * time.Hour).Unix()
	str := fmt.Sprintf("%v&%v", user.ID, ExpiresAt)
	encoded := base64.StdEncoding.EncodeToString([]byte(str)) + config.LoadConfig().AppKey
	token := base64.StdEncoding.EncodeToString([]byte(encoded)) + config.LoadConfig().AppKey
	token = hex.EncodeToString([]byte(token))
	return token, nil
}
