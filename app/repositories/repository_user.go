package repository

import (
	"ayo-indonesia-api/app/middlewares"
	"ayo-indonesia-api/app/models"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"ayo-indonesia-api/config"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

func CreateUser(data *reqres.GlobalUserRequest) (response models.GlobalUser, err error) {
		
		_, errUser := GetUserByEmail(strings.ToLower(data.Email))
		if errUser == nil {
			err = errors.New("email has been registered")
			return
		}

		response = models.GlobalUser{
			Fullname: data.Fullname,
			Email:    strings.ToLower(data.Email),
			Password: middlewares.BcryptPassword(data.Password),
		}

		var created bool
		for !created {
			respectiveID, _ := config.GetRespectiveID(config.DB, response.TableName(), true)
			response.ID = respectiveID
			err = config.DB.Create(&response).Error

			if err != nil {
				if !config.LoadConfig().EnableIDDuplicationHandling {
					return
				}
				if pqErr, ok := err.(*pq.Error); ok {
					if pqErr.Code != "23505" {
						return
					}
				}
			} else {
				created = true
			}
		}

	return
}

func BuildUserResponse(data models.GlobalUser) (response reqres.GlobalUserResponse) {

	response.CustomGormModel = data.CustomGormModel
	response.Fullname = data.Fullname
	response.Email = strings.ToLower(data.Email)

	return response
}

func GetUsers(param reqres.ReqPaging, withTotalSale, withTotalPurchase bool) (data reqres.ResPaging) {
	var responses []models.GlobalUser
	where := "deleted_at IS NULL"

	var modelTotal []models.GlobalUser

	type TotalResult struct {
		Total       int64
		LastUpdated time.Time
	}
	var totalResult TotalResult
	config.DB.Model(&modelTotal).Select("COUNT(*) AS total, MAX(updated_at) AS last_updated").Scan(&totalResult)

	if param.Search != "" {
		where += " AND (fullname ILIKE '%" + param.Search + "%' OR email ILIKE '%" + param.Search + "%')"
	}

	var totalFiltered int64
	config.DB.Model(&modelTotal).Where(where).Count(&totalFiltered)

	config.DB.Limit(param.Limit).Offset(param.Offset).Order(param.Sort + " " + param.Order).Where(where).Find(&responses)

	var responsesRefined []reqres.GlobalUserResponse
	for _, item := range responses {
		responseRefined := BuildUserResponse(item)

		responsesRefined = append(responsesRefined, responseRefined)
	}

	data = utils.PopulateResPaging(&param, responsesRefined, totalResult.Total, totalFiltered)

	return
}

func GetAllUsers(companyID, status int) (users []models.GlobalUser, err error) {
	where := "deleted_at IS NULL"

	err = config.DB.Where(where).Find(&users).Error

	return
}

func GetUserByID(id, companyID int) (response reqres.GlobalUserResponse, err error) {
	fmt.Println("companyID:", companyID)
	var data models.GlobalUser
	err = config.DB.First(&data, id).Error

	response = BuildUserResponse(data)

	return
}

func GetUserByIDPlain(id int) (response models.GlobalUser, err error) {
	err = config.DB.First(&response, id).Error

	return
}

func GetUserByEmail(email string) (response models.GlobalUser, err error) {
	err = config.DB.Where("email = ?", strings.ToLower(email)).First(&response).Error

	return
}

func DeleteUser(request models.GlobalUser) (models.GlobalUser, error) {
	var err error
	err = config.DB.Delete(&request).Error

	return request, err
}

func UpdateUser(request models.GlobalUser) (response models.GlobalUser, err error) {
	err = config.DB.Save(&request).Scan(&response).Error

	return
}
