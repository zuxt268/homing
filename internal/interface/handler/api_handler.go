package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zuxt268/homing/internal/domain"
	_ "github.com/zuxt268/homing/internal/interface/dto/res"
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
// @Summary      全顧客データ同期
// @Description  全ての顧客のデータを同期します
// @Tags         sync
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "全顧客同期完了"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/sync [post]
func (h *APIHandler) SyncAll(c echo.Context) error {
	err := h.customerUsecase.SyncAll(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "sync all")
}

// SyncOne godoc
// @Summary      指定顧客データ同期
// @Description  指定された顧客IDのデータを同期します
// @Tags         sync
// @Accept       json
// @Produce      json
// @Param        customer_id  path      int     true  "顧客ID"
// @Success      200          {string}  string  "顧客同期完了"
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
// @Summary      顧客情報取得
// @Description  指定された顧客IDの情報を取得します
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customer_id  path      int     true  "顧客ID"
// @Success      200          {object}  res.Customer  "顧客情報"
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

// GetInstagramAccount godoc
// @Summary      Instagramアカウント情報取得
// @Description  指定された顧客のInstagramアカウント情報を取得します
// @Tags         instagram
// @Accept       json
// @Produce      json
// @Param        customer_id  path      int     true  "顧客ID"
// @Success      200          {object}  res.InstagramAccounts  "Instagramアカウント情報"
// @Router       /api/instagram/{customer_id} [get]
func (h *APIHandler) GetInstagramAccount(c echo.Context) error {
	customerIDStr := c.Param("customer_id")
	if customerIDStr == "" {
		return c.JSON(http.StatusBadRequest, "customer_id is required")
	}
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid customer_id format")
	}
	resp, err := h.customerUsecase.GetInstagramAccount(c.Request().Context(), customerID)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

// SyncInstagramAccount godoc
// @Summary      Instagramアカウント同期
// @Description  指定された顧客のInstagramアカウントを同期します
// @Tags         instagram
// @Accept       json
// @Produce      json
// @Param        customer_id  path      int     true  "顧客ID"
// @Success      200          {object}  res.InstagramAccounts  "Instagramアカウント情報"
// @Router       /api/instagram/sync/{customer_id} [post]
func (h *APIHandler) SyncInstagramAccount(c echo.Context) error {
	customerIDStr := c.Param("customer_id")
	if customerIDStr == "" {
		return c.JSON(http.StatusBadRequest, "customer_id is required")
	}
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid customer_id format")
	}
	resp, err := h.customerUsecase.SyncAccount(c.Request().Context(), customerID)
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
