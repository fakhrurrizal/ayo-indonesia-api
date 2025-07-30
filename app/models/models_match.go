package models

import "github.com/guregu/null"

type Match struct {
	CustomGormModel
	Date       null.Time `json:"date" gorm:"type:timestamptz"`
	HomeTeamID uint      `json:"home_team_id" gorm:"column:home_team_id;not null"`
	AwayTeamID uint      `json:"away_team_id" gorm:"column:away_team_id;not null"`
	HomeScore  *int      `json:"home_score" gorm:"column:home_score"`
	AwayScore  *int      `json:"away_score" gorm:"column:away_score"`
	Status     string    `json:"status" gorm:"column:status;default:'scheduled'"`
}
