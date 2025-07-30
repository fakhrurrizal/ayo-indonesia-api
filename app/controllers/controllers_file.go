package controllers

import (
	"ayo-indonesia-api/app/models"
	repository "ayo-indonesia-api/app/repositories"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"ayo-indonesia-api/config"
	"fmt"
	"io"
	"mime/multipart"
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
func UploadFile(c *gin.Context) error {

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
		err = nil
	}
	storefolder := c.Query("folder")
	userIDData, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Wrong Authorization Token"))
	}
	userID := userIDData.(int)
	acceptedTypes := []string{
		"image/png",
		"image/jpeg",
		"image/gif",
		"image/webp",
		"video/quicktime",
		"video/mp4",
		"application/pdf",
		"text/csv",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-excel.sheet.macroenabled.12",
		"image/svg+xml",
		"image/svg",
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Get the MIME type of the file
	fileType := file.Header.Get("Content-Type")
	extension := ".jpg"
	if fileType == "image/png" {
		extension = ".png"
	}
	if fileType == "image/jpeg" {
		extension = ".jpg"
	}
	if fileType == "image/webp" {
		extension = ".webp"
	}
	if fileType == "image/gif" {
		extension = ".gif"
	}
	if fileType == "video/quicktime" {
		extension = ".mov"
	}
	if fileType == "video/mp4" {
		extension = ".mov"
	}
	if fileType == "application/pdf" {
		extension = ".pdf"
	}
	if fileType == "application/vnd.ms-excel" {
		extension = ".xls"
	}
	if fileType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		extension = ".xlsx"
	}
	if fileType == "application/vnd.ms-excel.sheet.macroenabled.12" {
		extension = ".et"
	}
	if fileType == "text/csv" {
		extension = ".csv"
	}
	if fileType == "image/svg+xml" || fileType == "image/svg" {
		extension = ".svg"
	}
	// Check if the MIME type is accepted
	var isAccepted bool
	for _, t := range acceptedTypes {
		if t == fileType {
			isAccepted = true
			break
		}
	}

	if !isAccepted {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"accepted_type": acceptedTypes,
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create a destination file
	t := time.Now().In(location)

	time := t.Format("2006-01")
	if storefolder != "" {
		storefolder = storefolder + "/"
	}
	folder := storefolder + strconv.Itoa(int(userID)) + "/" + time
	err = os.MkdirAll(config.RootPath()+"/assets/uploads/"+folder, os.ModePerm)
	if err != nil {
		return err
	}
	timestamp := strconv.Itoa(int(t.Unix()))

	filePath := filepath.Join(config.RootPath()+"/assets/uploads/", folder, timestamp+extension)

	fmt.Println(filePath)
	dst, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	data, err := SaveFileToDatabase(userID, folder+"/"+timestamp+extension, filePath)
	if err != nil {
		c.JSON(utils.ParseHttpError(err))
	}
	data.FullUrl = config.LoadConfig().BaseUrl + "/assets/uploads/" + folder + "/" + timestamp + extension

	if config.LoadConfig().EnableEncodeID {
		if data.EncodedID == "" {
			data.EncodedID = utils.EndcodeID(int(data.ID))
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    data,
		"message": "File uploaded successfully",
	})
	return nil
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
	return
}

func GetFileControl(id int, param reqres.ReqPaging) (data reqres.ResPaging, err error) {

	data, err = repository.GetFile(id, param)
	if err != nil {
		return
	}

	return
}


func getExtensionForFileType(fileType string) string {
	switch fileType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "video/quicktime":
		return ".mov"
	case "video/mp4":
		return ".mp4"
	case "application/pdf":
		return ".pdf"
	default:
		return ".unknown"
	}
}

func saveFile(src multipart.File, destPath string) error {
	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}
