package repository

import (
	"ayo-indonesia-api/app/models"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"ayo-indonesia-api/config"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
)

func GetTeams(param reqres.ReqPaging) (data reqres.ResPaging) {
	var responses []models.Team
	where := "deleted_at IS NULL"

	var modelTotal []models.Team

	type TotalResult struct {
		Total       int64
		LastUpdated time.Time
	}
	var totalResult TotalResult
	config.DB.Model(&modelTotal).Select("COUNT(*) AS total, MAX(updated_at) AS last_updated").Scan(&totalResult)

	if param.Custom != "" {
		where += " AND status = " + param.Custom.(string)
	}
	if param.Search != "" {
		where += " AND (name ILIKE '%" + param.Search + "%' OR city ILIKE '%" + param.Search + "%')"
	}

	var totalFiltered int64
	config.DB.Model(&modelTotal).Where(where).Count(&totalFiltered)

	config.DB.Limit(param.Limit).Offset(param.Offset).Order(param.Sort + " " + param.Order).Where(where).Find(&responses)

	var responsesRefined []reqres.TeamResponse
	for _, item := range responses {
		responseRefined := BuildTeamResponse(item)

		responsesRefined = append(responsesRefined, responseRefined)
	}

	data = utils.PopulateResPaging(&param, responsesRefined, totalResult.Total, totalFiltered, null.TimeFrom(totalResult.LastUpdated))

	return
}

func GetTeamByID(id uint) (responseRefined reqres.TeamResponse, err error) {
	var response models.Team
	err = config.DB.First(&response, id).Error

	responseRefined = BuildTeamResponse(response)

	return
}

func GetTeamByIDPlain(id uint) (response models.Team, err error) {
	err = config.DB.First(&response, id).Error

	return
}

func CreateTeam(data *reqres.TeamRequest) (response models.Team, err error) {

	response = models.Team{
		Name:        data.Name,
		Logo:        data.Logo,
		FoundedYear: data.FoundedYear,
		Address:     data.Address,
		City:        data.City,
	}

	var created bool
	for !created {
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

func UpdateTeam(request models.Team) (response models.Team, err error) {
	err = config.DB.Save(&request).Scan(&response).Error

	return
}

func DeleteTeam(request models.Team) (models.Team, error) {
	err := config.DB.Delete(&request).Error
	return request, err
}

func BuildTeamResponse(data models.Team) (response reqres.TeamResponse) {

	response.CustomGormModel = data.CustomGormModel
	response.Name = data.Name
	response.Logo = data.Logo
	response.FoundedYear = data.FoundedYear
	response.Address = data.Address
	response.City = data.City

	totalPlayer := CountPlayer(data.ID)
	response.PlayersCount = int(totalPlayer)

	return response
}
