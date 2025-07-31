package controllers

import (
	repository "ayo-indonesia-api/app/repositories"
	"ayo-indonesia-api/app/reqres"
	"ayo-indonesia-api/app/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get All Matches
// @Description Get all matches with pagination
// @Tags Matches
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search by team names"
// @Produce json
// @Success 200 {object} reqres.ResPaging
// @Router /v1/matches [get]
func GetMatches(c *gin.Context) {
	param := utils.PopulatePaging(c, "status")
	data := repository.GetMatches(param)
	c.JSON(http.StatusOK, data)
}

// @Summary Get Match by ID
// @Description Get match details by ID
// @Tags Matches
// @Param id path int true "Match ID"
// @Produce json
// @Success 200 {object} reqres.MatchResponse
// @Router /v1/matches/{id} [get]
func GetMatchByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	match, err := repository.GetMatchByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	response := repository.BuildMatchResponse(match)
	c.JSON(http.StatusOK, response)
}

// @Summary Create Match
// @Description Create a new match schedule
// @Tags Matches
// @Accept json
// @Produce json
// @Param match body reqres.MatchRequest true "Match data"
// @Success 201 {object} reqres.MatchResponse
// @Router /v1/matches [post]
// @Security ApiKeyAuth
// @Security JwtToken
func CreateMatch(c *gin.Context) {
	var req reqres.MatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := repository.CreateMatch(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := repository.BuildMatchResponse(match)
	c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    response,
		"message": "Create Success",
	})
}

// @Summary Update Match Result
// @Description Update match result with scores and goals
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Param result body reqres.MatchResultRequest true "Match result data"
// @Success 200 {object} reqres.MatchResponse
// @Router /v1/matches/{id}/result [put]
// @Security ApiKeyAuth
// @Security JwtToken
func UpdateMatchResult(c *gin.Context) {
 id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
        return
    }

    var req reqres.MatchResultRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
        return
    }

    if req.HomeScore < 0 || req.AwayScore < 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Scores cannot be negative"})
        return
    }

    match, err := repository.UpdateMatchResult(uint(id), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	response := repository.BuildMatchResponse(match)

	c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    response,
		"message": "Update Success",
	})
}

// @Summary Delete Match
// @Description Delete match by ID (soft delete)
// @Tags Matches
// @Param id path int true "Match ID"
// @Success 200 {object} map[string]string
// @Router /v1/matches/{id} [delete]
// @Security ApiKeyAuth
// @Security JwtToken
func DeleteMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	if err := repository.DeleteMatch(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, map[string]interface{}{
		"status":  200,
		"message": "Match deleted successfully",
	})
	c.JSON(http.StatusOK, gin.H{"message": "Match deleted successfully"})
}

// @Summary Get Match Report
// @Description Get detailed match report with statistics
// @Tags Matches
// @Param id path int true "Match ID"
// @Produce json
// @Success 200 {object} reqres.MatchReportResponse
// @Router /v1/matches/{id}/report [get]
func GetMatchReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	report, err := repository.GetMatchReport(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
