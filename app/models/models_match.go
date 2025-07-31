package models

import "time"

type Match struct {
	CustomGormModel
	Date       time.Time `json:"date" gorm:"column:date;not null"`
	Time       string    `json:"time" gorm:"column:time;not null"`
	HomeTeamID uint      `json:"home_team_id" gorm:"column:home_team_id;not null"`
	AwayTeamID uint      `json:"away_team_id" gorm:"column:away_team_id;not null"`
	HomeScore  *int      `json:"home_score" gorm:"column:home_score"`
	AwayScore  *int      `json:"away_score" gorm:"column:away_score"`
	Status     string    `json:"status" gorm:"column:status;default:'scheduled'"` // scheduled, completed
	HomeTeam   Team      `json:"home_team" gorm:"foreignKey:HomeTeamID"`
	AwayTeam   Team      `json:"away_team" gorm:"foreignKey:AwayTeamID"`
	Goals      []Goal    `json:"goals" gorm:"foreignKey:MatchID"`
}