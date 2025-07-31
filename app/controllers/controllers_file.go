package controllers

import (
	"ayo-indonesia-api/app/models"
	repository "ayo-indonesia-api/app/repositories"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"ayo-indonesia-api/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadFile godoc
// @Summary File Uploader
// @Description File Uploader
// @Tags File
// @Accept mpfd
// @Param file formData file true "File to upload"
// @Param folder query string false "Folder to store"
// @Produce json
// @Success 200
// @Router /v1/file [post]
// @Security ApiKeyAuth
// @Security JwtToken
func UploadFile(c *gin.Context) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
	}

	storefolder := c.Query("folder")

	userIDData, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Wrong Authorization Token"))
		return
	}
	userID := userIDData.(int)

	// File validation
	acceptedTypes := []string{
		"image/png", "image/jpeg", "image/gif", "image/webp", "video/quicktime", "video/mp4",
		"application/pdf", "text/csv", "application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-excel.sheet.macroenabled.12", "image/svg+xml", "image/svg",
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	fileType := file.Header.Get("Content-Type")
	extension := map[string]string{
		"image/png": ".png", "image/jpeg": ".jpg", "image/gif": ".gif", "image/webp": ".webp",
		"video/quicktime": ".mov", "video/mp4": ".mp4", "application/pdf": ".pdf", "text/csv": ".csv",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
		"application/vnd.ms-excel.sheet.macroenabled.12":                    ".et",
		"image/svg+xml": ".svg", "image/svg": ".svg",
	}[fileType]

	if extension == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed", "accepted_types": acceptedTypes})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	t := time.Now().In(location)
	timeStr := t.Format("2006-01")
	if storefolder != "" {
		storefolder += "/"
	}
	folder := fmt.Sprintf("%s%d/%s", storefolder, userID, timeStr)

	err = os.MkdirAll(filepath.Join(config.RootPath(), "assets/uploads/", folder), os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload folder"})
		return
	}

	timestamp := strconv.FormatInt(t.Unix(), 10)
	fileName := timestamp + extension
	filePath := filepath.Join(config.RootPath(), "assets/uploads", folder, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	data, err := SaveFileToDatabase(userID, folder+"/"+fileName, filePath)
	if err != nil {
		c.JSON(utils.ParseHttpError(err))
		return
	}

	data.FullUrl = config.LoadConfig().BaseUrl + "/assets/uploads/" + folder + "/" + fileName

	if config.LoadConfig().EnableEncodeID && data.EncodedID == "" {
		data.EncodedID = utils.EndcodeID(int(data.ID))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    data,
		"message": "File uploaded successfully",
	})
}

func SaveFileToDatabase(id int, filename, path string) (data models.GlobalFile, err error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
		err = nil
	}
	t := time.Now().In(location).Unix()
	data = models.GlobalFile{
		Token:    strconv.Itoa(int(t)) + utils.GenerateRandomString(5),
		UserID:   id,
		Filename: filename,
		Path:     path,
	}

	err = repository.SaveFile(&data)
	return
}

// GetFile godoc
// @Summary Mendapatkan List Files
// @Description Mendapatkan List Files
// @Tags File
// @Accept json
// @Param search query string false "search (string)"
// @Param page query integer false "page (int)"
// @Param limit query integer false "limit (int)"
// @Param token query string false "token (string)"
// @Produce json
// @Success 200
// @Router /v1/file [get]
// @Security ApiKeyAuth
// @Security JwtToken
func GetFile(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	userID, ok := userIDAny.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user_id"})
		return
	}

	param := utils.PopulatePaging(c, "token")

	data, err := GetFileControl(userID, param)
	if err != nil {
		code, res := utils.ParseHttpError(err)
		c.JSON(code, res)
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetFileControl(id int, param reqres.ReqPaging) (data reqres.ResPaging, err error) {

	data, err = repository.GetFile(id, param)
	if err != nil {
		return
	}

	return
}

// func getExtensionForFileType(fileType string) string {
// 	switch fileType {
// 	case "image/png":
// 		return ".png"
// 	case "image/jpeg":
// 		return ".jpg"
// 	case "image/gif":
// 		return ".gif"
// 	case "video/quicktime":
// 		return ".mov"
// 	case "video/mp4":
// 		return ".mp4"
// 	case "application/pdf":
// 		return ".pdf"
// 	default:
// 		return ".unknown"
// 	}
// }

// func saveFile(src multipart.File, destPath string) error {
// 	dst, err := os.Create(destPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer dst.Close()

// 	_, err = io.Copy(dst, src)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
