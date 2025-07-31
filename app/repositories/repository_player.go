package repository

import (
	"ayo-indonesia-api/app/models"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"ayo-indonesia-api/config"
	"errors"
)

func GetPlayers(teamID uint, param reqres.ReqPaging) reqres.ResPaging {
	var players []models.Player
	var total int64
	var totalFiltered int64

	query := config.DB.Model(&models.Player{})

	if teamID > 0 {
		query = query.Where("team_id = ?", teamID)
	}

	if param.Search != "" {
		query = query.Where("name ILIKE ? OR position ILIKE ?", "%"+param.Search+"%", "%"+param.Search+"%")
	}

	query.Count(&totalFiltered)
	config.DB.Model(&models.Player{}).Count(&total)

	query.Limit(param.Limit).Offset(param.Offset).Order(param.Sort + " " + param.Order).
		Preload("Team").Find(&players)

	var responses []reqres.PlayerResponse
	for _, player := range players {
		responses = append(responses, BuildPlayerResponse(player))
	}

	return utils.PopulateResPaging(&param, responses, total, totalFiltered)
}

func GetPlayerByID(id uint) (models.Player, error) {
	var player models.Player
	err := config.DB.Preload("Team").First(&player, id).Error
	return player, err
}

func CreatePlayer(req reqres.PlayerRequest) (models.Player, error) {
	var existingPlayer models.Player
	if err := config.DB.Where("team_id = ? AND jersey_number = ?", req.TeamID, req.JerseyNumber).First(&existingPlayer).Error; err == nil {
		return models.Player{}, errors.New("jersey number already exists in this team")
	}

	player := models.Player{
		Name:         req.Name,
		Height:       req.Height,
		Weight:       req.Weight,
		Position:     req.Position,
		JerseyNumber: req.JerseyNumber,
		TeamID:       req.TeamID,
	}
	err := config.DB.Create(&player).Error
	if err == nil {
		config.DB.Preload("Team").First(&player, player.ID)
	}
	return player, err
}

func UpdatePlayer(id uint, req reqres.PlayerRequest) (models.Player, error) {
	var player models.Player
	if err := config.DB.First(&player, id).Error; err != nil {
		return player, err
	}

	// Check jersey number uniqueness if changed
	if player.JerseyNumber != req.JerseyNumber || player.TeamID != req.TeamID {
		var existingPlayer models.Player
		if err := config.DB.Where("team_id = ? AND jersey_number = ? AND id != ?", req.TeamID, req.JerseyNumber, id).First(&existingPlayer).Error; err == nil {
			return models.Player{}, errors.New("jersey number already exists in this team")
		}
	}

	player.Name = req.Name
	player.Height = req.Height
	player.Weight = req.Weight
	player.Position = req.Position
	player.JerseyNumber = req.JerseyNumber
	player.TeamID = req.TeamID

	err := config.DB.Save(&player).Error
	if err == nil {
		config.DB.Preload("Team").First(&player, player.ID)
	}
	return player, err
}

func DeletePlayer(id uint) error {
	return config.DB.Delete(&models.Player{}, id).Error
}

func BuildPlayerResponse(player models.Player) reqres.PlayerResponse {
	return reqres.PlayerResponse{
		ID:           player.ID,
		Name:         player.Name,
		Height:       player.Height,
		Weight:       player.Weight,
		Position:     player.Position,
		JerseyNumber: player.JerseyNumber,
		Team:         BuildTeamResponse(player.Team),
	}
}
