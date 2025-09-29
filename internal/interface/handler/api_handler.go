package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/usecase"
)

type APIHandler struct {
	customerUsecase usecase.CustomerUsecase
}

func NewAPIHandler(customerUsecase usecase.CustomerUsecase) APIHandler {
	return APIHandler{
		customerUsecase: customerUsecase,
	}
}

// SyncAll godoc
// @Summary      Sync all customers
// @Description  Sync data for all customers
// @Tags         sync
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "sync all"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/sync [post]
func (h *APIHandler) SyncAll(c echo.Context) error {
	err := h.customerUsecase.SyncAll(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "sync all")
}

// SyncOne godoc
// @Summary      Sync one customer
// @Description  Sync data for a specific customer by ID
// @Tags         sync
// @Accept       json
// @Produce      json
// @Param        customer_id  path      int     true  "Customer ID"
// @Success      200          {string}  string  "sync one customer"
// @Failure      400          {string}  string  "Bad request"
// @Failure      404          {string}  string  "Customer not found"
// @Failure      500          {string}  string  "Internal server error"
// @Router       /api/sync/{customer_id} [post]
func (h *APIHandler) SyncOne(c echo.Context) error {
	customerID := c.Param("customer_id")
	if customerID == "" {
		return c.JSON(http.StatusBadRequest, "customer_id is required")
	}

	// Convert string to int
	var id int
	if _, err := fmt.Sscanf(customerID, "%d", &id); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid customer_id format")
	}

	err := h.customerUsecase.SyncOne(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "sync one customer")
}

// GetCustomer godoc
// @Summary      Get customer by ID
// @Description  Get customer information by customer ID
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customer_id  path      int     true  "Customer ID"
// @Success      200          {object}  interface{}  "Customer information"
// @Failure      400          {string}  string  "Bad request"
// @Failure      404          {string}  string  "Customer not found"
// @Failure      500          {string}  string  "Internal server error"
// @Router       /api/customers/{customer_id} [get]
func (h *APIHandler) GetCustomer(c echo.Context) error {
	customerIDStr := c.Param("customer_id")
	if customerIDStr == "" {
		return c.JSON(http.StatusBadRequest, "customer_id is required")
	}
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid customer_id format")
	}

	resp, err := h.customerUsecase.GetCustomer(c.Request().Context(), customerID)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

func handleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, err.Error())
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}
