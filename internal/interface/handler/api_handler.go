package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIHandler struct {
}

func NewAPIHandler() APIHandler {
	return APIHandler{}
}

func (h *APIHandler) SyncByCustomer(c echo.Context) error {
	return c.JSON(http.StatusOK, "sync by customer")
}

func (h *APIHandler) SyncAll(c echo.Context) error {
	return c.JSON(http.StatusOK, "sync all")
}
