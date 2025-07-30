package reqres

import (
	"ayo-indonesia-api/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/guregu/null"
)

type SignupRequest struct {
	Fullname       string                          `json:"fullname" validate:"required"`
	Email          string                          `json:"email" validate:"required"`
	Password       string                          `json:"password" validate:"required"`
}

func (request SignupRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Fullname, validation.Required),
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required),
	)
}

type SignUpRequest struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (request SignUpRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Length(8, 30)),
		validation.Field(&request.Fullname, validation.Required, validation.Length(5, 50)),
	)
}


type GlobalUserRequest struct {
	Fullname     string `json:"fullname" validate:"required"`
	Email        string `json:"email" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

func (request GlobalUserRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Fullname, validation.Required),
		validation.Field(&request.Email, validation.Required),
		validation.Field(&request.Password, validation.Required),
	)
}

type GlobalUserResponse struct {
	models.CustomGormModel
	Fullname              string               `json:"fullname"`
	Email                 string               `json:"email"`
	Password              string               `json:"-"`
	EmailVerifiedAt       null.Time            `json:"-"`
}