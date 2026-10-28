package controller

import (
	"net/http"
	"strconv"

	"order-management/model"
	"order-management/service"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	svc service.OrderService
}

func NewOrderController(svc service.OrderService) *OrderController {
	return &OrderController{svc: svc}
}

func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success:   false,
			Message:   "invalid request: " + err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	order, err := ctrl.svc.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success:   false,
			Message:   err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{
		Success:   true,
		Data:      order,
		RequestID: c.GetString("request_id"),
	})
}

func (ctrl *OrderController) GetOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "invalid order ID",
		})
		return
	}

	order, err := ctrl.svc.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success:   false,
			Message:   err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success:   true,
		Data:      order,
		RequestID: c.GetString("request_id"),
	})
}

func (ctrl *OrderController) ListOrders(c *gin.Context) {
	orders, err := ctrl.svc.ListOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success:   false,
			Message:   err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success:   true,
		Data:      orders,
		RequestID: c.GetString("request_id"),
	})
}

func (ctrl *OrderController) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "invalid order ID",
		})
		return
	}

	var req model.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "invalid request: " + err.Error(),
		})
		return
	}

	order, err := ctrl.svc.UpdateOrderStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success:   false,
			Message:   err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success:   true,
		Data:      order,
		RequestID: c.GetString("request_id"),
	})
}
