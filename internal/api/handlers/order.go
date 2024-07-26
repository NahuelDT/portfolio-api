package handlers

import (
	"net/http"
	"strconv"

	"github.com/NahuelDT/portfolio-api/internal/models"
	"github.com/NahuelDT/portfolio-api/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) PlaceOrder(c *gin.Context) {
	var orderRequest struct {
		Order       models.Order `json:"order"`
		TotalAmount float64      `json:"totalAmount"`
	}
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderService.PlaceOrder(&orderRequest.Order, orderRequest.TotalAmount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderRequest.Order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := h.orderService.CancelOrder(uint(orderID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}
