package repository

import (
	"ayo-indonesia-api/app/models"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"ayo-indonesia-api/config"
)

func GetTeams(param reqres.ReqPaging) reqres.ResPaging {
	var teams []models.Team
	var total int64
	var totalFiltered int64

	query := config.DB.Model(&models.Team{})

	if param.Search != "" {
		query = query.Where("name ILIKE ? OR city ILIKE ?", "%"+param.Search+"%", "%"+param.Search+"%")
	}

	query.Count(&totalFiltered)
	config.DB.Model(&models.Team{}).Count(&total)

	query.Limit(param.Limit).Offset(param.Offset).Order(param.Sort + " " + param.Order).
		Preload("Players").Find(&teams)

	var responses []reqres.TeamResponse
	for _, team := range teams {
		responses = append(responses, BuildTeamResponse(team))
	}

	return utils.PopulateResPaging(&param, responses, total, totalFiltered)
}

func GetTeamByID(id uint) (models.Team, error) {
	var team models.Team
	err := config.DB.Preload("Players").First(&team, id).Error
	return team, err
}

func CreateTeam(req reqres.TeamRequest) (models.Team, error) {
	team := models.Team{
		Name:        req.Name,
		Logo:        req.Logo,
		FoundedYear: req.FoundedYear,
		Address:     req.Address,
		City:        req.City,
	}
	err := config.DB.Create(&team).Error
	return team, err
}

func UpdateTeam(id uint, req reqres.TeamRequest) (models.Team, error) {
	var team models.Team
	if err := config.DB.First(&team, id).Error; err != nil {
		return team, err
	}

	team.Name = req.Name
	team.Logo = req.Logo
	team.FoundedYear = req.FoundedYear
	team.Address = req.Address
	team.City = req.City

	err := config.DB.Save(&team).Error
	return team, err
}

func DeleteTeam(id uint) error {
	return config.DB.Delete(&models.Team{}, id).Error
}

func BuildTeamResponse(team models.Team) reqres.TeamResponse {
	return reqres.TeamResponse{
		ID:           team.ID,
		Name:         team.Name,
		Logo:         team.Logo,
		FoundedYear:  team.FoundedYear,
		Address:      team.Address,
		City:         team.City,
		PlayersCount: len(team.Players),
	}
}
