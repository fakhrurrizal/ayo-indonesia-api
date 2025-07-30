package reqres

import (
	"ayo-indonesia-api/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type TripRequest struct {
	Name              string   `json:"name"`
	Image             []string `json:"image" validate:"required"`
	Description       string   `json:"description"`
	Status            bool     `json:"status"`
	TripCategoryID    int      `json:"trip_category_id"`
	DestinationTypeID int      `json:"destination_type_id"`
	BasePrice         float64  `json:"base_price"`
	IsActive          bool     `json:"is_active"`
	Location          string   `json:"location" `
	Latitude          float64  `json:"latitude" `
	Longitude         float64  `json:"longitude" `
	AppID             int      `json:"app_id"`
}

func (request TripRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
		validation.Field(&request.TripCategoryID, validation.Required),
		validation.Field(&request.BasePrice, validation.Required),
		validation.Field(&request.Image, validation.Required),
	)
}

type TripResponse struct {
	models.CustomGormModel
	Name            string               `json:"name"`
	Image           []string             `json:"image"`
	Description     string               `json:"description"`
	Status          bool                 `json:"status"`
	TripCategory    GlobalIDNameResponse `json:"trip_category"`
	DestinationType GlobalIDNameResponse `json:"destination_type"`
	BasePrice       float64              `json:"base_price"`
	IsActive        bool                 `json:"is_active"`
	Location        string               `json:"location" `
	Latitude        float64              `json:"latitude" `
	Longitude       float64              `json:"longitude" `
	AppID           int                  `json:"app_id"`
}
