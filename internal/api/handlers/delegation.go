package handlers

import (
	"net/http"
	"strconv"
	"tezos-delegation-service/internal/api"
	"tezos-delegation-service/internal/db"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for delegations.
type Handler struct {
	db db.Querier
}

// NewHandler creates a new Handler with the provided database querier.
func NewHandler(db db.Querier) *Handler {
	return &Handler{db: db}
}

// GetDelegations handles the request to get delegations for a specific year.
func (h *Handler) GetXtzDelegations(c *gin.Context, delegationsParams api.GetXtzDelegationsParams) {
	year := c.Query("year")

	// validator will intercept bad requests and yearInt will fallback to default value : 0
	yearInt, _ := strconv.Atoi(year)

	// Fetch delegations from the database for the given year.
	delegations, err := h.db.GetDelegationsByYear(c, yearInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "failed to fetch delegations"})
		return
	}

	// Respond with the fetched delegations.
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": delegations})
}
