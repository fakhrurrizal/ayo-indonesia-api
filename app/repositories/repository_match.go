package repository

import (
	"ayo-indonesia-api/app/models"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"ayo-indonesia-api/config"
	"errors"
	"fmt"

	"time"
)

func GetMatches(param reqres.ReqPaging) reqres.ResPaging {
	var matches []models.Match
	var total int64
	var totalFiltered int64

	query := config.DB.Model(&models.Match{})

	if param.Search != "" {
		query = query.Joins("JOIN teams ht ON ht.id = matches.home_team_id").
			Joins("JOIN teams at ON at.id = matches.away_team_id").
			Where("ht.name ILIKE ? OR at.name ILIKE ?", "%"+param.Search+"%", "%"+param.Search+"%")
	}

	query.Count(&totalFiltered)
	config.DB.Model(&models.Match{}).Count(&total)

	query.Limit(param.Limit).Offset(param.Offset).Order(param.Sort + " " + param.Order).
		Preload("HomeTeam").Preload("AwayTeam").Preload("Goals.Player.Team").Find(&matches)

	var responses []reqres.MatchResponse
	for _, match := range matches {
		responses = append(responses, BuildMatchResponse(match))
	}

	return utils.PopulateResPaging(&param, responses, total, totalFiltered)
}

func GetMatchByID(id uint) (models.Match, error) {
	var match models.Match
	err := config.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Goals.Player.Team").First(&match, id).Error
	return match, err
}

func CreateMatch(req reqres.MatchRequest) (models.Match, error) {
	if req.HomeTeamID == req.AwayTeamID {
		return models.Match{}, errors.New("home team and away team cannot be the same")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return models.Match{}, errors.New("invalid date format, use YYYY-MM-DD")
	}

	match := models.Match{
		Date:       date,
		Time:       req.Time,
		HomeTeamID: req.HomeTeamID,
		AwayTeamID: req.AwayTeamID,
		Status:     "scheduled",
	}

	err = config.DB.Create(&match).Error
	if err == nil {
		config.DB.Preload("HomeTeam").Preload("AwayTeam").First(&match, match.ID)
	}
	return match, err
}

func UpdateMatchResult(id uint, req reqres.MatchResultRequest) (models.Match, error) {
	var match models.Match
	if err := config.DB.First(&match, id).Error; err != nil {
		return match, err
	}

	for i, goalReq := range req.Goals {
		var player models.Player
		if err := config.DB.First(&player, goalReq.PlayerID).Error; err != nil {
			return match, fmt.Errorf("player with ID %d not found for goal %d", goalReq.PlayerID, i+1)
		}

		if player.TeamID != match.HomeTeamID && player.TeamID != match.AwayTeamID {
			return match, fmt.Errorf("player %s (ID: %d) does not belong to either team in this match", player.Name, goalReq.PlayerID)
		}

		if goalReq.Minute < 1 || goalReq.Minute > 120 {
			return match, fmt.Errorf("invalid minute %d for goal %d (must be between 1-120)", goalReq.Minute, i+1)
		}
	}

	tx := config.DB.Begin()

	match.HomeScore = &req.HomeScore
	match.AwayScore = &req.AwayScore
	match.Status = "completed"

	if err := tx.Save(&match).Error; err != nil {
		tx.Rollback()
		return match, err
	}

	if err := tx.Unscoped().Where("match_id = ?", id).Delete(&models.Goal{}).Error; err != nil {
		tx.Rollback()
		return match, err
	}

	// Create new goals
	for i, goalReq := range req.Goals {
		goal := models.Goal{
			MatchID:  id,
			PlayerID: goalReq.PlayerID,
			Minute:   goalReq.Minute,
		}
		if err := tx.Create(&goal).Error; err != nil {
			tx.Rollback()
			return match, fmt.Errorf("failed to create goal %d: %v", i+1, err)
		}
	}

	tx.Commit()

	// Reload match with relations
	config.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Goals.Player.Team").First(&match, match.ID)
	return match, nil
}

func DeleteMatch(id uint) error {
	return config.DB.Delete(&models.Match{}, id).Error
}

func GetMatchReport(id uint) (reqres.MatchReportResponse, error) {
	match, err := GetMatchByID(id)
	if err != nil {
		return reqres.MatchReportResponse{}, err
	}

	if match.Status != "completed" {
		return reqres.MatchReportResponse{}, errors.New("match not completed yet")
	}

	// Determine match result
	var matchResult string
	if *match.HomeScore > *match.AwayScore {
		matchResult = "Tim Home Menang"
	} else if *match.HomeScore < *match.AwayScore {
		matchResult = "Tim Away Menang"
	} else {
		matchResult = "Draw"
	}

	// Get top scorer in this match
	var topScorer *reqres.PlayerResponse
	if len(match.Goals) > 0 {
		goalCounts := make(map[uint]int)
		playerMap := make(map[uint]models.Player)

		for _, goal := range match.Goals {
			goalCounts[goal.PlayerID]++
			playerMap[goal.PlayerID] = goal.Player
		}

		var maxGoals int
		var topScorerID uint
		for playerID, count := range goalCounts {
			if count > maxGoals {
				maxGoals = count
				topScorerID = playerID
			}
		}

		if topScorerID > 0 {
			player := BuildPlayerResponse(playerMap[topScorerID])
			topScorer = &player
		}
	}

	// Get total wins for both teams
	var homeWins, awayWins int64
	config.DB.Model(&models.Match{}).Where("(home_team_id = ? AND home_score > away_score) OR (away_team_id = ? AND away_score > home_score)", match.HomeTeamID, match.HomeTeamID).Count(&homeWins)
	config.DB.Model(&models.Match{}).Where("(home_team_id = ? AND home_score > away_score) OR (away_team_id = ? AND away_score > home_score)", match.AwayTeamID, match.AwayTeamID).Count(&awayWins)

	return reqres.MatchReportResponse{
		Match:        BuildMatchResponse(match),
		MatchResult:  matchResult,
		TopScorer:    topScorer,
		HomeTeamWins: homeWins,
		AwayTeamWins: awayWins,
	}, nil
}

func BuildMatchResponse(match models.Match) reqres.MatchResponse {
	var goals []reqres.GoalResponse
	for _, goal := range match.Goals {
		goals = append(goals, reqres.GoalResponse{
			ID:     goal.ID,
			Player: BuildPlayerResponse(goal.Player),
			Minute: goal.Minute,
		})
	}

	return reqres.MatchResponse{
		ID:        match.ID,
		Date:      match.Date,
		Time:      match.Time,
		HomeTeam:  BuildTeamResponse(match.HomeTeam),
		AwayTeam:  BuildTeamResponse(match.AwayTeam),
		HomeScore: match.HomeScore,
		AwayScore: match.AwayScore,
		Status:    match.Status,
		Goals:     goals,
	}
}
