package repository

import (
	"ayo-indonesia-api/app/models"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"ayo-indonesia-api/config"
	"strconv"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
)

func GetPlayers(teamID int, param reqres.ReqPaging) (data reqres.ResPaging) {
	var responses []models.Player
	where := "deleted_at IS NULL"

	var modelTotal []models.Player

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
		where += " AND (name ILIKE '%" + param.Search + "%' OR position ILIKE '%" + param.Search + "%')"
	}

	if teamID > 0 {
		where += " AND team_id = " + strconv.Itoa(teamID)
	}

	var totalFiltered int64
	config.DB.Model(&modelTotal).Where(where).Count(&totalFiltered)

	config.DB.Limit(param.Limit).Offset(param.Offset).Order(param.Sort + " " + param.Order).Where(where).Find(&responses)

	var responsesRefined []reqres.PlayerResponse
	for _, item := range responses {
		responseRefined := BuildPlayerResponse(item)

		responsesRefined = append(responsesRefined, responseRefined)
	}

	data = utils.PopulateResPaging(&param, responsesRefined, totalResult.Total, totalFiltered, null.TimeFrom(totalResult.LastUpdated))

	return
}

func GetPlayerByID(id uint) (responseRefined reqres.PlayerResponse, err error) {
	var response models.Player
	err = config.DB.First(&response, id).Error

	responseRefined = BuildPlayerResponse(response)

	return
}

func CreatePlayer(data *reqres.PlayerRequest) (response models.Player, err error) {

	response = models.Player{
		Name:         data.Name,
		Height:       data.Height,
		Position:     data.Position,
		Weight:       data.Weight,
		JerseyNumber: data.JerseyNumber,
		TeamID:       data.TeamID,
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

func UpdatePlayer(request models.Player) (response models.Player, err error) {
	err = config.DB.Save(&request).Scan(&response).Error

	return
}

func DeletePlayer(request models.Player) (models.Player, error) {
	err := config.DB.Delete(&request).Error
	return request, err
}

func CountPlayer(teamID uint) (total int64) {
	where := "deleted_at IS NULL"

	if teamID > 0 {
		where += " AND team_id = " + strconv.Itoa(int(teamID))
	}

	var modelTotal []models.GlobalUser

	config.DB.Model(&modelTotal).Where(where).Count(&total)

	return total
}

func BuildPlayerResponse(data models.Player) (response reqres.PlayerResponse) {

	var team models.Team

	response.CustomGormModel = data.CustomGormModel
	response.Name = data.Name
	response.Height = data.Height
	response.Position = data.Position
	response.Weight = data.Weight
	response.JerseyNumber = data.JerseyNumber

	if data.TeamID > 0 {
		team, _ = GetTeamByIDPlain(data.TeamID)
		response.Team = reqres.GlobalIDNameResponse{
			ID:   int(team.ID),
			Name: team.Name,
		}
	}

	return response
}
