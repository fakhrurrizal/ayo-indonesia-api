package models

type Goal struct {
    CustomGormModel
    MatchID    uint   `json:"match_id" gorm:"type: int8"`
    PlayerID   uint   `json:"player_id" gorm:"type: int8"`
    Minute     int    `json:"minute" gorm:"column:minute;not null"`
}