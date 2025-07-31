package models

type Goal struct {
    CustomGormModel
    MatchID    uint   `json:"match_id" gorm:"column:match_id;not null"`
    PlayerID   uint   `json:"player_id" gorm:"column:player_id;not null"`
    Minute     int    `json:"minute" gorm:"column:minute;not null"`
    Match      Match  `json:"match" gorm:"foreignKey:MatchID"`
    Player     Player `json:"player" gorm:"foreignKey:PlayerID"`
}
