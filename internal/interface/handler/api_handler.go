package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/req"
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

// SaveToken godoc
// @Summary      トークンを保存します。
// @Description
// @Tags         token
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "全顧客同期完了"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/token [post]
func (h *APIHandler) SaveToken(c echo.Context) error {

	var token req.Token
	if err := c.Bind(&token); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := h.customerUsecase.SaveToken(c.Request().Context(), token.Token)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "sync all")
}

func handleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, err.Error())
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}
