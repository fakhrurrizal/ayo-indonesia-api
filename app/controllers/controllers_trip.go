package controllers

import (
	repository "ayo-indonesia-api/app/repositories"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"encoding/json"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

// CreateTrip godoc
// @Summary Create Trip
// @Description Create New Trip
// @Tags Trip
// @Produce json
// @Param Body body reqres.TripRequest true "Create body"
// @Success 200
// @Router /v1/trip [post]
// @Security ApiKeyAuth
// @Security JwtToken
func CreateTrip(c echo.Context) error {
	var input reqres.TripRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&input)

	if err := input.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(400, utils.NewInvalidInputError(errVal))
	}

	data, err := repository.CreateTrip(&input)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to create"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Create Success",
	})
}

// GetTrips godoc
// @Summary Get All Trip With Pagination
// @Description Get All Trip With Pagination
// @Tags Trip
// @Param search query string false "search (string)"
// @Param page query integer false "page (int)"
// @Param limit query integer false "limit (int)"
// @Param sort query string false "sort (ASC/DESC)"
// @Param order query string false "order by (default: id)"
// @Param status query boolean false "status (true (active) or false (inactive))"
// @Param category_id query integer false "category_id (int)"
// @Param app_id query integer false "app_id (int)"
// @Param destination_type_id query integer false "destination_type_id (int)"
// @Param created_at_margin_top query string false "created_at_margin_top (format: 2006-01-02)"
// @Param created_at_margin_bottom query string false "created_at_margin_top (format: 2006-01-02)"
// @Param code query string false "code (string)"
// @Produce json
// @Success 200
// @Router /v1/trip [get]
// @Security ApiKeyAuth
func GetTrips(c echo.Context) error {
	categoryId, _ := strconv.Atoi(c.QueryParam("category_id"))
	destinationTypeId, _ := strconv.Atoi(c.QueryParam("destination_type_id"))
	appId, _ := strconv.Atoi(c.QueryParam("app_id"))
	param := utils.PopulatePaging(c, "status")
	data := repository.GetTrips(categoryId, destinationTypeId, appId, param)

	return c.JSON(http.StatusOK, data)
}

// GetTripByID godoc
// @Summary Get Single Trip
// @Description Get Single Trip
// @Tags Trip
// @Param id path integer true "ID"
// @Produce json
// @Success 200
// @Router /v1/trip/{id} [get]
// @Security ApiKeyAuth
func GetTripByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetTripByID(id)
	if err != nil {
		return c.JSON(404, utils.Respond(404, err, "Failed to get"))
	}
	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to get",
	})
}

// DeleteTripByID godoc
// @Summary Delete Single Trip by ID
// @Description Delete Single Trip by ID
// @Tags Trip
// @Produce json
// @Param id path integer true "ID"
// @Success 200
// @Router /v1/trip/{id} [delete]
// @Security ApiKeyAuth
// @Security JwtToken
func DeleteTripByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetTripByIDPlain(id)
	if err != nil {
		return c.JSON(500, utils.Respond(404, err, "Failed to get"))
	}

	_, err = repository.DeleteTrip(data)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to delete"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to delete",
	})
}

// UpdateTripByID godoc
// @Summary Update Single Trip by ID
// @Description Update Single Trip by ID
// @Tags Trip
// @Produce json
// @Param id path integer true "ID"
// @Param Body body reqres.TripRequest true "Update body"
// @Success 200
// @Router /v1/trip/{id} [put]
// @Security ApiKeyAuth
// @Security JwtToken
func UpdateTripByID(c echo.Context) error {
	var input reqres.TripRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&input)

	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetTripByIDPlain(id)
	if err != nil {
		return c.JSON(500, utils.Respond(404, err, "Failed to get"))
	}

	images, _ := json.Marshal(input.Image)
	data.Image = string(images)

	if input.Name != "" {
		data.Name = input.Name
	}
	if input.Description != "" {
		data.Description = input.Description
	}

	if input.TripCategoryID != 0 {
		data.TripCategoryID = input.TripCategoryID
	}

	if input.DestinationTypeID != 0 {
		data.DestinationTypeID = input.DestinationTypeID
	}

	if input.BasePrice != 0 {
		data.BasePrice = input.BasePrice
	}

	if input.Location != "" {
		data.Location = input.Location
	}

	if input.Latitude != 0 {
		data.Latitude = input.Latitude
	}

	if input.Longitude != 0 {
		data.Longitude = input.Longitude
	}

	data.Status = input.Status

	data.IsActive = input.IsActive

	dataUpdate, err := repository.UpdateTrip(data)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to update"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    dataUpdate,
		"message": "Success to update",
	})
}
