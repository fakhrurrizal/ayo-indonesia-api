package reqres

import (
    validation "github.com/go-ozzo/ozzo-validation"
)

type TeamRequest struct {
    Name        string `json:"name"`
    Logo        string `json:"logo"`
    FoundedYear int    `json:"founded_year"`
    Address     string `json:"address"`
    City        string `json:"city"`
}

func (r TeamRequest) Validate() error {
    return validation.ValidateStruct(&r,
        validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
        validation.Field(&r.City, validation.Required, validation.Length(1, 100)),
        validation.Field(&r.FoundedYear, validation.Required, validation.Min(1800), validation.Max(2025)),
    )
}

type TeamResponse struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Logo        string `json:"logo"`
    FoundedYear int    `json:"founded_year"`
    Address     string `json:"address"`
    City        string `json:"city"`
    PlayersCount int   `json:"players_count"`
}
