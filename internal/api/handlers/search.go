package handlers

import (
	"net/http"

	"github.com/NahuelDT/portfolio-api/internal/service"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *service.SearchService
}

func NewSearchHandler(searchService *service.SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

func (h *SearchHandler) SearchAssets(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	results, err := h.searchService.SearchAssets(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search assets"})
		return
	}

	c.JSON(http.StatusOK, results)
}
