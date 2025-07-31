package models

type Player struct {
    CustomGormModel
    Name         string  `json:"name" gorm:"column:name;not null"`
    Height       float64 `json:"height" gorm:"column:height"`
    Weight       float64 `json:"weight" gorm:"column:weight"`
    Position     string  `json:"position" gorm:"column:position;not null"` // penyerang, gelandang, bertahan, penjaga gawang
    JerseyNumber int     `json:"jersey_number" gorm:"column:jersey_number;not null"`
    TeamID       uint    `json:"team_id" gorm:"column:team_id;not null"`
    Team         Team    `json:"team" gorm:"foreignKey:TeamID"`
    Goals        []Goal  `json:"goals" gorm:"foreignKey:PlayerID"`
}