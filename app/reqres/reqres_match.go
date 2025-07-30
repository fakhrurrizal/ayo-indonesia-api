package reqres

import (
	"ayo-indonesia-api/app/models"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type MatchRequest struct {
	Date       string `json:"date"`
	HomeTeamID uint   `json:"home_team_id"`
	AwayTeamID uint   `json:"away_team_id"`
}

func (r MatchRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Date, validation.Required),
		validation.Field(&r.HomeTeamID, validation.Required),
		validation.Field(&r.AwayTeamID, validation.Required),
	)
}

type MatchResultRequest struct {
	HomeScore int           `json:"home_score"`
	AwayScore int           `json:"away_score"`
	Goals     []GoalRequest `json:"goals"`
}

type GoalRequest struct {
	PlayerID uint `json:"player_id"`
	Minute   int  `json:"minute"`
}

type MatchResponse struct {
	models.CustomGormModel
	Date      time.Time      `json:"date"`
	HomeTeam  TeamResponse   `json:"home_team"`
	AwayTeam  TeamResponse   `json:"away_team"`
	HomeScore *int           `json:"home_score"`
	AwayScore *int           `json:"away_score"`
	Status    string         `json:"status"`
	Goals     []GoalResponse `json:"goals"`
}

type GoalResponse struct {
	models.CustomGormModel
	Player PlayerResponse `json:"player"`
	Minute int            `json:"minute"`
}
