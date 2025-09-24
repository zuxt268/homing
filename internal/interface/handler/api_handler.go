package handler

import (
	"errors"
	"fmt"
	"net/http"

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

func (h *APIHandler) SyncByCustomer(c echo.Context) error {
	err := h.customerUsecase.SyncOne(c.Request().Context(), 1) // TODO
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "sync by customer")
}

func (h *APIHandler) SyncAll(c echo.Context) error {
	err := h.customerUsecase.SyncAll(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "sync all")
}

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

func handleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, err.Error())
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}
