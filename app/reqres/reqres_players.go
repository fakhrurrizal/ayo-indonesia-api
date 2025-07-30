package reqres

import (
	"ayo-indonesia-api/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type PlayerRequest struct {
	Name         string  `json:"name"`
	Height       float64 `json:"height"`
	Weight       float64 `json:"weight"`
	Position     string  `json:"position"`
	JerseyNumber int     `json:"jersey_number"`
	TeamID       uint    `json:"team_id"`
}

func (r PlayerRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Position, validation.Required, validation.In("penyerang", "gelandang", "bertahan", "penjaga gawang")),
		validation.Field(&r.JerseyNumber, validation.Required, validation.Min(1), validation.Max(99)),
		validation.Field(&r.TeamID, validation.Required),
		validation.Field(&r.Height, validation.Min(1.0), validation.Max(3.0)),
		validation.Field(&r.Weight, validation.Min(30.0), validation.Max(200.0)),
	)
}

type PlayerResponse struct {
	models.CustomGormModel
	Name         string               `json:"name"`
	Height       float64              `json:"height"`
	Weight       float64              `json:"weight"`
	Position     string               `json:"position"`
	JerseyNumber int                  `json:"jersey_number"`
	Team         GlobalIDNameResponse `json:"team"`
}
