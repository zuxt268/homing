package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/dto/req"
	"github.com/zuxt268/homing/internal/interface/dto/res"
	_ "github.com/zuxt268/homing/internal/interface/dto/res"
	"github.com/zuxt268/homing/internal/usecase"
)

type APIHandler struct {
	customerUsecase           usecase.CustomerUsecase
	tokenUsecase              usecase.TokenUsecase
	wordpressInstagramUsecase usecase.WordpressInstagramUsecase
}

func NewAPIHandler(
	customerUsecase usecase.CustomerUsecase,
	tokenUsecase usecase.TokenUsecase,
	wordpressInstagramUsecase usecase.WordpressInstagramUsecase,
) APIHandler {
	return APIHandler{
		customerUsecase:           customerUsecase,
		tokenUsecase:              tokenUsecase,
		wordpressInstagramUsecase: wordpressInstagramUsecase,
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
// @Summary      全顧客データ同期
// @Description  全ての顧客のデータを同期します
// @Tags         sync
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "全顧客同期完了"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/sync/{id} [post]
func (h *APIHandler) SyncOne(c echo.Context) error {
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err := h.customerUsecase.SyncOne(c.Request().Context(), id)
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
// @Success      200  {string}  string  "更新完了"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/token [post]
func (h *APIHandler) SaveToken(c echo.Context) error {

	var token req.UpdateToken
	if err := c.Bind(&token); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := h.tokenUsecase.UpdateToken(c.Request().Context(), token)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "ok")
}

// GetToken godoc
// @Summary      トークンを取得します。
// @Description
// @Tags         token
// @Accept       json
// @Produce      json
// @Success 200 {object} res.Token "トークン情報"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/token [get]
func (h *APIHandler) GetToken(c echo.Context) error {
	token, err := h.tokenUsecase.GetToken(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, token)
}

// CheckToken godoc
// @Summary      トークンの認証情報を取得する
// @Description
// @Tags         token
// @Accept       json
// @Produce      json
// @Success 200 {object} string "ok"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/token/check [post]
func (h *APIHandler) CheckToken(c echo.Context) error {
	err := h.tokenUsecase.CheckToken(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, "ok")
}

// GetWordpressInstagramList godoc
// @Summary      Wordpress Instagram一覧取得
// @Description  Wordpress Instagramの一覧を取得します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Param        limit          query     int     false  "取得件数"
// @Param        offset         query     int     false  "オフセット"
// @Param        name           query     string  false  "名前"
// @Param        wordpress      query     string  false  "WordPress URL"
// @Param        instagram_id   query     string  false  "Instagram ID"
// @Param        status         query     int     false  "ステータス"
// @Param        delete_hash    query     bool    false  "削除ハッシュ"
// @Param        customer_type  query     int     false  "顧客タイプ"
// @Success      200  {object}  res.WordpressInstagramList  "Wordpress Instagram一覧"
// @Failure      400  {string}  string  "不正なリクエスト"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram [get]
func (h *APIHandler) GetWordpressInstagramList(c echo.Context) error {
	var params req.GetWordpressInstagram
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	list, err := h.wordpressInstagramUsecase.GetWordpressInstagramList(c.Request().Context(), params)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, list)
}

// GetWordpressInstagram godoc
// @Summary      Wordpress Instagram詳細取得
// @Description  Wordpress Instagramの詳細を取得します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Wordpress Instagram ID"
// @Success      200  {object}  res.WordpressInstagram  "Wordpress Instagram詳細"
// @Failure      404  {string}  string  "見つかりません"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram/{id} [get]
func (h *APIHandler) GetWordpressInstagram(c echo.Context) error {
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	item, err := h.wordpressInstagramUsecase.GetWordpressInstagram(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, item)
}

// CreateWordpressInstagram godoc
// @Summary      Wordpress Instagram作成
// @Description  Wordpress Instagramを作成します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Param        body  body      req.CreateWordpressInstagram  true  "作成データ"
// @Success      201   {object}  res.WordpressInstagram  "作成されたWordpress Instagram"
// @Failure      400   {string}  string  "不正なリクエスト"
// @Failure      500   {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram [post]
func (h *APIHandler) CreateWordpressInstagram(c echo.Context) error {
	var body req.CreateWordpressInstagram
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println("CreateWordpressInstagram", body)

	item, err := h.wordpressInstagramUsecase.CreateWordpressInstagram(c.Request().Context(), body)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, item)
}

// UpdateWordpressInstagram godoc
// @Summary      Wordpress Instagram更新
// @Description  Wordpress Instagramを更新します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Param        id    path      int                           true  "Wordpress Instagram ID"
// @Param        body  body      req.UpdateWordpressInstagram  true  "更新データ"
// @Success      200   {object}  res.WordpressInstagram  "更新されたWordpress Instagram"
// @Failure      400   {string}  string  "不正なリクエスト"
// @Failure      404   {string}  string  "見つかりません"
// @Failure      500   {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram/{id} [put]
func (h *APIHandler) UpdateWordpressInstagram(c echo.Context) error {
	var body req.UpdateWordpressInstagram
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	body.ID = &id

	item, err := h.wordpressInstagramUsecase.UpdateWordpressInstagram(c.Request().Context(), body)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, item)
}

// DeleteWordpressInstagram godoc
// @Summary      Wordpress Instagram削除
// @Description  Wordpress Instagramを削除します
// @Tags         wordpress-instagram
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Wordpress Instagram ID"
// @Success      204  {string}  string  "削除成功"
// @Failure      404  {string}  string  "見つかりません"
// @Failure      500  {string}  string  "内部サーバーエラー"
// @Router       /api/wordpress-instagram/{id} [delete]
func (h *APIHandler) DeleteWordpressInstagram(c echo.Context) error {
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := h.wordpressInstagramUsecase.DeleteWordpressInstagram(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func handleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, res.ErrorResponse{Message: err.Error()})
	case errors.Is(err, domain.ErrBadRequest):
		return c.JSON(http.StatusBadRequest, res.ErrorResponse{Message: err.Error()})
	case errors.Is(err, domain.ErrWordpressConnection):
		return c.JSON(http.StatusBadRequest, res.ErrorResponse{Message: err.Error()})
	case errors.Is(err, domain.ErrInstagramConnection):
		return c.JSON(http.StatusBadRequest, res.ErrorResponse{Message: err.Error()})
	default:
		return c.JSON(http.StatusInternalServerError, res.ErrorResponse{Message: err.Error()})
	}
}
