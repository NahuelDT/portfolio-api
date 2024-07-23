package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/NahuelDT/portfolio-api/internal/service"
)

type PortfolioHandler struct {
    portfolioService *service.PortfolioService
}

func NewPortfolioHandler(portfolioService *service.PortfolioService) *PortfolioHandler {
    return &PortfolioHandler{portfolioService: portfolioService}
}

func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
    userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    portfolio, err := h.portfolioService.GetPortfolio(uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, portfolio)
}

