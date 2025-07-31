package controllers

import (
	repository "ayo-indonesia-api/app/repositories"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get All Players
// @Description Get all players with pagination and optional team filter
// @Tags Players
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search by name or position"
// @Param team_id query int false "Filter by team ID"
// @Produce json
// @Success 200 {object} reqres.ResPaging
// @Router /v1/players [get]
func GetPlayers(c *gin.Context) {
	teamID, _ := strconv.ParseUint(c.Query("team_id"), 10, 32)
	param := utils.PopulatePaging(c, "status")
	data := repository.GetPlayers(uint(teamID), param)
	c.JSON(http.StatusOK, data)
}

// @Summary Get Player by ID
// @Description Get player details by ID
// @Tags Players
// @Param id path int true "Player ID"
// @Produce json
// @Success 200 {object} reqres.PlayerResponse
// @Router /v1/players/{id} [get]
func GetPlayerByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	player, err := repository.GetPlayerByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	response := repository.BuildPlayerResponse(player)
	c.JSON(http.StatusOK, response)
}

// @Summary Create Player
// @Description Create a new player
// @Tags Players
// @Accept json
// @Produce json
// @Param player body reqres.PlayerRequest true "Player data"
// @Success 201 {object} reqres.PlayerResponse
// @Router /v1/players [post]
// @Security ApiKeyAuth
// @Security JwtToken
func CreatePlayer(c *gin.Context) {
	var req reqres.PlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := repository.CreatePlayer(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := repository.BuildPlayerResponse(player)
	c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    response,
		"message": "Create Success",
	})
}

// @Summary Update Player
// @Description Update player by ID
// @Tags Players
// @Accept json
// @Produce json
// @Param id path int true "Player ID"
// @Param player body reqres.PlayerRequest true "Player data"
// @Success 200 {object} reqres.PlayerResponse
// @Router /v1/players/{id} [put]
// @Security ApiKeyAuth
// @Security JwtToken
func UpdatePlayer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	var req reqres.PlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := repository.UpdatePlayer(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := repository.BuildPlayerResponse(player)
	c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    response,
		"message": "Create Success",
	})
}

// @Summary Delete Player
// @Description Delete player by ID (soft delete)
// @Tags Players
// @Param id path int true "Player ID"
// @Success 200 {object} map[string]string
// @Router /v1/players/{id} [delete]
// @Security ApiKeyAuth
// @Security JwtToken
func DeletePlayer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	if err := repository.DeletePlayer(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, map[string]interface{}{
		"status":  200,
		"message": "Player deleted successfull",
	})
}
