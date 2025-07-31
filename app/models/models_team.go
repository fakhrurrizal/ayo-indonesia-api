package models

type Team struct {
    CustomGormModel
    Name         string    `json:"name" gorm:"column:name;not null"`
    Logo         string    `json:"logo" gorm:"column:logo"`
    FoundedYear  int       `json:"founded_year" gorm:"column:founded_year"`
    Address      string    `json:"address" gorm:"column:address"`
    City         string    `json:"city" gorm:"column:city"`
    Players      []Player  `json:"players" gorm:"foreignKey:TeamID"`
    HomeMatches  []Match   `json:"home_matches" gorm:"foreignKey:HomeTeamID"`
    AwayMatches  []Match   `json:"away_matches" gorm:"foreignKey:AwayTeamID"`
}