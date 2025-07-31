package controllers

import (
	repository "ayo-indonesia-api/app/repositories"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get All Teams
// @Description Get all teams with pagination
// @Tags Teams
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search by name or city"
// @Produce json
// @Success 200 {object} reqres.ResPaging
// @Router /v1/teams [get]
func GetTeams(c *gin.Context) {
	param := utils.PopulatePaging(c, "status")
	data := repository.GetTeams(param)
	c.JSON(http.StatusOK, data)
}

// @Summary Get Team by ID
// @Description Get team details by ID
// @Tags Teams
// @Param id path int true "Team ID"
// @Produce json
// @Success 200 {object} reqres.TeamResponse
// @Router /v1/teams/{id} [get]
func GetTeamByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	team, err := repository.GetTeamByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	response := repository.BuildTeamResponse(team)
	c.JSON(http.StatusOK, response)
}

// @Summary Create Team
// @Description Create a new team
// @Tags Teams
// @Accept json
// @Produce json
// @Param team body reqres.TeamRequest true "Team data"
// @Success 200 {object} reqres.TeamResponse
// @Router /v1/teams [post]
// @Security ApiKeyAuth
// @Security JwtToken
func CreateTeam(c *gin.Context) {
	var req reqres.TeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := repository.CreateTeam(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := repository.BuildTeamResponse(team)
	c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    response,
		"message": "Create Success",
	})
}

// @Summary Update Team
// @Description Update team by ID
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Param team body reqres.TeamRequest true "Team data"
// @Success 200 {object} reqres.TeamResponse
// @Router /v1/teams/{id} [put]
// @Security ApiKeyAuth
// @Security JwtToken
func UpdateTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req reqres.TeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := repository.UpdateTeam(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := repository.BuildTeamResponse(team)
	c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    response,
		"message": "Update Success",
	})
}

// @Summary Delete Team
// @Description Delete team by ID (soft delete)
// @Tags Teams
// @Param id path int true "Team ID"
// @Success 200 {object} map[string]string
// @Router /v1/teams/{id} [delete]
// @Security ApiKeyAuth
// @Security JwtToken
func DeleteTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	if err := repository.DeleteTeam(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, map[string]interface{}{
		"status":  200,
		"message": "Team deleted successfully",
	})
}
