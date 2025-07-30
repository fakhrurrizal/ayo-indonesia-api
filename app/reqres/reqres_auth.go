package reqres

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type SignInRequest struct {
	Email      string                          `json:"email" validate:"required"`
	Password   string                          `json:"password" validate:"required"`
}

func (request SignInRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required),
		validation.Field(&request.Password, validation.Required),
	)
}


type GlobalUserAuthResponse struct {
	ID            int                  `json:"id"`
	Fullname      string               `json:"fullname"`
	Email         string               `json:"email"`
}


type GlobalIDNameResponse struct {
	ID                 int     `json:"id,omitempty"`
	Name               string  `json:"name,omitempty"`
	Code               string  `json:"code,omitempty"`
}
