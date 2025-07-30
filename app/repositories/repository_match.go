package repository

import (
	"ayo-indonesia-api/app/models"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"ayo-indonesia-api/config"
	"errors"
	"strconv"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
)

func GetMatchs(teamID int, param reqres.ReqPaging) (data reqres.ResPaging) {
	var responses []models.Match
	where := "deleted_at IS NULL"

	var modelTotal []models.Match

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

	var responsesRefined []reqres.MatchResponse
	for _, item := range responses {
		responseRefined := BuildMatchResponse(item)

		responsesRefined = append(responsesRefined, responseRefined)
	}

	data = utils.PopulateResPaging(&param, responsesRefined, totalResult.Total, totalFiltered, null.TimeFrom(totalResult.LastUpdated))

	return
}

func GetMatchByID(id uint) (responseRefined reqres.MatchResponse, err error) {
	var response models.Match
	err = config.DB.First(&response, id).Error

	responseRefined = BuildMatchResponse(response)

	return
}

func CreateMatch(Date null.Time,data *reqres.MatchRequest) (response models.Match, err error) {

	response = models.Match{
		AwayTeamID:         data.AwayTeamID,
		HomeTeamID:         data.HomeTeamID,
		Date:       Date,
		Status:       "scheduled",
		
	}

	 if data.HomeTeamID == data.AwayTeamID {
        return models.Match{}, errors.New("home team and away team cannot be the same")
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

func UpdateMatch(request models.Match) (response models.Match, err error) {
	err = config.DB.Save(&request).Scan(&response).Error

	return
}

func DeleteMatch(request models.Match) (models.Match, error) {
	err := config.DB.Delete(&request).Error
	return request, err
}

func CountMatch(teamID uint) (total int64) {
	where := "deleted_at IS NULL"

	if teamID > 0 {
		where += " AND team_id = " + strconv.Itoa(int(teamID))
	}

	var modelTotal []models.GlobalUser

	config.DB.Model(&modelTotal).Where(where).Count(&total)

	return total
}

func BuildMatchResponse(data models.Match) (response reqres.MatchResponse) {


	response.CustomGormModel = data.CustomGormModel
	

	return response
}
